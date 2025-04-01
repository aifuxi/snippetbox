package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// 邮箱或密码错误
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// 邮箱重复
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
