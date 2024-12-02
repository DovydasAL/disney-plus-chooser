package movieservice

import (
	"context"

	"math/rand"

	"github.com/movieofthenight/go-streaming-availability/v4"
)

type DisneyMovieService struct {
	client     *streaming.APIClient
	dataAccess *MovieDataAccess
}

type MovieService interface {
	GetSuggestion(*SuggestionRequestV1) (*Suggestion, error)
	GetMovies() (*[]streaming.Show, error)
}

func CreateMovieService(apiKey string, dataAccess *MovieDataAccess) MovieService {
	return &DisneyMovieService{
		client:     streaming.NewAPIClientFromRapidAPIKey(apiKey, nil),
		dataAccess: dataAccess,
	}
}

func (service *DisneyMovieService) GetMovies() (*[]streaming.Show, error) {
	results := []streaming.Show{}
	apiResult, err := service.client.ShowsAPI.SearchShowsByFilters(context.Background()).Country("us").Catalogs([]string{"disney"}).ShowType("movie").ExecuteWithAutoPagination(1)
	if err != nil {
		return nil, err
	}
	for apiResult.Next() {
		movie := apiResult.Get()
		results = append(results, movie)
	}
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func (service *DisneyMovieService) GetSuggestion(suggestionRequest *SuggestionRequestV1) (*Suggestion, error) {
	suggestion := &Suggestion{}
	movies, err := (*(service.dataAccess)).GetMovies()
	if err != nil {
		return nil, err
	}
	movieTitles := make([]string, len(*movies))
	for i := 0; i < len(movieTitles); i++ {
		movieTitles[i] = (*movies)[i].title
	}
	suggestion.Movie = movieTitles[rand.Intn(len(movieTitles))]
	suggestion.AllMovies = movies
	return suggestion, nil
}
