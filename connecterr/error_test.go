package connecterr

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/actatum/errs"
	"github.com/bufbuild/connect-go"
)

func TestNewFromError(t *testing.T) {
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
				err: fmt.Errorf("error"),
			},
			want: connect.NewError(connect.CodeInternal, fmt.Errorf("error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFromError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromError() = %v, want %v", got, tt.want)
			}
		})
	}
}
