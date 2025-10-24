package order

type logger interface {
	Warn(msg string, args ...any)
}
