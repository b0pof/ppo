package seller

type logger interface {
	Warn(msg string, args ...any)
}
