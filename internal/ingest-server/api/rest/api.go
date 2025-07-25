package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	generatedapi "github.com/rom8726/warden/internal/generated/ingestserver"
	"github.com/rom8726/warden/internal/ingest-server/contract"
)

var _ generatedapi.Handler = (*RestAPI)(nil)

type RestAPI struct {
	envelopeUseCase   contract.EnvelopeUseCase
	storeEventUseCase contract.StoreEventUseCase
}

func New(
	envelopeUseCase contract.EnvelopeUseCase,
	storeEventUseCase contract.StoreEventUseCase,
) *RestAPI {
	return &RestAPI{
		envelopeUseCase:   envelopeUseCase,
		storeEventUseCase: storeEventUseCase,
	}
}

func (r *RestAPI) NewError(_ context.Context, err error) *generatedapi.ErrorStatusCode {
	code := http.StatusInternalServerError
	errMessage := err.Error()

	var secError *ogenerrors.SecurityError
	if errors.As(err, &secError) {
		code = http.StatusUnauthorized
		errMessage = "unauthorized"
	}

	return &generatedapi.ErrorStatusCode{
		StatusCode: code,
		Response: generatedapi.Error{
			Error: generatedapi.ErrorError{
				Message: generatedapi.NewOptString(errMessage),
			},
		},
	}
}
