package user

type logger interface {
	Warn(msg string, args ...any)
}
