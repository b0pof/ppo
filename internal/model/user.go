package model

const (
	RoleGuest  = "guest"
	RoleBuyer  = "buyer"
	RoleSeller = "seller"
)

type User struct {
	ID       int64
	Name     string
	Login    string
	Role     string
	Phone    string
	Password string
}

type UserMetaInfo struct {
	Name            string
	CartItemsAmount int
}

type Seller struct {
	ID   int64
	Name string
}
