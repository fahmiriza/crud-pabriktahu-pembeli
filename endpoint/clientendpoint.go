package endpoint

import (
	"context"
	"fmt"

	sv "MiniProject/git.bluebird.id/mini/pembeli/server"
)

func (ke PembeliEndpoint) AddPembeliService(ctx context.Context, pembeli sv.Pembeli) error {
	_, err := ke.AddPembeliEndpoint(ctx, pembeli)
	return err
}

func (ke PembeliEndpoint) ReadPembeliByKotaService(ctx context.Context, kodekota string) (sv.Pembeli, error) {
	req := sv.Pembeli{KodeKota: kodekota}
	fmt.Println(req)
	resp, err := ke.ReadPembeliByKotaEndpoint(ctx, req)
	if err != nil {
		fmt.Println("error pada endpoint: ", err)
	}
	cus := resp.(sv.Pembeli)
	return cus, err
}

func (ke PembeliEndpoint) ReadPembeliService(ctx context.Context) (sv.Pembelis, error) {
	resp, err := ke.ReadPembeliEndpoint(ctx, nil)
	fmt.Println("ke resp", resp)
	if err != nil {
		fmt.Println("error pada endpoint: ", err)
	}
	return resp.(sv.Pembelis), err
}

func (ke PembeliEndpoint) UpdatePembeliService(ctx context.Context, kar sv.Pembeli) error {
	_, err := ke.UpdatePembeliEndpoint(ctx, kar)
	if err != nil {
		fmt.Println("error pada endpoint:", err)
	}
	return err
}
