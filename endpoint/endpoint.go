package endpoint

import (
	"context"

	svc "MiniProject/git.bluebird.id/mini/pembeli/server"

	kit "github.com/go-kit/kit/endpoint"
)

type PembeliEndpoint struct {
	AddPembeliEndpoint        kit.Endpoint
	ReadPembeliByKotaEndpoint kit.Endpoint
	ReadPembeliEndpoint       kit.Endpoint
	UpdatePembeliEndpoint     kit.Endpoint
}

func NewPembeliEndpoint(service svc.PembeliService) PembeliEndpoint {
	addPembeliEp := makeAddPembeliEndpoint(service)
	readPembeliByKotaEp := makeReadPembeliByKotaEndpoint(service)
	readPembeliEp := makeReadPembeliEndpoint(service)
	updatePembeliEp := makeUpdatePembeliEndpoint(service)
	return PembeliEndpoint{AddPembeliEndpoint: addPembeliEp,
		ReadPembeliByKotaEndpoint: readPembeliByKotaEp,
		ReadPembeliEndpoint:       readPembeliEp,
		UpdatePembeliEndpoint:     updatePembeliEp,
	}
}

func makeAddPembeliEndpoint(service svc.PembeliService) kit.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(svc.Pembeli)
		err := service.AddPembeliService(ctx, req)
		return nil, err
	}
}

func makeReadPembeliByKotaEndpoint(service svc.PembeliService) kit.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(svc.Pembeli)
		result, err := service.ReadPembeliByKotaService(ctx, req.KodeKota)
		return result, err
	}
}

func makeReadPembeliEndpoint(service svc.PembeliService) kit.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		result, err := service.ReadPembeliService(ctx)
		return result, err
	}
}

func makeUpdatePembeliEndpoint(service svc.PembeliService) kit.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(svc.Pembeli)
		err := service.UpdatePembeliService(ctx, req)
		return nil, err
	}
}
