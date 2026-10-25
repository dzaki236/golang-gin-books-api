package helpers

import "errors"

var ErrCategoryNotFound = errors.New("kategori tidak ditemukan")
var ErrBookNotFound = errors.New("buku tidak ditemukan")
var ErrCategoryExists = errors.New("kategori sudah ada")
var ErrUserNotFound = errors.New("pengguna tidak ditemukan")
var ErrUsernameTaken = errors.New("nama pengguna sudah digunakan")
