package main

import (
	"context"
	"flag"
	"log"

	"github.com/ebisaan/inventory/internal/adapter/postgres"
)

func main() {
	dsn := flag.String("dsn", "", "Data connection string")
	flag.Parse()

	db, err := postgres.NewAdapter(*dsn)
	if err != nil {
		log.Println("failed to create postgres adapter" + err.Error())
	}

	log.Println("starting auto migrate on inventory database...")
	err = db.AutoMigration(context.Background())
	if err != nil {
		log.Println("failed to migrate" + err.Error())
	}
}
