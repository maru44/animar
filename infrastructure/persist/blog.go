package persist

import (
	"animar/v1/domain"
	"animar/v1/repository"
	"database/sql"
)

type blogPersist struct {}

func NewBlogPersist repository.blogPersist