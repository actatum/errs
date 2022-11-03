package httperr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/actatum/errs"
)

func TestRenderError(t *testing.T) {
	type args struct {
		err error
		w   http.ResponseWriter
	}
	tests := []struct {
		name       string
		args       args
		want       ErrorResponse
		wantStatus int
	}{
		{
			name: "invalid",
			args: args{
				err: errs.Errorf(errs.Invalid, "bad request: %s", "cheese is required"),
				w:   httptest.NewRecorder(),
			},
			want: ErrorResponse{
				Code:    errs.Invalid.String(),
				Message: "bad request: cheese is required",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "unauthorized",
			args: args{
				err: errs.Errorf(errs.Unauthorized, "invalid jwt: %s", "expired"),
				w:   httptest.NewRecorder(),
			},
			want: ErrorResponse{
				Code:    errs.Unauthorized.String(),
				Message: "invalid jwt: expired",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "permission denied",
			args: args{
				err: errs.Errorf(errs.PermissionDenied, "only spongebob or patrick can use the conch"),
				w:   httptest.NewRecorder(),
			},
			want: ErrorResponse{
				Code:    errs.PermissionDenied.String(),
				Message: "only spongebob or patrick can use the conch",
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "not found",
			args: args{
				err: errs.Errorf(errs.NotFound, "no thing with id: %s", "0"),
				w:   httptest.NewRecorder(),
			},
			want: ErrorResponse{
				Code:    errs.NotFound.String(),
				Message: "no thing with id: 0",
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal",
			args: args{
				err: errs.Errorf(errs.Internal, "super secret error message"),
				w:   httptest.NewRecorder(),
			},
			want: ErrorResponse{
				Code:    errs.Internal.String(),
				Message: "internal error",
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "non errs.Error",
			args: args{
				err: fmt.Errorf("another super secret error message"),
				w:   httptest.NewRecorder(),
			},
			want: ErrorResponse{
				Code:    errs.Internal.String(),
				Message: "internal error",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RenderError(tt.args.err, tt.args.w)
			w := tt.args.w.(*httptest.ResponseRecorder)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("RenderError() gotStatus = %v, wantStatus %v", resp.StatusCode, tt.wantStatus)
			}

			want, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatal(err)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(bytes.TrimSpace(got), bytes.TrimSpace(want)) {
				t.Fatalf("RenderError() got = %v, want %v", got, want)
			}
		})
	}
}
