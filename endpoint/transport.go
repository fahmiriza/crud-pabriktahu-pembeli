package endpoint

import (
	"context"

	scv "MiniProject/git.bluebird.id/mini/pembeli/server"

	pb "MiniProject/git.bluebird.id/mini/pembeli/grpc"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	stdopentracing "github.com/opentracing/opentracing-go"
	oldcontext "golang.org/x/net/context"
)

type grpcPembeliServer struct {
	addPembeli        grpctransport.Handler
	readPembeliByKota grpctransport.Handler
	readPembeli       grpctransport.Handler
	updatePembeli     grpctransport.Handler
}

func NewGRPCPembeliServer(endpoints PembeliEndpoint, tracer stdopentracing.Tracer,
	logger log.Logger) pb.PembeliServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcPembeliServer{
		addPembeli: grpctransport.NewServer(endpoints.AddPembeliEndpoint,
			decodeAddPembeliRequest,
			encodeEmptyResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "AddPembeli", logger)))...),
		readPembeliByKota: grpctransport.NewServer(endpoints.ReadPembeliByKotaEndpoint,
			decodeReadPembeliByKotaRequest,
			encodeReadPembeliByKotaResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "ReadPembeliByKota", logger)))...),
		readPembeli: grpctransport.NewServer(endpoints.ReadPembeliEndpoint,
			decodeReadPembeliRequest,
			encodeReadPembeliResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "ReadPembeli", logger)))...),
		updatePembeli: grpctransport.NewServer(endpoints.UpdatePembeliEndpoint,
			decodeUpdatePembeliRequest,
			encodeEmptyResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "UpdatePembeli", logger)))...),
	}
}

func decodeAddPembeliRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.AddPembeliReq)
	return scv.Pembeli{IDPembeli: req.GetIDPembeli(), NamaPembeli: req.GetNamaPembeli(), NoTelpon: req.GetNoTelpon(), KodeKota: req.GetKodeKota(),
		Status: req.GetStatus(), CreateBy: req.GetCreateBy()}, nil
}

func encodeEmptyResponse(_ context.Context, response interface{}) (interface{}, error) {
	return &google_protobuf.Empty{}, nil
}

func (s *grpcPembeliServer) AddPembeli(ctx oldcontext.Context, pembeli *pb.AddPembeliReq) (*google_protobuf.Empty, error) {
	_, resp, err := s.addPembeli.ServeGRPC(ctx, pembeli)
	if err != nil {
		return nil, err
	}
	return resp.(*google_protobuf.Empty), nil
}

func decodeReadPembeliByKotaRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ReadPembeliByKotaReq)
	return scv.Pembeli{KodeKota: req.KodeKota}, nil
}

func decodeReadPembeliRequest(_ context.Context, request interface{}) (interface{}, error) {
	return nil, nil
}

func encodeReadPembeliByKotaResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(scv.Pembeli)
	return &pb.ReadPembeliByKotaResp{IDPembeli: resp.IDPembeli, NamaPembeli: resp.NamaPembeli, NoTelpon: resp.NoTelpon, KodeKota: resp.KodeKota,
		Status: resp.Status, CreateBy: resp.CreateBy}, nil
}

func encodeReadPembeliResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(scv.Pembelis)

	rsp := &pb.ReadPembeliResp{}

	for _, v := range resp {
		itm := &pb.ReadPembeliByKotaResp{
			IDPembeli:   v.IDPembeli,
			NamaPembeli: v.NamaPembeli,
			NoTelpon:    v.NoTelpon,
			KodeKota:    v.KodeKota,
			Status:      v.Status,
			CreateBy:    v.CreateBy,
		}
		rsp.AllPembeli = append(rsp.AllPembeli, itm)
	}
	return rsp, nil
}

func (s *grpcPembeliServer) ReadPembeliByKota(ctx oldcontext.Context, kodekota *pb.ReadPembeliByKotaReq) (*pb.ReadPembeliByKotaResp, error) {
	_, resp, err := s.readPembeliByKota.ServeGRPC(ctx, kodekota)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ReadPembeliByKotaResp), nil
}

func (s *grpcPembeliServer) ReadPembeli(ctx oldcontext.Context, e *google_protobuf.Empty) (*pb.ReadPembeliResp, error) {
	_, resp, err := s.readPembeli.ServeGRPC(ctx, e)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ReadPembeliResp), nil
}

func decodeUpdatePembeliRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UpdatePembeliReq)
	return scv.Pembeli{IDPembeli: req.IDPembeli, NamaPembeli: req.NamaPembeli, NoTelpon: req.NoTelpon, KodeKota: req.KodeKota, Status: req.Status, UpdateBy: req.UpdateBy}, nil
}

func (s *grpcPembeliServer) UpdatePembeli(ctx oldcontext.Context, cus *pb.UpdatePembeliReq) (*google_protobuf.Empty, error) {
	_, resp, err := s.updatePembeli.ServeGRPC(ctx, cus)
	if err != nil {
		return &google_protobuf.Empty{}, err
	}
	return resp.(*google_protobuf.Empty), nil
}
