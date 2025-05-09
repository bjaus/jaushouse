package errs

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"

	"pkg/validate"
	"pkg/web"
)

type Response struct {
	Error  string            `json:"error"`
	Detail map[string]string `json:"detail,omitempty"`
}

func Respond(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		web.Respond(w, http.StatusRequestTimeout, Response{Error: "request canceled or deadline exceeded"})
		return
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		payload := Response{
			Error: "invalid request payload",
			Detail: lo.Associate(ve, func(fe validator.FieldError) (string, string) {
				return fe.Field(), fe.Translate(validate.Translator)
			}),
		}
		web.Respond(w, http.StatusUnprocessableEntity, payload)
		return
	}

	e := As(err)
	payload := Response{
		Error: e.Message(),
	}
	web.Respond(w, e.status(), payload)
}
