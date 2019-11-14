package service

import (
	"context"
	"net/http"

	"apiboy/backend/src/config"
	"apiboy/backend/src/firebase"
	"apiboy/backend/src/logger"
	"apiboy/backend/src/store"

	"firebase.google.com/go/auth"
	"github.com/facebookgo/grace/gracehttp"
)

// Service implements the service logic
type Service struct {
	Config             *config.Config
	Logger             *logger.Logger
	Store              *store.Store
	FirebaseAuthClient *auth.Client
}

// New returns a new Service
func New(conf *config.Config) (*Service, error) {
	ctx := context.Background()

	firebaseApp, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, err
	}

	firestoreClient, err := firebaseApp.NewFirestoreClient(ctx)
	if err != nil {
		return nil, err
	}

	firebaseAuthClient, err := firebaseApp.NewAuthClient(ctx)
	if err != nil {
		return nil, err
	}

	log := logger.New(conf)
	st := store.New(conf, firestoreClient)

	return &Service{
		Config:             conf,
		Logger:             log,
		Store:              st,
		FirebaseAuthClient: firebaseAuthClient,
	}, nil
}

// Run executes the service
func (s *Service) Run() {
	ctx := context.Background()
	addr := ":" + s.Config.Port

	httpEndpoints := MakeHTTPEndpoints(s)
	httpHandler := MakeHTTPHandler(ctx, s.Logger, httpEndpoints)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: httpHandler,
	}

	if err := gracehttp.Serve(httpServer); err != nil {
		s.Logger.Error("listening error", logger.Field{Key: "err", Val: err})
	}
}
