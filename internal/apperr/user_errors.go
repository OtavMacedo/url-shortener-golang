package apperr

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("invalid password")
