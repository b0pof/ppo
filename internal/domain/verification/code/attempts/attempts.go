package attempts

const (
	maxAttempts = 3
)

func MaxReached(usedAttempts int64) bool {
	return usedAttempts >= maxAttempts
}
