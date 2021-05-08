package main

import (
	"animar/v1/anime"
	"animar/v1/helper"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello world!")

	/*   Anime database   */
	http.HandleFunc("/db/", helper.Handle(anime.AnimeView))

	http.ListenAndServe(":8080", nil)
}
