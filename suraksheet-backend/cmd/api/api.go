package api

import (
	"database/sql"
	"net/http"

	"github.com/LikheKeto/Suraksheet/service/bin"
	"github.com/LikheKeto/Suraksheet/service/document"
	"github.com/LikheKeto/Suraksheet/service/user"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

type APIServer struct {
	addr     string
	db       *sql.DB
	minio    *minio.Client
	rmqChan  *amqp.Channel
	rmq      amqp.Queue
	esClient *elasticsearch.Client
}

func NewAPIServer(addr string, db *sql.DB, minio *minio.Client, rmqChan *amqp.Channel, rmq amqp.Queue, esClient *elasticsearch.Client) *APIServer {
	return &APIServer{
		addr:     addr,
		db:       db,
		minio:    minio,
		rmqChan:  rmqChan,
		rmq:      rmq,
		esClient: esClient,
	}
}

func (s *APIServer) Run() error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	subrouter := chi.NewRouter()
	router.Mount("/api/v1", subrouter)

	userStore := user.NewStore(s.db)
	binStore := bin.NewStore(s.db)
	documentStore := document.NewStore(s.db)

	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	binHandler := bin.NewHandler(binStore, userStore, documentStore, s.minio)
	binHandler.RegisterRoutes(subrouter)

	documentHandler := document.NewHandler(documentStore, userStore, binStore, s.minio, s.rmqChan, s.rmq, s.esClient)
	documentHandler.RegisterRoutes(subrouter)

	err := http.ListenAndServe(s.addr, router)
	return err
}
