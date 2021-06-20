package usecase

import "animar/v1/pkg/domain"

type AdminAnimeInteractor struct {
	AdminAnimeRepository AdminAnimeRepository
}

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
