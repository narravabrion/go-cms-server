package main

import (

	"net/http"

	"github.com/narravabrion/go-cms-server/internal/store"
)

func (api *api) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	pq := store.PaginationFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}
	pq, err := pq.Parse(r)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate.Struct(pq); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	feed, err := api.store.Posts.GetUserFeed(ctx, int64(9), pq)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := api.jsonResponse(w, http.StatusOK, feed); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
