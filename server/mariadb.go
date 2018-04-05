package server

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	addPembeli          = `insert into pembeli(IDPembeli,NamaPembeli,NoTelpon,KodeKota,Status,CreateBy,CreateOn)values (?,?,?,?,?,?,?)`
	selectPembeli       = `select IDPembeli,NamaPembeli,NoTelpon,KodeKota,Status,CreateBy from pembeli `
	updatePembeli       = `update pembeli set NamaPembeli=?,NoTelpon=?,KodeKota=?,Status=?,UpdateOn=?, UpdateBy=? where IDPembeli=?`
	selectPembeliByKota = `select IDPembeli,NamaPembeli,NoTelpon,KodeKota,Status,CreateOn from pembeli where KodeKota=?`
)

type dbReadWriter struct {
	db *sql.DB
}

func NewDBReadWriter(url string, schema string, user string, password string) ReadWriter {
	schemaURL := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, url, schema)
	db, err := sql.Open("mysql", schemaURL)
	if err != nil {
		panic(err)
	}
	return &dbReadWriter{db: db}
}

func (rw *dbReadWriter) AddPembeli(pembeli Pembeli) error {
	fmt.Println("insert")
	tx, err := rw.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(addPembeli, pembeli.IDPembeli, pembeli.NamaPembeli, pembeli.NoTelpon, pembeli.KodeKota, OnAdd, pembeli.CreateBy, time.Now())
	//fmt.Println(err)
	if err != nil {
		return err

	}
	return tx.Commit()
}

func (rw *dbReadWriter) ReadPembeliByKota(kodekota string) (Pembeli, error) {
	fmt.Println("show by nama")
	pembeli := Pembeli{KodeKota: kodekota}
	err := rw.db.QueryRow(selectPembeliByKota, kodekota).Scan(&pembeli.IDPembeli, &pembeli.NamaPembeli, &pembeli.NoTelpon, &pembeli.KodeKota,
		&pembeli.Status, &pembeli.CreateBy)

	if err != nil {
		return Pembeli{}, err
	}

	return pembeli, nil
}

func (rw *dbReadWriter) ReadPembeli() (Pembelis, error) {
	fmt.Println("show all")
	pembeli := Pembelis{}
	rows, _ := rw.db.Query(selectPembeli)
	defer rows.Close()
	for rows.Next() {
		var k Pembeli
		err := rows.Scan(&k.IDPembeli, &k.NamaPembeli, &k.NoTelpon, &k.KodeKota, &k.Status, &k.CreateBy)
		if err != nil {
			fmt.Println("error query:", err)
			return pembeli, err
		}
		pembeli = append(pembeli, k)
	}
	//fmt.Println("db nya:", mahasiswa)
	return pembeli, nil
}

func (rw *dbReadWriter) UpdatePembeli(p Pembeli) error {
	fmt.Println("update successfuly")
	tx, err := rw.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(updatePembeli, p.NamaPembeli, p.NoTelpon, p.KodeKota, p.Status, p.UpdateBy, time.Now(), p.IDPembeli)

	fmt.Println("name:", p.NamaPembeli, p.Status)
	if err != nil {
		return err
	}

	return tx.Commit()
}
