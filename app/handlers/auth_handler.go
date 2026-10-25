package handlers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"books-api/app/helpers"
	"books-api/app/logics"
	"books-api/app/structs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const accessTokenLifetime = 24 * time.Hour

func Register(context *gin.Context) {
	var request structs.AuthRequest
	if err := context.ShouldBindJSON(&request); err != nil || !validCredentials(request) {
		helpers.Respond(context, http.StatusBadRequest, "username must be 3-100 characters and password 8-72 characters", nil)
		return
	}
	request.Username = strings.TrimSpace(request.Username)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "could not create account", nil)
		return
	}
	user, err := logics.CreateUser(context.Request.Context(), request.Username, string(passwordHash))
	if errors.Is(err, helpers.ErrUsernameTaken) {
		helpers.Respond(context, http.StatusConflict, "username is already registered", nil)
		return
	}
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "could not create account", nil)
		return
	}
	helpers.Respond(context, http.StatusCreated, "registration successful", user)
}

func Login(context *gin.Context) {
	var request structs.AuthRequest
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be set and contain at least 32 characters")
	}
	if err := context.ShouldBindJSON(&request); err != nil || !validCredentials(request) {
		helpers.Respond(context, http.StatusBadRequest, "username and password are required", nil)
		return
	}

	request.Username = strings.TrimSpace(request.Username)
	user, err := logics.FindUserByUsername(context.Request.Context(), request.Username)
	if err != nil {
		if errors.Is(err, helpers.ErrUserNotFound) {
			helpers.Respond(context, http.StatusUnauthorized, "invalid username or password", nil)
			return
		}
		helpers.Respond(context, http.StatusInternalServerError, "could not authenticate", nil)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		helpers.Respond(context, http.StatusUnauthorized, "invalid username or password", nil)
		return
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(accessTokenLifetime)
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.ID),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		helpers.Respond(context, http.StatusInternalServerError, "could not create access token", nil)
		return
	}

	user.Password = ""
	helpers.Respond(context, http.StatusOK, "login successful", structs.AuthResponse{
		User:        user,
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(accessTokenLifetime.Seconds()),
	})
}

func validCredentials(request structs.AuthRequest) bool {
	username := strings.TrimSpace(request.Username)
	return len(username) >= 3 && len(username) <= 100 && len(request.Password) >= 8 && len(request.Password) <= 72
}
