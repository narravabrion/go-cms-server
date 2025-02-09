package main

import (
	"log"
	"net/http"
)

 func (api *api) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("feed handler")
	ctx := r.Context()
	feed, err := api.store.Posts.GetUserFeed(ctx, int64(9))
	log.Printf("the user feed: %v", feed)
	if err != nil{
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := api.jsonResponse(w, http.StatusOK, feed); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
 }