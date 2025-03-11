package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/narravabrion/go-cms-server/internal/mailer"
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
type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=5,max=72"`
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
		Role: models.Role{
			Name: "user",
		},
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
		User:  user,
		Token: token,
	}
	activationURL := fmt.Sprintf("%sconfirm/%s", api.config.frontEndURL, token)

	isProdEnv := api.config.env == "production"

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	err = api.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)

	if err != nil {
		api.logger.Errorw("error sending email", "error", err)
		// rollback
		if err := api.store.Users.Delete(ctx, user.ID); err != nil {
			api.logger.Errorw("error deleting user", "error", err)
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJson(w, http.StatusCreated, userWithToken); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}

func (api *api) createTokenHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			writeJSONError(w, http.StatusUnauthorized, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = user.Password.ComparePassword(payload.Password)
    if err != nil {
        writeJSONError(w, http.StatusUnauthorized, err.Error())
        return
    }
	claims := jwt.MapClaims{
		"sub" : user.ID,
		"exp": time.Now().Add(api.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": api.config.auth.token.iss,
		"aud": api.config.auth.token.iss,

	}
	token, err := api.authenticator.GenerateToken(claims)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := api.jsonResponse(w, http.StatusCreated, token); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
