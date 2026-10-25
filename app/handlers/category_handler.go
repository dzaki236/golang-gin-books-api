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

func ListCategories(context *gin.Context) {
	categories, err := logics.ListCategories(context.Request.Context())
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "daftar kategori tidak dapat diambil", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "daftar kategori berhasil diambil", categories)
}

func GetCategory(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	category, err := logics.FindCategory(context.Request.Context(), id)
	if errors.Is(err, helpers.ErrCategoryNotFound) {
		helpers.Respond(context, http.StatusNotFound, "kategori tidak ditemukan", nil)
		return
	}
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "kategori tidak dapat diambil", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "kategori berhasil diambil", category)
}

func ListCategoryBooks(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id < 1 {
		helpers.Respond(context, http.StatusBadRequest, "ID kategori tidak valid", nil)
		return
	}
	category, err := logics.FindCategory(context.Request.Context(), id)
	if errors.Is(err, helpers.ErrCategoryNotFound) {
		helpers.Respond(context, http.StatusNotFound, "kategori tidak ditemukan", nil)
		return
	}
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "kategori tidak dapat diambil", nil)
		return
	}
	books, err := logics.ListBooksByCategory(context.Request.Context(), id)
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "daftar buku kategori tidak dapat diambil", nil)
		return
	}
	categoryBooks := gin.H{
		"id": category.ID,
		"category": gin.H{
			"name":       category.Name,
			"created_at": category.CreatedAt,
		},
		"books": books,
	}
	helpers.Respond(context, http.StatusOK, "daftar buku kategori berhasil diambil", categoryBooks)
}

func CreateCategory(context *gin.Context) {
	var request structs.CategoryRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		helpers.Respond(context, http.StatusBadRequest, "nama kategori wajib diisi", nil)
		return
	}
	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" || len(request.Name) > 150 {
		helpers.Respond(context, http.StatusBadRequest, "nama kategori harus terdiri dari 1-150 karakter", nil)
		return
	}
	category, err := logics.CreateCategory(context.Request.Context(), request.Name)
	if errors.Is(err, helpers.ErrCategoryExists) {
		helpers.Respond(context, http.StatusConflict, "nama kategori sudah digunakan", nil)
		return
	}
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "kategori tidak dapat dibuat", nil)
		return
	}
	helpers.Respond(context, http.StatusCreated, "kategori berhasil dibuat", category)
}

func UpdateCategory(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	var request structs.CategoryRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		helpers.Respond(context, http.StatusBadRequest, "nama kategori wajib diisi", nil)
		return
	}
	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" || len(request.Name) > 150 {
		helpers.Respond(context, http.StatusBadRequest, "nama kategori harus terdiri dari 1-150 karakter", nil)
		return
	}
	category, err := logics.UpdateCategory(context.Request.Context(), id, request.Name)
	if errors.Is(err, helpers.ErrCategoryNotFound) {
		helpers.Respond(context, http.StatusNotFound, "kategori tidak ditemukan", nil)
		return
	}
	if errors.Is(err, helpers.ErrCategoryExists) {
		helpers.Respond(context, http.StatusConflict, "nama kategori sudah digunakan", nil)
		return
	}
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "kategori tidak dapat diperbarui", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "kategori berhasil diperbarui", category)
}

func DeleteCategory(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	if err := logics.DeleteCategory(context.Request.Context(), id); errors.Is(err, helpers.ErrCategoryNotFound) {
		helpers.Respond(context, http.StatusNotFound, "kategori tidak ditemukan", nil)
		return
	} else if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "kategori tidak dapat dihapus", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "kategori berhasil dihapus", nil)
}
