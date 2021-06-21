package usecase

import "animar/v1/pkg/domain"

type PlatformRepository interface {
	ListAll() (domain.TPlatforms, err error)

}