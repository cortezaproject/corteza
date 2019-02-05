package service

type (
	serviceError string
)

func (e serviceError) Error() string {
	return "crust.messaging.service." + string(e)
}
