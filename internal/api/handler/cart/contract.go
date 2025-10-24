package cart

type logger interface {
	Warn(msg string, args ...any)
}
