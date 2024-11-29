package movieservice

import (
	"errors"
	"net/http"
)

type Suggestion struct {
	Movie string `json:"movie"`
}

type SuggestionRequestV1 struct {
	TestParam string `json:"test_param"`
}

func (s *SuggestionRequestV1) Bind(r *http.Request) error {
	if s.TestParam == "" {
		return errors.New("TestParam cannot be nil")
	}
	return nil
}
