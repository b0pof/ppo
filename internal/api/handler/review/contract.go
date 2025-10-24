package review

type logger interface {
	Warn(msg string, args ...any)
}
