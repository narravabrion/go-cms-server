package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/narravabrion/go-cms-server/internal/models"
	"github.com/narravabrion/go-cms-server/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=2000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=2000"`
}

func (api *api) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := getUserFromCtx(r)

	post := &models.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  user.ID,
		Tags:    payload.Tags,
	}

	ctx := r.Context()
	if err := api.store.Posts.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := api.jsonResponse(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (api *api) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	if err := api.jsonResponse(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (api *api) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}
	ctx := r.Context()

	err = api.store.Posts.Delete(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func (api *api) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if err := api.store.Posts.Update(r.Context(), post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := api.jsonResponse(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (api *api) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		ctx := r.Context()
		post, err := api.store.Posts.GetByID(ctx, id)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				writeJSONError(w, http.StatusNotFound, err.Error())
			default:
				writeJSONError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		ctx = context.WithValue(ctx, "post", post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *models.Post {
	post, _ := r.Context().Value("post").(*models.Post)
	return post
}
