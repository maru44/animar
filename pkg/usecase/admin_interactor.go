package usecase

import "animar/v1/pkg/domain"

type AdminAnimeInteractor struct {
	AdminAnimeRepository AdminAnimeRepository
}

type AdminPlatformInteractor struct {
	AdminPlatformRepository AdminPlatformRepository
}

/************************
         anime
*************************/

func (interactor *AdminAnimeInteractor) AnimesAllAdmin() (animes domain.TAnimes, err error) {
	animes, err = interactor.AdminAnimeRepository.ListAll()
	return
}
func (interactor *AdminAnimeInteractor) AnimeDetailAdmin(id int) (anime domain.TAnimeAdmin, err error) {
	anime, err = interactor.AdminAnimeRepository.FindById(id)
	return
}
func (interactor *AdminAnimeInteractor) AnimeInsert(anime domain.TAnimeInsert) (lastInsertId int, err error) {
	lastInsertId, err = interactor.AdminAnimeRepository.Insert(anime)
	return
}
func (interactor *AdminAnimeInteractor) AnimeUpdate(id int, anime domain.TAnimeInsert) (rowsAffected int, err error) {
	rowsAffected, err = interactor.AdminAnimeRepository.Update(id, anime)
	return
}
func (interactor *AdminAnimeInteractor) AnimeDelete(id int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.AdminAnimeRepository.Delete(id)
	return
}

/************************
          platform
*************************/

func (interactor *AdminPlatformInteractor) PlatformAllAdmin() (platforms domain.TPlatforms, err error) {
	platforms, err = interactor.AdminPlatformRepository.ListAll()
	return
}

func (interactor *AdminPlatformInteractor) PlatformDetail(id int) (platform domain.TPlatform, err error) {
	platform, err = interactor.AdminPlatformRepository.FindById(id)
	return
}

func (interactor *AdminPlatformInteractor) PlatformInsert(platform domain.TPlatform) (lastInserted int, err error) {
	lastInserted, err = interactor.AdminPlatformRepository.Insert(platform)
	return
}

func (interactor *AdminPlatformInteractor) PlatformUpdate(platform domain.TPlatform, id int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.AdminPlatformRepository.Update(platform, id)
	return
}

func (interactor *AdminPlatformInteractor) PlatformDelete(id int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.AdminPlatformRepository.Delete(id)
	return
}

// relation

func (interactor *AdminPlatformInteractor) RelationPlatformInsert(platform domain.TRelationPlatformInput) (lastInserted int, err error) {
	lastInserted, err = interactor.AdminPlatformRepository.InsertRelation(platform)
	return
}

func (interactor *AdminPlatformInteractor) RelationPlatformDelete(animeId int, platformId int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.AdminPlatformRepository.DeleteRelation(animeId, platformId)
	return
}

func (interactor *AdminPlatformInteractor) RelationPlatformByAnime(animeId int) (platforms domain.TRelationPlatforms, err error) {
	platforms, err = interactor.AdminPlatformRepository.FilterByAnime(animeId)
	return
}
