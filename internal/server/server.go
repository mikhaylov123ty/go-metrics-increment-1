package server

import (
	"log"
	"net/http"

	"metrics/internal/server/api"
	"metrics/internal/storage"

	"github.com/go-chi/chi/v5"
)

// Структура сервера
type Server struct {
	storage *storage.Storage
}

// Конструктор инстанса сервера
func New(storage *storage.Storage) *Server {
	return &Server{storage: storage}
}

// Метод запуска сервера
func (s *Server) Start(address string) {
	// Создание роутера
	router := chi.NewRouter()

	// Назначение соответствий хендлеров
	s.addHandlers(router, api.NewHandler(*s.storage))

	// Старт сервера
	log.Printf("Starting server on port %s", address)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatal(err)
	}
}

// Наполнение сервера методами хендлера
func (s *Server) addHandlers(router *chi.Mux, handler *api.Handler) {
	// /update
	router.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", handler.UpdatePost)
	})

	// /value
	router.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handler.ValueGet)
	})

	// index
	router.Get("/", handler.IndexGet)
}
