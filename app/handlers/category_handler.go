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
		helpers.Respond(context, http.StatusInternalServerError, "could not retrieve categories", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "categories retrieved", categories)
}

func GetCategory(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	category, err := logics.FindCategory(context.Request.Context(), id)
	if errors.Is(err, helpers.ErrCategoryNotFound) {
		helpers.Respond(context, http.StatusNotFound, "category not found", nil)
		return
	}
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "could not retrieve category", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "category retrieved", category)
}

func CreateCategory(context *gin.Context) {
	var request structs.CategoryRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		helpers.Respond(context, http.StatusBadRequest, "category name is required", nil)
		return
	}
	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" || len(request.Name) > 150 {
		helpers.Respond(context, http.StatusBadRequest, "category name must be 1-150 characters", nil)
		return
	}
	category, err := logics.CreateCategory(context.Request.Context(), request.Name)
	if errors.Is(err, helpers.ErrCategoryExists) {
		helpers.Respond(context, http.StatusConflict, "category name already exists", nil)
		return
	}
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "could not create category", nil)
		return
	}
	helpers.Respond(context, http.StatusCreated, "category created", category)
}

func UpdateCategory(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	var request structs.CategoryRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		helpers.Respond(context, http.StatusBadRequest, "category name is required", nil)
		return
	}
	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" || len(request.Name) > 150 {
		helpers.Respond(context, http.StatusBadRequest, "category name must be 1-150 characters", nil)
		return
	}
	category, err := logics.UpdateCategory(context.Request.Context(), id, request.Name)
	if errors.Is(err, helpers.ErrCategoryNotFound) {
		helpers.Respond(context, http.StatusNotFound, "category not found", nil)
		return
	}
	if errors.Is(err, helpers.ErrCategoryExists) {
		helpers.Respond(context, http.StatusConflict, "category name already exists", nil)
		return
	}
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "could not update category", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "category updated", category)
}

func DeleteCategory(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	if err := logics.DeleteCategory(context.Request.Context(), id); errors.Is(err, helpers.ErrCategoryNotFound) {
		helpers.Respond(context, http.StatusNotFound, "category not found", nil)
		return
	} else if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "could not delete category", nil)
		return
	}
	helpers.Respond(context, http.StatusOK, "category deleted", nil)
}
