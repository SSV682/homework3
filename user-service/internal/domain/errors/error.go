package errors

type UserError string

func (ee UserError) Error() string {
	return string(ee)
}

var (
	ErrNonExistentId   = UserError("non-existent id")
	ErrIncorrectParams = UserError("incorrect params")
	ErrConflict        = UserError("conflict")
)
