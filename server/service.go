package server

import "context"

type Status int32

const (
	//ServiceID is dispatch service ID
	ServiceID        = "pembeli.bluebird.id"
	OnAdd     Status = 1
)

type Pembeli struct {
	IDPembeli   int32
	NamaPembeli string
	NoTelpon    string
	KodeKota    string
	Status      int32
	CreateBy    string
	UpdateBy    string
}
type Pembelis []Pembeli

/*type Location struct {
	customerID   int64
	label        []int32
	locationType []int32
	name         []string
	street       string
	village      string
	district     string
	city         string
	province     string
	latitude     float64
	longitude    float64
}*/

type ReadWriter interface {
	AddPembeli(Pembeli) error
	ReadPembeli() (Pembelis, error)
	UpdatePembeli(Pembeli) error
	ReadPembeliByKota(string) (Pembeli, error)
}

type PembeliService interface {
	AddPembeliService(context.Context, Pembeli) error
	ReadPembeliService(context.Context) (Pembelis, error)
	UpdatePembeliService(context.Context, Pembeli) error
	ReadPembeliByKotaService(context.Context, string) (Pembeli, error)
}
