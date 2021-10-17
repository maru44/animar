package main

import (
	"animar/v1/internal/pkg/infrastructure"
	"animar/v1/internal/pkg/interfaces/batch"
	"log"
)

func main() {
	// dbbackup.Test()
	sqlHandler := infrastructure.NewSqlHandler()

	platformBatch := batch.NewPlatformBatch(sqlHandler)

	if err := platformBatch.SendBatch(); err != nil {
		log.Println(err)
	}

	// @TODO remove comment out
	// dbbackup.BackupMainDatabase()
}
