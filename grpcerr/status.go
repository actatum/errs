package grpcerr

import (
	"github.com/actatum/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StatusFromError returns a grpc status type from a given error interface.
func StatusFromError(err error) *status.Status {
	code, message := errs.ErrorCode(err), errs.ErrorMessage(err)

	if code == errs.Internal {
		message = "internal error"
	}

	return status.New(codeToGRPCCode(code), message)
}

// GRPCErrorFromError returns a grpc error from a given error.
func GRPCErrorFromError(err error) error {
	return StatusFromError(err).Err()
}

func codeToGRPCCode(code errs.Code) codes.Code {
	switch code {
	case errs.Invalid:
		return codes.InvalidArgument
	case errs.Unauthorized:
		return codes.Unauthenticated
	case errs.PermissionDenied:
		return codes.PermissionDenied
	case errs.NotFound:
		return codes.NotFound
	default:
		return codes.Internal
	}
}
