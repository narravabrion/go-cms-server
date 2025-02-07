package main

import "net/http"

 func (api *api) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	feed, err := api.store.Posts.GetUserFeed(ctx, int64(49))
	if err != nil{
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := api.jsonResponse(w, http.StatusOK, feed); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
 }