package errs

import "errors"

var ErrOwnershipTransferRequired = errors.New(
	"community ownership must be transferred before leaving",
)
