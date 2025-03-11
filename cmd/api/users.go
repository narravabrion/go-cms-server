package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/narravabrion/go-cms-server/internal/models"
	"github.com/narravabrion/go-cms-server/internal/store"
)

// type FollowerUser struct {
// 	UserID int64 `json:"user_id"`
// }

type ctxKey string

const userCtx ctxKey = "user"

// ShowAccount godoc
//
//	@Summary		Fetches user profile
//	@Description	gets the user by userID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (api *api) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// user := getUserFromCtx(r)

	idParam := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ctx := r.Context()
	user, err := api.getUser(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return

		}
	}
	if err := api.jsonResponse(w, http.StatusOK, user); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
func (api *api) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ctx := r.Context()
	err = api.store.Users.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
func (api *api) updateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (api *api) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	followedUserID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}

	ctx := r.Context()
	if err := api.store.Followers.Follow(ctx, followerUser.ID, followedUserID); err != nil {
		switch err {
		case store.ErrAlreadyFollowing:
			writeJSONError(w, http.StatusConflict, store.ErrAlreadyFollowing.Error())
			return
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err := api.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (api *api) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowerUser := getUserFromCtx(r)
	unfollowedUserID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}

	ctx := r.Context()
	if err := api.store.Followers.UnFollow(ctx, unFollowerUser.ID, unfollowedUserID); err != nil {
		switch err {
		case store.ErrNotFollowing:
			writeJSONError(w, http.StatusConflict, store.ErrNotFollowing.Error())
			return
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err := api.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (api *api) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	ctx := r.Context()
	if err := api.store.Users.Activate(ctx, token); err != nil {
		switch err {
		case store.ErrNotFound:
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err := writeJson(w, http.StatusNoContent, ""); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}

// func (api *api) userContextMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		idParam := chi.URLParam(r, "userID")
// 		id, err := strconv.ParseInt(idParam, 10, 64)
// 		if err != nil {
// 			writeJSONError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		ctx := r.Context()
// 		user, err := api.store.Users.GetByID(ctx, id)
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, store.ErrNotFound):
// 				writeJSONError(w, http.StatusNotFound, err.Error())
// 				return
// 			default:
// 				writeJSONError(w, http.StatusInternalServerError, err.Error())
// 				return

// 			}
// 		}
// 		ctx = context.WithValue(ctx, "user", user)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})

// }

func getUserFromCtx(r *http.Request) *models.User {
	user, _ := r.Context().Value("user").(*models.User)
	return user
}
