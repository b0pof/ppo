package model

import (
	"regexp"
	"unicode/utf8"
)

const (
	userNamePattern = `^[А-Яа-яA-Za-z0-9_]*$`
	loginPattern    = `^[A-Za-z0-9_]*$`
	phonePattern    = `^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$`
	passwordPattern = `^[A-Za-z0-9_]*$`

	itemNamePattern        = `^[А-Яа-яA-Za-z0-9_ ]*$`
	itemDescriptionPattern = `^[А-Яа-яA-Za-z0-9_ ]*$`
	itemImgSrcPattern      = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
)

// User

func ValidateUser(user User) error {
	var err error

	//if user.Name != "" {
	err = ValidateUserName(user.Name)
	if err != nil {
		return err
	}
	//}

	if user.Phone != "" {
		err = ValidateUserPhone(user.Phone)
		if err != nil {
			return err
		}
	}

	if user.Login != "" {
		err = ValidateUserLogin(user.Login)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateUserName(name string) error {
	length := utf8.RuneCountInString(name)

	if length < 1 || length > 32 {
		return NewValidationError("Допустимая длина имени - от 1 до 32 символов")
	}

	if !fitsPattern(name, userNamePattern) {
		return NewValidationError("Используйте в имени только буквы, цифры и нижнее подчеркивание")
	}

	return nil
}

func ValidateUserPassword(password string) error {
	length := utf8.RuneCountInString(password)

	if length < 8 || length > 32 {
		return NewValidationError("Допустимая длина пароля - от 8 до 32 символов")
	}

	if !fitsPattern(password, passwordPattern) {
		return NewValidationError("Используйте в пароле только латинские символы, цифры и нижнее подчеркивание")
	}

	return nil
}

func ValidateUserPhone(phone string) error {
	if !fitsPattern(phone, phonePattern) {
		return NewValidationError("Неверный формат номер телефона")
	}

	return nil
}

func ValidateUserLogin(login string) error {
	length := utf8.RuneCountInString(login)

	if length < 4 || length > 32 {
		return NewValidationError("Допустимая длина пароля - от 8 до 32 символов")
	}

	if !fitsPattern(login, loginPattern) {
		return NewValidationError("Используйте в логине только латинские символы, цифры и нижнее подчеркивание")
	}

	return nil
}

// Item

func ValidateItem(item Item) error {
	var err error

	if item.Name != "" {
		err = ValidateItemName(item.Name)
		if err != nil {
			return err
		}
	}

	if item.Description != "" {
		err = ValidateItemDescription(item.Description)
		if err != nil {
			return err
		}
	}

	if item.ImgSrc != "" {
		err = ValidateItemImgSrc(item.ImgSrc)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateItemName(name string) error {
	length := utf8.RuneCountInString(name)

	if length < 4 || length > 128 {
		return NewValidationError("Допустимая длина названия товара - от 4 до 128 символов")
	}

	if !fitsPattern(name, itemNamePattern) {
		return NewValidationError("Используйте в имени товара только буквы, цифры и нижнее подчеркивание")
	}

	return nil
}

func ValidateItemDescription(description string) error {
	length := utf8.RuneCountInString(description)

	if length < 4 || length > 512 {
		return NewValidationError("Допустимая длина описания товара - от 4 до 512 символов")
	}

	if !fitsPattern(description, itemDescriptionPattern) {
		return NewValidationError("Используйте в описании товара только буквы, цифры и нижнее подчеркивание")
	}

	return nil
}

func ValidateItemImgSrc(url string) error {
	length := utf8.RuneCountInString(url)

	if length < 4 || length > 512 {
		return NewValidationError("Допустимая длина ссылки на фото товара - от 4 до 512 символов")
	}

	if !fitsPattern(url, itemImgSrcPattern) {
		return NewValidationError("Невалидный адрес ссылки")
	}

	return nil
}

func fitsPattern(field string, pattern string) bool {
	return regexp.MustCompile(pattern).MatchString(field)
}
