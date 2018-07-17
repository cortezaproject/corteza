package service

type (
	serviceError string
)

func (e serviceError) Error() string {
	return "crust.sam.service." + string(e)
}
