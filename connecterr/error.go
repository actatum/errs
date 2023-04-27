package connecterr

import (
	"github.com/actatum/errs"
	"github.com/bufbuild/connect-go"
)

// NewFromError returns a new connect error from a given error interface.
func NewFromError(err error) *connect.Error {
	code := errs.ErrorCode(err)

	return connect.NewError(codeToConnectCode(code), err)
}

func codeToConnectCode(code errs.Code) connect.Code {
	switch code {
	case errs.Invalid:
		return connect.CodeInvalidArgument
	case errs.Unauthorized:
		return connect.CodeUnauthenticated
	case errs.PermissionDenied:
		return connect.CodePermissionDenied
	case errs.NotFound:
		return connect.CodeNotFound
	case errs.Conflict:
		return connect.CodeAlreadyExists
	default:
		return connect.CodeInternal
	}
}
