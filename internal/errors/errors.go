package errors

type DiscordError struct {
	Message string
}

func (de DiscordError) Error() string {
	return de.Message
}

func CreateUnauthorizedUserError(userId string) error {
	return DiscordError{
		Message: "User " + userId + " is not authorized to use this command",
	}
}

func CreateInvalidArgumentError(arg string) error {
	return DiscordError{
		Message: "Invalid argument " + arg,
	}
}
