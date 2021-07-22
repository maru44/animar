package main

import (
	"animar/v1/internal/pkg/batch/dbbackup"
)

func main() {
	dbbackup.Test()
	dbbackup.BackupMainDatabase()
}
