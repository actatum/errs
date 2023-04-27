package errs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestCode_String(t *testing.T) {
	type fields struct {
		slug string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				slug: "new",
			},
			want: "new",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Code{
				slug: tt.fields.slug,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Code.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodeFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want Code
	}{
		{
			name: "not found",
			args: args{
				s: "not_found",
			},
			want: NotFound,
		},
		{
			name: "invalid",
			args: args{
				s: "invalid",
			},
			want: Invalid,
		},
		{
			name: "unauthorized",
			args: args{
				s: "unauthorized",
			},
			want: Unauthorized,
		},
		{
			name: "permission denied",
			args: args{
				s: "permission_denied",
			},
			want: PermissionDenied,
		},
		{
			name: "conflict",
			args: args{
				s: "conflict",
			},
			want: Conflict,
		},
		{
			name: "internal",
			args: args{
				s: "internal",
			},
			want: Internal,
		},
		{
			name: "unknown code",
			args: args{
				s: "random",
			},
			want: Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CodeFromString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CodeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Code(t *testing.T) {
	type fields struct {
		code    Code
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   Code
	}{
		{
			name: "ok",
			fields: fields{
				code: Code{slug: "me"},
			},
			want: Code{slug: "me"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				code:    tt.fields.code,
				message: tt.fields.message,
			}
			if got := e.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Message(t *testing.T) {
	type fields struct {
		code    Code
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				message: "message",
			},
			want: "message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				code:    tt.fields.code,
				message: tt.fields.message,
			}
			if got := e.Message(); got != tt.want {
				t.Errorf("Error.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		code    Code
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				code:    Internal,
				message: "message here",
			},
			want: "error: code=internal message=message here",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				code:    tt.fields.code,
				message: tt.fields.message,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Code
	}{
		{
			name: "Error type",
			args: args{
				Errorf(Unauthorized, "invalid jwt: %s", "expired"),
			},
			want: Unauthorized,
		},
		{
			name: "Non Error type",
			args: args{
				fmt.Errorf("non Error"),
			},
			want: Internal,
		},
		{
			name: "wrapped with fmt package",
			args: args{
				fmt.Errorf("wrapped: %w", Errorf(Unauthorized, "unauthorized")),
			},
			want: Unauthorized,
		},
		{
			name: "wrapped with pkg/errors package",
			args: args{
				errors.Wrap(Errorf(Unauthorized, "unauthorized"), "wrapped"),
			},
			want: Unauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrorCode(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Error type",
			args: args{
				Errorf(PermissionDenied, "only admins can delete this thing"),
			},
			want: "only admins can delete this thing",
		},
		{
			name: "Non Error type",
			args: args{
				fmt.Errorf("non Error"),
			},
			want: "internal error",
		},
		{
			name: "wrapped with fmt package",
			args: args{
				fmt.Errorf("wrapped: %w", Errorf(Unauthorized, "unauthorized")),
			},
			want: "unauthorized",
		},
		{
			name: "wrapped with pkg/errors package",
			args: args{
				errors.Wrap(Errorf(Unauthorized, "unauthorized"), "wrapped"),
			},
			want: "unauthorized",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrorMessage(tt.args.err); got != tt.want {
				t.Errorf("ErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		code   Code
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
		want Error
	}{
		{
			name: "ok",
			args: args{
				code:   Invalid,
				format: "bad request: %s",
				args: []interface{}{
					"wow",
				},
			},
			want: Error{
				code:    Invalid,
				message: "bad request: wow",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Errorf(tt.args.code, tt.args.format, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Errorf() = %v, want %v", got, tt.want)
			}
		})
	}
}
