package user_test

import "github.com/b0pof/ppo/internal/model"

type UserBuilder struct {
	user model.User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		user: model.User{
			ID:       1,
			Name:     "testName",
			Login:    "testLogin",
			Role:     "testRole",
			Phone:    "88005553535",
			Password: "testPassword",
		},
	}
}

func (b *UserBuilder) WithID(id int64) *UserBuilder {
	b.user.ID = id
	return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.user.Name = name
	return b
}

func (b *UserBuilder) WithLogin(login string) *UserBuilder {
	b.user.Login = login
	return b
}

func (b *UserBuilder) WithRole(role string) *UserBuilder {
	b.user.Role = role
	return b
}

func (b *UserBuilder) WithPhone(phone string) *UserBuilder {
	b.user.Phone = phone
	return b
}

func (b *UserBuilder) WithPassword(password string) *UserBuilder {
	b.user.Password = password
	return b
}

func (b *UserBuilder) Build() model.User {
	return b.user
}
