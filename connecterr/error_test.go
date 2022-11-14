package connecterr

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/actatum/errs"
	"github.com/bufbuild/connect-go"
)

func TestNewFromError(t *testing.T) {
	var testError = fmt.Errorf("error")
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *connect.Error
	}{
		{
			name: "invalid",
			args: args{
				err: errs.Errorf(errs.Invalid, "error"),
			},
			want: connect.NewError(connect.CodeInvalidArgument, errs.Errorf(errs.Invalid, "error")),
		},
		{
			name: "unauthorized",
			args: args{
				err: errs.Errorf(errs.Unauthorized, "error"),
			},
			want: connect.NewError(connect.CodeUnauthenticated, errs.Errorf(errs.Unauthorized, "error")),
		},
		{
			name: "permission denied",
			args: args{
				err: errs.Errorf(errs.PermissionDenied, "error"),
			},
			want: connect.NewError(connect.CodePermissionDenied, errs.Errorf(errs.PermissionDenied, "error")),
		},
		{
			name: "not found",
			args: args{
				err: errs.Errorf(errs.NotFound, "error"),
			},
			want: connect.NewError(connect.CodeNotFound, errs.Errorf(errs.NotFound, "error")),
		},
		{
			name: "internal",
			args: args{
				err: errs.Errorf(errs.Internal, "error"),
			},
			want: connect.NewError(connect.CodeInternal, errs.Errorf(errs.Internal, "error")),
		},
		{
			name: "non errs.Error",
			args: args{
				err: testError,
			},
			want: connect.NewError(connect.CodeInternal, testError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFromError(tt.args.err)
			if got.Code() != tt.want.Code() {
				t.Fatalf("Code() = %v, want %v", got.Code(), tt.want.Code())
			}
			if got.Message() != tt.want.Message() {
				t.Fatalf("Message() = %v, want %v", got.Message(), tt.want.Message())
			}
			if !reflect.DeepEqual(got.Details(), tt.want.Details()) {
				t.Fatalf("Details() = %v, want %v", got.Details(), tt.want.Details())
			}
			if !reflect.DeepEqual(got.Meta(), tt.want.Meta()) {
				t.Fatalf("Meta() = %v, want %v", got.Meta(), tt.want.Meta())
			}

			underlying := got.Unwrap()
			wantUnderlying := tt.want.Unwrap()
			if !errors.Is(underlying, wantUnderlying) {
				t.Errorf("Underlying = %v, want %v", underlying, wantUnderlying)
			}
		})
	}
}
