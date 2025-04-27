package common

type BadRequest struct{}
type FailedAddUser struct{}
type FailedAddTracking struct{}
type FailedUntrack struct{}

func (br *BadRequest) Error() string {
	return "bad request"
}

func (fau *FailedAddUser) Error() string {
	return "failed to add user"
}

func (fat *FailedAddTracking) Error() string {
	return "failed to add tracking"
}

func (fu *FailedUntrack) Error() string {
	return "failed to untrack question"
}
