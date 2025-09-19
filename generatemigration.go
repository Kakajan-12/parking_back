package main

import (
	"backend/contrib/models"
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	// Use postgres dialect instead of sqlite
	stmts, err := gormschema.New("postgres").Load(
		&models.User{},
		&models.UserSession{},
		&models.MacUser{},
		&models.Camera{},
		&models.Car{},
		&models.CarSubscription{},
		&models.CarSession{},
		&models.Tariff{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	// Write schema to stdout (Atlas will consume this)
	io.WriteString(os.Stdout, stmts)
}
