package main

import (
	"net/http"
	"./utils"
	"./routes"
	"./models"
)


func main() {
	models.Init()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

