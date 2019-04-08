package goauth

type onAuthFunction func(authCode string)

type OAuth struct {
	// The client ID used within the OAuth protocol.
	ClientId string

	// The client secret used within the OAuth protocol.
	ClientSecret string

	// Space-delimited list of APIs you wish to access
	Scope string

	// OAuth for the local server that starts.
	ServerConfig ServerConfig

	provider oAuthProvider

	error error
}

// Configuration for the local server that starts to handle the OAuth process
type ServerConfig struct {
	// The port to start the local server on. Defaults to port 80.
	Port int

	// The URL that the local server will redirect to upon auth success.
	AuthFailedUrl string

	// The URL that the local server will redirect to upon auth failure.
	AuthSuccessUrl string

	oAuthUrl string

	redirectUri string
}

// Google OAuth 2
func Google(oAuth OAuth) *GoogleOAuth {
	oAuth.provider = GOOGLE
	if err := validateConfig(&oAuth); err != nil {
		oAuth.error = err
	}

	return &GoogleOAuth{OAuth: &oAuth}
}

func validateConfig(config *OAuth) error {
	if config.ClientId == "" {
		return InvalidParameterError{Msg: "Client ID is required for OAuth"}
	}

	if config.ClientSecret == "" {
		return InvalidParameterError{Msg: "Client secret is required for OAuth"}
	}

	if config.Scope == "" {
		return InvalidParameterError{Msg: "A scope is required for OAuth"}
	}

	if config.provider == "" {
		config.provider = GOOGLE
	}

	if config.ServerConfig.Port == 0 {
		config.ServerConfig.Port = 80
	}

	return nil
}
