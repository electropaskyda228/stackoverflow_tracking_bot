package common

type ChangeExisting struct{}

func (ce *ChangeExisting) Error() string {
	return "Tracking is already existing"
}

type BadRequest struct{}

func (br *BadRequest) Error() string {
	return "Bad request"
}

type APIError struct{}

func (ae *APIError) Error() string {
	return "Problems with api"
}
