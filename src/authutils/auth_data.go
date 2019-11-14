package authutils

// AuthData contains the authentication data that is sent to the endpoints inside the context
type AuthData struct {
	JwtID     string
	UserID    string
	UserName  string
	UserEmail string
	UserRole  string
}
