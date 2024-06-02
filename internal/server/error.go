package server

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Error          error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
	AppCode        int64  `json:"code,omitempty"`
	ErrorText      string `json:"error,omitempty"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrorInvalidRequest(err error) render.Renderer {
	return &ErrorResponse{
		Error:          err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrorRender(err error) render.Renderer {
	return &ErrorResponse{
		Error:          err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrorNotFound(err error) render.Renderer {
	return &ErrorResponse{
		Error:          err,
		HTTPStatusCode: 404,
		StatusText:     "Resource not found.",
		ErrorText:      err.Error(),
	}
}
