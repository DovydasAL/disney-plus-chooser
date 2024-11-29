package main

import (
	"net/http"

	movieservice "github.com/DovydasAL/disneyapi/movieservice"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Howdy"))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response, err := movieservice.GetSuggestion(nil)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		render.Status(r, http.StatusOK)
		render.Render(w, r, NewSuggestionResponse(response))
	})

	r.Post("/api/v1/get_suggestion", func(w http.ResponseWriter, r *http.Request) {
		data := &movieservice.SuggestionRequestV1{}
		err := render.Bind(r, data)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		response, err := movieservice.GetSuggestion(data)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		render.Status(r, http.StatusOK)
		render.Render(w, r, NewSuggestionResponse(response))
	})
	http.ListenAndServe(":3000", r)
}

func NewSuggestionResponse(suggestion *movieservice.Suggestion) *SuggestionResponseV1 {
	resp := &SuggestionResponseV1{Suggestion: suggestion}
	return resp
}

type SuggestionResponseV1 struct {
	Suggestion *movieservice.Suggestion `json:"suggestion"`
}

func (s *SuggestionResponseV1) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrInternalServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error.",
		ErrorText:      err.Error(),
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render implements render.Renderer.
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}