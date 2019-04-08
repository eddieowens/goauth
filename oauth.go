package goauth

type oAuthProvider string

const (
	GOOGLE oAuthProvider = "GOOGLE"
)

type OAuthToken struct {
	Error error
}
