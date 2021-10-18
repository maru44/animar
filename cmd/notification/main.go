package notification

import (
	"animar/v1/internal/pkg/infrastructure"
	"animar/v1/internal/pkg/interfaces/batch"
	"log"
)

func main() {
	sqlHandler := infrastructure.NewSqlHandler()

	platformBatch := batch.NewPlatformBatch(sqlHandler)
	if err := platformBatch.SendBatch(); err != nil {
		log.Println(err)
	}
}
