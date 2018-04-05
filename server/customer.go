package server

import (
	"context"
)

type pembeli struct {
	writer ReadWriter
}

func NewPembeli(writer ReadWriter) PembeliService {
	return &pembeli{writer: writer}
}

//Methode pada interface CustomerService di service.go
func (c *pembeli) AddPembeliService(ctx context.Context, pembeli Pembeli) error {
	//fmt.Println("pembeli")
	err := c.writer.AddPembeli(pembeli)
	if err != nil {
		return err
	}

	return nil
}

func (c *pembeli) ReadPembeliService(ctx context.Context) (Pembelis, error) {
	cus, err := c.writer.ReadPembeli()
	//fmt.Println("customer", cus)
	if err != nil {
		return cus, err
	}
	return cus, nil
}

func (c *pembeli) UpdatePembeliService(ctx context.Context, cus Pembeli) error {
	err := c.writer.UpdatePembeli(cus)
	if err != nil {
		return err
	}
	return nil
}

func (c *pembeli) ReadPembeliByKotaService(ctx context.Context, KodeKota string) (Pembeli, error) {
	cus, err := c.writer.ReadPembeliByKota(KodeKota)
	//fmt.Println("customer:", cus)
	if err != nil {
		return cus, err
	}
	return cus, nil
}
