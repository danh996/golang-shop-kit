package requestinfo

type Info struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	Ip        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

const (
	CookieKey = "h5token"
	// InfoKey ...
	InfoKey = "info_key"
	// Authorization ...
	Authorization = "Authorization"
	// BearerPrefix ...
	Bearer = "Bearer"
	// AuthHeaderKey ...
	AuthHeaderKey = "authorization"
	// ClientIPKey ...
	ClientIPKey = "grpcgateway-client-ip"
	// UserAgentKey ...
	UserAgentKey = "grpcgateway-user-agent"
	// UserIDRequestKey ...
	UserIDRequestKey = "grpcgateway-user-id-request"
)
