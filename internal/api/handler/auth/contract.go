package auth

type logger interface {
	Warn(msg string, args ...any)
}
