package main

import (
	"context"
	"fmt"
	"time"

	cli "MiniProject/git.bluebird.id/mini/pembeli/endpoint"
	opt "MiniProject/git.bluebird.id/mini/util/grpc"
	util "MiniProject/git.bluebird.id/mini/util/microservice"
	tr "github.com/opentracing/opentracing-go"
)

func main() {
	logger := util.Logger()
	tracer := tr.GlobalTracer()
	option := opt.ClientOption{Retry: 3, RetryTimeout: 500 * time.Millisecond, Timeout: 30 * time.Second}

	client, err := cli.NewGRPCPembeliClient([]string{"127.0.0.1:2181"}, nil, option, tracer, logger)
	if err != nil {
		logger.Log("error", err)
	}

	//Add Customer
	//client.AddPembeliService(context.Background(), svc.Pembeli{IDPembeli: 22, NamaPembeli: "jajat", NoTelpon: "087669", KodeKota: "jakarta", Status: 1, CreateBy: "tas"})

	//Get Customer By Mobile No
	//cusMobile, _ := client.ReadPembeliByKotaService(context.Background(), "Bandung")
	//fmt.Println("pembeli based on nama:", cusMobile)

	//List Customer
	cuss, _ := client.ReadPembeliService(context.Background())
	fmt.Println("all pembeli:", cuss)

	//Update Customer
	//client.UpdatePembeliService(context.Background(), svc.Pembeli{IDPembeli: 22, NamaPembeli: "Drajajat", NoTelpon: "087669", KodeKota: "Bandung", Status: 0, CreateBy: "tas"})

	//Get Customer By Email
	//cuss, _ := client.ReadPembeliService(context.Background())
	//fmt.Println("all pembelis:", cuss)
}
