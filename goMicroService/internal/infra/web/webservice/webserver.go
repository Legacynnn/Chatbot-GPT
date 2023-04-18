package webservice

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebService struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebService(router chi.Router, handlers map[string]http.HandlerFunc, webServerPort string) *WebService {
	return &WebService{
		WebServerPort: webServerPort,
	}
}

func (s *WebService) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebService) Start() {
	s.Router.Use(middleware.Logger)
	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	if err := http.ListenAndServe(s.WebServerPort, s.Router); err != nil {
		panic(err.Error())
	}
}
