package expiration

import "time"

const (
	codeExpirationTTL = 2 * time.Minute
)

type Expiration struct{}

func New() *Expiration {
	return &Expiration{}
}

func (e *Expiration) Get() time.Time {
	return time.Now().Add(codeExpirationTTL)
}
