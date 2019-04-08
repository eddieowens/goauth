package goauth

type InvalidParameterError struct {
	Msg string
}

func (c InvalidParameterError) Error() string {
	return c.Msg
}

type AuthFailedError struct {
	Msg string
}

func (a AuthFailedError) Error() string {
	return a.Msg
}
