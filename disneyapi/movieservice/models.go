package movieservice

import (
	"net/http"
)

type Suggestion struct {
	AllMovies *[]MovieDBObject `json:"all_movies"`
	Movie     string           `json:"movie"`
}

type SuggestionRequestV1 struct {
}

func (s *SuggestionRequestV1) Bind(r *http.Request) error {
	// Validation on request object here
	return nil
}
