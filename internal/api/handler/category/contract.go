package category

type logger interface {
	Warn(msg string, args ...any)
}
