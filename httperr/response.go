package httperr

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/actatum/errs"
)

// ErrorResponse represents an error presented to http clients.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RenderError returns a json formatted error response to the caller.
func RenderError(err error, w http.ResponseWriter) {
	code, message := errs.ErrorCode(err), errs.ErrorMessage(err)

	if code == errs.Internal {
		message = "internal error"
	}

	resp := ErrorResponse{
		Code:    code.String(),
		Message: message,
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeToStatus(code))
	_, _ = w.Write(buf.Bytes())
}

func codeToStatus(code errs.Code) int {
	switch code {
	case errs.Invalid:
		return http.StatusBadRequest
	case errs.Unauthorized:
		return http.StatusUnauthorized
	case errs.PermissionDenied:
		return http.StatusForbidden
	case errs.NotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
