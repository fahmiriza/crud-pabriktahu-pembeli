package endpoint

import (
	"context"
	"time"

	svc "MiniProject/git.bluebird.id/mini/pembeli/server"

	pb "MiniProject/git.bluebird.id/mini/pembeli/grpc"

	util "MiniProject/git.bluebird.id/mini/util/grpc"
	disc "MiniProject/git.bluebird.id/mini/util/microservice"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	grpcName = "grpc.PembeliService"
)

func NewGRPCPembeliClient(nodes []string, creds credentials.TransportCredentials, option util.ClientOption,
	tracer stdopentracing.Tracer, logger log.Logger) (svc.PembeliService, error) {

	instancer, err := disc.ServiceDiscovery(nodes, svc.ServiceID, logger)
	if err != nil {
		return nil, err
	}

	retryMax := option.Retry
	retryTimeout := option.RetryTimeout
	timeout := option.Timeout

	var addPembeliEp endpoint.Endpoint
	{
		factory := util.EndpointFactory(makeClientAddPembeliEndpoint, creds, timeout, tracer, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		addPembeliEp = retry
	}

	var readPembeliByKotaEp endpoint.Endpoint
	{
		factory := util.EndpointFactory(makeClientReadPembeliByKotaEndpoint, creds, timeout, tracer, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		readPembeliByKotaEp = retry
	}

	var readPembeliEp endpoint.Endpoint
	{
		factory := util.EndpointFactory(makeClientReadPembeliEndpoint, creds, timeout, tracer, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		readPembeliEp = retry
	}

	var updatePembeliEp endpoint.Endpoint
	{
		factory := util.EndpointFactory(makeClientUpdatePembeli, creds, timeout, tracer, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		updatePembeliEp = retry
	}
	return PembeliEndpoint{AddPembeliEndpoint: addPembeliEp, ReadPembeliByKotaEndpoint: readPembeliByKotaEp,
		ReadPembeliEndpoint: readPembeliEp, UpdatePembeliEndpoint: updatePembeliEp}, nil
}

func encodeAddPembeliRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(svc.Pembeli)
	return &pb.AddPembeliReq{
		IDPembeli:   req.IDPembeli,
		NamaPembeli: req.NamaPembeli,
		NoTelpon:    req.NoTelpon,
		KodeKota:    req.KodeKota,
		Status:      req.Status,
		CreateBy:    req.CreateBy,
	}, nil
}

func encodeReadPembeliByKotaRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(svc.Pembeli)
	return &pb.ReadPembeliByKotaReq{KodeKota: req.KodeKota}, nil
}

func encodeReadPembeliRequest(_ context.Context, request interface{}) (interface{}, error) {
	return &google_protobuf.Empty{}, nil
}

func encodeUpdatePembeliRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(svc.Pembeli)
	return &pb.UpdatePembeliReq{
		IDPembeli:   req.IDPembeli,
		NamaPembeli: req.NamaPembeli,
		NoTelpon:    req.NoTelpon,
		KodeKota:    req.KodeKota,
		Status:      req.Status,
		UpdateBy:    req.UpdateBy,
	}, nil
}

func decodePembeliResponse(_ context.Context, response interface{}) (interface{}, error) {
	return nil, nil
}

func decodeReadPembeliByKotaRespones(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ReadPembeliByKotaResp)
	return svc.Pembeli{
		IDPembeli:   resp.IDPembeli,
		NamaPembeli: resp.NamaPembeli,
		NoTelpon:    resp.NoTelpon,
		KodeKota:    resp.KodeKota,
		Status:      resp.Status,
		CreateBy:    resp.CreateBy,
	}, nil
}

func decodeReadPembeliResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ReadPembeliResp)
	var rsp svc.Pembelis

	for _, v := range resp.AllPembeli {

		itm := svc.Pembeli{
			IDPembeli:   v.IDPembeli,
			NamaPembeli: v.NamaPembeli,
			NoTelpon:    v.NoTelpon,
			KodeKota:    v.KodeKota,
			Status:      v.Status,
			CreateBy:    v.CreateBy,
		}
		rsp = append(rsp, itm)
	}
	return rsp, nil
}

func makeClientAddPembeliEndpoint(conn *grpc.ClientConn, timeout time.Duration, tracer stdopentracing.Tracer,
	logger log.Logger) endpoint.Endpoint {

	endpoint := grpctransport.NewClient(
		conn,
		grpcName,
		"AddPembeli",
		encodeAddPembeliRequest,
		decodePembeliResponse,
		google_protobuf.Empty{},
		grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
	).Endpoint()

	endpoint = opentracing.TraceClient(tracer, "AddPembeli")(endpoint)
	endpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "AddPembeli",
		Timeout: timeout,
	}))(endpoint)

	return endpoint
}

func makeClientReadPembeliByKotaEndpoint(conn *grpc.ClientConn, timeout time.Duration, tracer stdopentracing.Tracer,
	logger log.Logger) endpoint.Endpoint {

	endpoint := grpctransport.NewClient(
		conn, grpcName,
		"ReadPembeliByKota",
		encodeReadPembeliByKotaRequest,
		decodeReadPembeliByKotaRespones,
		pb.ReadPembeliByKotaResp{},
		grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
	).Endpoint()

	endpoint = opentracing.TraceClient(tracer, "ReadPembeliByKota")(endpoint)
	endpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "ReadPembeliByKota",
		Timeout: timeout,
	}))(endpoint)

	return endpoint
}

func makeClientReadPembeliEndpoint(conn *grpc.ClientConn, timeout time.Duration, tracer stdopentracing.Tracer,
	logger log.Logger) endpoint.Endpoint {

	endpoint := grpctransport.NewClient(
		conn, grpcName,
		"ReadPembeli",
		encodeReadPembeliRequest,
		decodeReadPembeliResponse,
		pb.ReadPembeliResp{},
		grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
	).Endpoint()

	endpoint = opentracing.TraceClient(tracer, "ReadPembeli")(endpoint)
	endpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "ReadPembeli",
		Timeout: timeout,
	}))(endpoint)

	return endpoint
}

func makeClientUpdatePembeli(conn *grpc.ClientConn, timeout time.Duration, tracer stdopentracing.Tracer,
	logger log.Logger) endpoint.Endpoint {
	endpoint := grpctransport.NewClient(
		conn, grpcName,
		"UpdatePembeli",
		encodeUpdatePembeliRequest,
		decodePembeliResponse,
		google_protobuf.Empty{},
		grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
	).Endpoint()

	endpoint = opentracing.TraceClient(tracer, "UpdatePembeli")(endpoint)
	endpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "UpdatePembeli",
		Timeout: timeout,
	}))(endpoint)

	return endpoint
}
