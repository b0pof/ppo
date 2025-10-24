package auth

import "time"

const (
	defaultSessionTTL = 12 * time.Hour
)

type Option func(*Repository)

type Repository struct {
	redis      redisClient
	sessionTTL time.Duration
}

func New(r redisClient, opts ...Option) *Repository {
	repo := &Repository{
		redis:      r,
		sessionTTL: defaultSessionTTL,
	}

	for _, opt := range opts {
		opt(repo)
	}

	return repo
}

func WithSessionTTL(ttl time.Duration) Option {
	return func(r *Repository) {
		r.sessionTTL = ttl
	}
}
