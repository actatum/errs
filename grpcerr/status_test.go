package grpcerr

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/actatum/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStatusFromError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *status.Status
	}{
		{
			name: "invalid",
			args: args{
				err: errs.Errorf(errs.Invalid, "error"),
			},
			want: status.New(codes.InvalidArgument, "error"),
		},
		{
			name: "unauthorized",
			args: args{
				err: errs.Errorf(errs.Unauthorized, "error"),
			},
			want: status.New(codes.Unauthenticated, "error"),
		},
		{
			name: "permission denied",
			args: args{
				err: errs.Errorf(errs.PermissionDenied, "error"),
			},
			want: status.New(codes.PermissionDenied, "error"),
		},
		{
			name: "not found",
			args: args{
				err: errs.Errorf(errs.NotFound, "error"),
			},
			want: status.New(codes.NotFound, "error"),
		},
		{
			name: "internal",
			args: args{
				err: errs.Errorf(errs.Internal, "error"),
			},
			want: status.New(codes.Internal, "internal error"),
		},
		{
			name: "non errs.Error",
			args: args{
				err: fmt.Errorf("error"),
			},
			want: status.New(codes.Internal, "internal error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusFromError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusFromError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCErrorFromError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantCode    codes.Code
		wantMessage string
	}{
		{
			name: "internal",
			args: args{
				err: errs.Errorf(errs.Internal, "error"),
			},
			wantErr:     true,
			wantCode:    codes.Internal,
			wantMessage: "internal error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GRPCErrorFromError(tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCErrorFromError() error = %v, wantErr %v", err, tt.wantErr)
			}

			st, ok := status.FromError(err)
			if !ok {
				t.Fatalf("error should be of type *status.Status, got = %T", err)
			}

			if st.Code() != tt.wantCode {
				t.Fatalf("st.Code() got = %v, want %v", st.Code(), tt.wantCode)
			}
		})
	}
}
