package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"books-api/app/helpers"
	"books-api/app/logics"
	"books-api/app/structs"

	"github.com/gin-gonic/gin"
)

func ListBooks(context *gin.Context) {
	books, err := logics.ListBooks(context.Request.Context())
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "daftar buku tidak dapat diambil", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "daftar buku berhasil diambil", books)
}

func GetBook(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id < 1 {
		helpers.Respond(context, http.StatusBadRequest, "ID buku tidak valid", nil)
		return
	}
	book, err := logics.FindBook(context.Request.Context(), id)
	if errors.Is(err, helpers.ErrBookNotFound) {
		helpers.Respond(context, http.StatusNotFound, "buku tidak ditemukan", nil)
		return
	}
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "buku tidak dapat diambil", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "buku berhasil diambil", book)
}

func CreateBook(context *gin.Context) {
	var request structs.BookRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		helpers.Respond(context, http.StatusBadRequest, "data buku tidak valid", nil)
		return
	}
	if message := validateBookRequest(&request); message != "" {
		helpers.Respond(context, http.StatusBadRequest, message, nil)
		return
	}
	book, err := logics.CreateBook(context.Request.Context(), request)
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "buku tidak dapat dibuat", nil)
		return
	}
	helpers.Respond(context, http.StatusCreated, "buku berhasil dibuat", book)
}

func UpdateBook(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id < 1 {
		helpers.Respond(context, http.StatusBadRequest, "ID buku tidak valid", nil)
		return
	}
	var request structs.BookRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		helpers.Respond(context, http.StatusBadRequest, "data buku tidak valid", nil)
		return
	}
	if message := validateBookRequest(&request); message != "" {
		helpers.Respond(context, http.StatusBadRequest, message, nil)
		return
	}
	book, err := logics.UpdateBook(context.Request.Context(), id, request)
	if errors.Is(err, helpers.ErrBookNotFound) {
		helpers.Respond(context, http.StatusNotFound, "buku tidak ditemukan", nil)
		return
	}
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "buku tidak dapat diperbarui", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "buku berhasil diperbarui", book)
}

func DeleteBook(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id < 1 {
		helpers.Respond(context, http.StatusBadRequest, "ID buku tidak valid", nil)
		return
	}
	if err := logics.DeleteBook(context.Request.Context(), id); errors.Is(err, helpers.ErrBookNotFound) {
		helpers.Respond(context, http.StatusNotFound, "buku tidak ditemukan", nil)
		return
	} else if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "buku tidak dapat dihapus", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "buku berhasil dihapus", nil)
}

func validateBookRequest(request *structs.BookRequest) string {
	request.Title = strings.TrimSpace(request.Title)
	if request.Title == "" || len(request.Title) > 255 {
		return "judul harus terdiri dari 1-255 karakter"
	}
	if request.ReleaseYear < 1980 || request.ReleaseYear > 2024 {
		return "tahun terbit (release_year) harus antara 1980 dan 2024"
	}
	if request.Price < 0 {
		return "harga tidak boleh kurang dari nol"
	}
	if request.TotalPage < 1 {
		return "jumlah halaman (total_page) harus lebih dari 0"
	}
	if request.CategoryID < 1 {
		return "ID kategori (category_id) wajib diisi"
	}
	return ""
}
