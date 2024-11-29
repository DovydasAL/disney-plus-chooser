package movieservice

import (
	"context"

	"math/rand"

	"github.com/movieofthenight/go-streaming-availability/v4"
)

const apiKey = "98249cbed4msh3601788fa047ff9p195816jsn601ab369576a"

var client *streaming.APIClient = nil
var movies *[]streaming.Show = nil

func getClient() *streaming.APIClient {
	if client == nil {
		client = streaming.NewAPIClientFromRapidAPIKey(apiKey, nil)
	}
	return client
}

func getMovies() (*[]streaming.Show, error) {
	if movies != nil {
		return movies, nil
	}
	results := []streaming.Show{}
	apiResult, err := getClient().ShowsAPI.SearchShowsByFilters(context.Background()).Country("us").Catalogs([]string{"disney"}).ShowType("movie").ExecuteWithAutoPagination(1)
	for apiResult.Next() {
		movie := apiResult.Get()
		results = append(results, movie)
	}
	if err != nil {
		return nil, err
	}
	movies = &results
	return movies, nil
}

func GetSuggestion(suggestionRequest *SuggestionRequestV1) (*Suggestion, error) {
	suggestion := &Suggestion{}
	movies, err := getMovies()
	if err != nil {
		return nil, err
	}
	movieTitles := make([]string, len(*movies))
	for i := 0; i < len(movieTitles); i++ {
		movieTitles[i] = (*movies)[i].Title
	}
	suggestion.Movie = movieTitles[rand.Intn(len(movieTitles))]
	return suggestion, nil
}
