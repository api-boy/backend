package enums

const (
	// RequestTypeGet is the type GET
	RequestTypeGet = "GET"

	// RequestTypePost is the type POST
	RequestTypePost = "POST"

	// RequestTypeDelete is the type DELETE
	RequestTypeDelete = "DELETE"

	// RequestTypeHead is the type HEAD
	RequestTypeHead = "HEAD"

	// RequestTypeOptions is the type OPTIONS
	RequestTypeOptions = "OPTIONS"

	// RequestTypePut is the type PUT
	RequestTypePut = "PUT"

	// RequestTypePatch is the type PATCH
	RequestTypePatch = "PATCH"
)

// IsValidRequestType return valid type request
func IsValidRequestType(requestType string) bool {
	if requestType == RequestTypeGet || requestType == RequestTypePost ||
		requestType == RequestTypeDelete || requestType == RequestTypeHead ||
		requestType == RequestTypeOptions || requestType == RequestTypePut ||
		requestType == RequestTypePatch {
		return true
	}

	return false
}
