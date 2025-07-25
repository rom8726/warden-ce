package rest

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/stretchr/testify/assert"
)

func TestRestAPI_NewError(t *testing.T) {
	r := &RestAPI{}

	t.Run("internal error", func(t *testing.T) {
		err := errors.New("some error")
		resp := r.NewError(nil, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "some error", resp.Response.Error.Message.Value)
	})

	t.Run("security error", func(t *testing.T) {
		err := &ogenerrors.SecurityError{}
		resp := r.NewError(nil, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Equal(t, "unauthorized", resp.Response.Error.Message.Value)
	})
}
