package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/narravabrion/go-cms-server/internal/models"
)

func (api *api) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				writeJSONError(w, http.StatusUnauthorized, "authorization header is missing")
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				writeJSONError(w, http.StatusUnauthorized, "authorization header is malformed")
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, err.Error())
				return
			}

			username := api.config.auth.basic.user
			password := api.config.auth.basic.password

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
				return

			}
			next.ServeHTTP(w, r)
		})
	}
}

func (api *api) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, http.StatusUnauthorized, "authorization header is missing")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSONError(w, http.StatusUnauthorized, "authorization header is malformed")
			return
		}
		token := parts[1]
		JWTToken, err := api.authenticator.ValidateToken(token)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, err.Error())
			return
		}
		claims, _ := JWTToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx := r.Context()
		user, err := api.store.Users.GetByID(ctx, userID)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *api) checkPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user:= getUserFromCtx(r)
		post :=  getPostFromCtx(r)

		if post.UserID == user.ID {
			next.ServeHTTP(w,r)
			return
		}
		allowed, err := api.checkRolePrecedence(r.Context(), user, role)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return 
		}
		if !allowed {
			writeJSONError(w, http.StatusForbidden, "the user does not have access")
			return
		}
		next.ServeHTTP(w,r)
	})
}

func (api *api) checkRolePrecedence(ctx context.Context, user *models.User, roleName string ) (bool, error) {
	role, err := api.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}