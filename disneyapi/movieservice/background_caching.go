package movieservice

import (
	"log/slog"
	"time"

	"github.com/movieofthenight/go-streaming-availability/v4"
)

type BackgroundCachingService struct {
	dataAccess   *MovieDataAccess
	movieService *MovieService
}

type CachingService interface {
	Start()
}

func CreateBackgroundCachingService(dataAccess *MovieDataAccess, service *MovieService) CachingService {
	return &BackgroundCachingService{
		dataAccess:   dataAccess,
		movieService: service,
	}
}

func (s *BackgroundCachingService) Start() {
	go s.loop()
}

func (s *BackgroundCachingService) loop() {
	for {
		time.Sleep(time.Hour * 24)
		slog.Info("Starting background loop")
		movies, err := (*(s.movieService)).GetMovies()
		if err != nil {
			slog.Error("Error getting movies from movie service")
			time.Sleep(time.Hour * 24)
			continue
		}
		slog.Info("Successfully got movies from movie service")
		dbMovies := s.mapMovies(movies)
		slog.Info("Successfully mapped movies")
		err = (*(s.dataAccess)).InsertMovies(dbMovies)
		if err != nil {
		}
		slog.Info("Successfully inserted movies")
	}
}

func (s *BackgroundCachingService) mapMovies(shows *[]streaming.Show) *[]MovieDBObject {
	results := make([]MovieDBObject, len(*shows))
	for i, v := range *shows {
		newShow := MovieDBObject{
			title:                 v.Title,
			overview:              v.Overview,
			horizontalPosterw1080: v.ImageSet.HorizontalPoster.W1080,
			verticalPosterw720:    v.ImageSet.VerticalPoster.W720,
		}
		results[i] = newShow
	}
	return &results
}
