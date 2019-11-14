package firebase

import (
	"context"
	"fmt"

	"apiboy/backend/src/config"

	"cloud.google.com/go/firestore"
	fb "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

// App wraps a Firebase application
type App struct {
	FirebaseApp *fb.App
}

// NewApp returns a new App
func NewApp(ctx context.Context, conf *config.Config) (*App, error) {
	serviceAccountFile := fmt.Sprintf(".team/%s/firebase-service-account.json", conf.UpStage)
	clientOption := option.WithCredentialsFile(serviceAccountFile)

	app, err := fb.NewApp(ctx, &fb.Config{}, clientOption)
	if err != nil {
		return nil, err
	}

	return &App{
		FirebaseApp: app,
	}, nil
}

// NewFirestoreClient returns a Firestore client
func (a *App) NewFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	client, err := a.FirebaseApp.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewAuthClient returns a Firebase Auth client
func (a *App) NewAuthClient(ctx context.Context) (*auth.Client, error) {
	client, err := a.FirebaseApp.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
