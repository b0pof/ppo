package model

type Permission struct {
	role string
}

type Permissions []Permission

func (p Permissions) HasAccess(role string) bool {
	for _, perm := range p {
		if perm.role == role {
			return true
		}
	}

	return false
}
