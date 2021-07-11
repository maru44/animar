package usecase

import "animar/v1/internal/pkg/domain"

type AdminInteractor struct {
	animeRepository    AdminAnimeRepository
	platformRepository AdminPlatformRepository
	seasonRepository   AdminSeasonRepository
	seriesRepository   AdminSeriesRepository
}

func NewAdminAnimeInteractor(anime AdminAnimeRepository, platform AdminPlatformRepository, season AdminSeasonRepository, series AdminSeriesRepository) domain.AdminInteractor {
	return &AdminInteractor{
		animeRepository:    anime,
		platformRepository: platform,
		seasonRepository:   season,
		seriesRepository:   series,
	}
}

/************************
        repository
************************/

type AdminAnimeRepository interface {
	ListAll() (domain.TAnimes, error)
	FindById(int) (domain.TAnimeAdmin, error)
	Insert(domain.TAnimeInsert) (int, error)
	Update(int, domain.TAnimeInsert) (int, error)
	Delete(int) (int, error)
}

type AdminPlatformRepository interface {
	ListAll() (domain.TPlatforms, error)
	Insert(domain.TPlatform) (int, error)
	FindById(int) (domain.TPlatform, error)
	Update(domain.TPlatform, int) (int, error)
	Delete(int) (int, error)
	// relation
	FilterByAnime(int) (domain.TRelationPlatforms, error)
	InsertRelation(domain.TRelationPlatformInput) (int, error)
	DeleteRelation(int, int) (int, error)
}

type AdminSeasonRepository interface {
	ListAll() ([]domain.TSeason, error)
	FindById(int) (domain.TSeason, error)
	Insert(domain.TSeasonInput) (int, error)
	InsertRelation(domain.TSeasonRelationInput) (int, error)
}

type AdminSeriesRepository interface {
	ListAll() ([]domain.TSeries, error)
	FindById(int) (domain.TSeries, error)
	Insert(domain.TSeriesInput) (int, error)
	Update(domain.TSeriesInput, int) (int, error)
	Delete(int) (int, error)
}

/************************
         anime
*************************/

func (interactor *AdminInteractor) AnimesAllAdmin() (animes domain.TAnimes, err error) {
	animes, err = interactor.animeRepository.ListAll()
	return
}
func (interactor *AdminInteractor) AnimeDetailAdmin(id int) (anime domain.TAnimeAdmin, err error) {
	anime, err = interactor.animeRepository.FindById(id)
	return
}
func (interactor *AdminInteractor) AnimeInsert(anime domain.TAnimeInsert) (lastInsertId int, err error) {
	lastInsertId, err = interactor.animeRepository.Insert(anime)
	return
}
func (interactor *AdminInteractor) AnimeUpdate(id int, anime domain.TAnimeInsert) (rowsAffected int, err error) {
	rowsAffected, err = interactor.animeRepository.Update(id, anime)
	return
}
func (interactor *AdminInteractor) AnimeDelete(id int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.animeRepository.Delete(id)
	return
}

/************************
          platform
*************************/

func (interactor *AdminInteractor) PlatformAllAdmin() (platforms domain.TPlatforms, err error) {
	platforms, err = interactor.platformRepository.ListAll()
	return
}

func (interactor *AdminInteractor) PlatformDetail(id int) (platform domain.TPlatform, err error) {
	platform, err = interactor.platformRepository.FindById(id)
	return
}

func (interactor *AdminInteractor) PlatformInsert(platform domain.TPlatform) (lastInserted int, err error) {
	lastInserted, err = interactor.platformRepository.Insert(platform)
	return
}

func (interactor *AdminInteractor) PlatformUpdate(platform domain.TPlatform, id int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.platformRepository.Update(platform, id)
	return
}

func (interactor *AdminInteractor) PlatformDelete(id int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.platformRepository.Delete(id)
	return
}

// relation

func (interactor *AdminInteractor) RelationPlatformInsert(platform domain.TRelationPlatformInput) (lastInserted int, err error) {
	lastInserted, err = interactor.platformRepository.InsertRelation(platform)
	return
}

func (interactor *AdminInteractor) RelationPlatformDelete(animeId int, platformId int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.platformRepository.DeleteRelation(animeId, platformId)
	return
}

func (interactor *AdminInteractor) RelationPlatformByAnime(animeId int) (platforms domain.TRelationPlatforms, err error) {
	platforms, err = interactor.platformRepository.FilterByAnime(animeId)
	return
}

// season

func (interactor *AdminInteractor) ListSeason() (seasons []domain.TSeason, err error) {
	seasons, err = interactor.seasonRepository.ListAll()
	return
}

func (interactor *AdminInteractor) DetailSeason(id int) (season domain.TSeason, err error) {
	season, err = interactor.seasonRepository.FindById(id)
	return
}

func (interactor *AdminInteractor) InsertSeason(season domain.TSeasonInput) (lastInserted int, err error) {
	lastInserted, err = interactor.seasonRepository.Insert(season)
	return
}

func (interactor *AdminInteractor) InsertRelationSeasonAnime(rel domain.TSeasonRelationInput) (lastInserted int, err error) {
	lastInserted, err = interactor.seasonRepository.InsertRelation(rel)
	return
}

// series

func (interactor *AdminInteractor) ListSeries() (series []domain.TSeries, err error) {
	series, err = interactor.seriesRepository.ListAll()
	return
}

func (interactor *AdminInteractor) DetailSeries(id int) (series domain.TSeries, err error) {
	series, err = interactor.seriesRepository.FindById(id)
	return
}

func (interactor *AdminInteractor) InsertSeries(series domain.TSeriesInput) (lastInserted int, err error) {
	lastInserted, err = interactor.seriesRepository.Insert(series)
	return
}

func (interactor *AdminInteractor) UpdateSeries(series domain.TSeriesInput, id int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.seriesRepository.Update(series, id)
	return
}

func (interactor *AdminInteractor) DeleteSeries(id int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.seriesRepository.Delete(id)
	return
}
