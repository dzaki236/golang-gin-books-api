package helpers

import "errors"

var ErrCategoryNotFound = errors.New("category not found")
var ErrCategoryExists = errors.New("category already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrUsernameTaken = errors.New("username already exists")
