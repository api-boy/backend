package httputils

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"apiboy/backend/src/errors"

	kithttp "github.com/go-kit/kit/transport/http"
)

// DecodeRPCRequest decodes an http request into an input type
func DecodeRPCRequest(inPtr interface{}) kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		// create the input variable from the given type
		inputType := reflect.TypeOf(inPtr).Elem()
		input := reflect.New(inputType).Interface()

		// parse request body
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(input); err != nil {
			return nil, errors.BadRequest{}
		}
		defer r.Body.Close()

		return input, nil
	}
}
