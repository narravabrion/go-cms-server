package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
	"github.com/narravabrion/go-cms-server/internal/models"
	"github.com/narravabrion/go-cms-server/internal/store"
)

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,min=5,max=14"`
}

type UserWithToken struct {
	*models.User
	Token string `json:"token"`
}

func (api *api) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := &models.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := r.Context()

	token := uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	err := api.store.Users.CreateAndInvite(ctx, user, hashedToken, api.config.mail.exp)
	if err != nil {

		switch err {
		case store.ErrDuplicateEmail:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		case store.ErrDuplicateUsername:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	userWithToken := UserWithToken{
		User: user,
		Token: token,
	}

	if err := writeJson(w, http.StatusCreated, userWithToken); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
