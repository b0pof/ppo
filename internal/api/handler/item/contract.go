package item

type logger interface {
	Warn(msg string, args ...any)
}
