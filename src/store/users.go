package store

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// UsersCollection is the name of the collection
const UsersCollection = "users"

// User represents a model in the database
type User struct {
	ID       string `json:"id" firestore:"id"`
	Name     string `json:"name" firestore:"name"`
	Email    string `json:"email" firestore:"email"`
	Password string `json:"-" firestore:"password"`
	Role     string `json:"role" firestore:"role"`
	TempCode string `json:"-" firestore:"temp_code"`
	Created  *Event `json:"created" firestore:"created"`
	Updated  *Event `json:"updated" firestore:"updated"`
	Deleted  *Event `json:"deleted" firestore:"deleted"`
}

// NewUserID generates a UUID for users
func (s *Store) NewUserID() string {
	return "usr-" + uuid.New().String()
}

// CreateUser creates a new user
func (s *Store) CreateUser(ctx context.Context, userID string, user *User) error {
	user.Created = NewEvent(userID)
	_, err := s.Client.Collection(UsersCollection).Doc(user.ID).Set(ctx, user)
	return err
}

// UpdateUser updates an existing user
func (s *Store) UpdateUser(ctx context.Context, userID string, user *User) error {
	user.Updated = NewEvent(userID)
	_, err := s.Client.Collection(UsersCollection).Doc(user.ID).Set(ctx, user)
	return err
}

// DeleteUser deletes an existing user
func (s *Store) DeleteUser(ctx context.Context, userID string, user *User) error {
	user.Deleted = NewEvent(userID)
	_, err := s.Client.Collection(UsersCollection).Doc(user.ID).Set(ctx, user)
	return err
}

// GetUserByID gets a user by id
func (s *Store) GetUserByID(ctx context.Context, id string) (*User, error) {
	return s.getUserByField(ctx, "id", id)
}

// GetUserByEmail gets a user by email
func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.getUserByField(ctx, "email", email)
}

// getUserByField gets a user by a given field
func (s *Store) getUserByField(ctx context.Context, field, value string) (*User, error) {
	iter := s.Client.Collection(UsersCollection).Where(field, "==", value).Limit(1).Documents(ctx)

	snapshot, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	user := &User{}
	snapshot.DataTo(user)

	if user.Deleted != nil {
		return nil, nil
	}

	return user, nil
}
