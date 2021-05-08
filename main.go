package main

import (
	"animar/v1/anime"
	"animar/v1/auth"
	"animar/v1/helper"
	"net/http"
)

func main() {

	helper.SetEnviron()

	/*   Anime database   */
	http.HandleFunc("/db/anime/", helper.Handle(anime.AnimeView))
	http.HandleFunc("/db/anime/post/", helper.Handle(anime.AnimePostView))

	/**/

	/*   auth   */
	http.HandleFunc("/auth/sample/", auth.UserListView)

	http.ListenAndServe(":8080", nil)
}
