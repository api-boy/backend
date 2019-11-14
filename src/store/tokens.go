package store

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// TokensCollection is the name of the collection
const TokensCollection = "tokens"

// Token represents a model in the database
type Token struct {
	ID      string `json:"id" firestore:"id"`
	UserID  string `json:"user_id" firestore:"user_id"`
	Created *Event `json:"created" firestore:"created"`
}

// NewTokenID generates a UUID for tokens
func (s *Store) NewTokenID() string {
	return "tok-" + uuid.New().String()
}

// CreateToken creates a new token
func (s *Store) CreateToken(ctx context.Context, token *Token) error {
	token.Created = NewEvent(token.UserID)
	_, err := s.Client.Collection(TokensCollection).Doc(token.ID).Set(ctx, token)
	return err
}

// DeleteToken deletes a token
func (s *Store) DeleteToken(ctx context.Context, id string) error {
	_, err := s.Client.Collection(TokensCollection).Doc(id).Delete(ctx)
	return err
}

// GetTokenByID gets a token by id
func (s *Store) GetTokenByID(ctx context.Context, id string) (*Token, error) {
	iter := s.Client.Collection(TokensCollection).Where("id", "==", id).Limit(1).Documents(ctx)

	snapshot, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	token := &Token{}
	snapshot.DataTo(token)

	return token, nil
}
