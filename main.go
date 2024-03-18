package main

import (
	"backend/ssg"
	"fmt"
	"log"
	"net/http"
)

var frontendPath = "frontend/live/public"
var jsPath = "frontend/live/public/js"
var host = "127.0.0.1:8080"

func main() {
	go useStaticSiteGenerator() // autorefresher uses static site generator
	go useWebserver()

	select {} // block until termination
}

func useWebserver() {
	fileServer := http.FileServer(http.Dir(frontendPath))
	http.Handle("/", fileServer)

	log.Println("Listening on http://" + host)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func useStaticSiteGenerator() {
	ssg.JsPath = jsPath
	err := ssg.Run()
	if err != nil {
		fmt.Println(err)
	}
}
