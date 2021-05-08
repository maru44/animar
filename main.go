package main

import (
	"animar/v1/anime"
	"animar/v1/helper"
	"net/http"
)

func main() {

	/*   Anime database   */
	http.HandleFunc("/db/anime/", helper.Handle(anime.AnimeView))
	http.HandleFunc("/db/anime/post/", helper.Handle(anime.AnimePostView))

	/**/

	http.ListenAndServe(":8080", nil)
}
