package controllers

import (
	"html/template"
	"net/http"

	"github.com/sirupsen/logrus"
)

type IndexData struct {
	Title      string
	IsLoggedIn bool
}

func Index(w http.ResponseWriter, r *http.Request) {
	data := IndexData{
		Title:      "Home",
		IsLoggedIn: true,
	}

	temp, err := template.New("").ParseFiles("./templates/index.html", "./templates/base.html")
	if err != nil {
		logrus.WithError(err).Fatalln("Could not render Template")
	}
	temp.ExecuteTemplate(w, "base", data)
}
