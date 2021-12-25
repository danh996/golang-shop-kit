package enum

type AuthentiCationType string

var (
	Cookies             AuthentiCationType = "cookies"
	BearerAuthorization AuthentiCationType = "bearer_authorization"
)
