package main

import (
	"fmt"
	"log"
	"net/http"
	"project/ssg"
	"project/ssg/autoreloader"
	"project/ssg/config"
	"project/ssg/production"
)

var configWebserverFilePath = "./webserver-config.json"
var configWebserverData config.WebserverData
var configSSGFilePath = "./ssg-config.json"
var configSSGData config.SSGData

func main() {
	useWebserverConfig()
	useSSGConfig()

	go useWebserver()
	go useSSG()

	select {} // block until termination
}

func useWebserverConfig() {
	configWebserverData = config.WebserverData{FilePath: configWebserverFilePath}
	configWebserverData.Load()
}

func useWebserver() {
	data := configWebserverData
	var fileServer http.Handler
	if production.Enabled {
		fileServer = http.FileServer(http.Dir(data.Production.WebPath))
	} else {
		fileServer = http.FileServer(http.Dir(data.Development.WebPath))
	}

	http.Handle("/", fileServer)

	if production.Enabled {
		log.Println("Listening on https://" + configWebserverData.Production.Host)
		err := http.ListenAndServeTLS(data.Production.Host, data.Production.TLS.CertPath, data.Production.TLS.KeyPath, nil)
		if err != nil {
			log.Fatal("ListenAndServeTLS: ", err)
		}
	} else {
		log.Println("Listening on http://" + configWebserverData.Development.Host)
		err := http.ListenAndServe(configWebserverData.Development.Host, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}

func useSSGConfig() {
	configSSGData = config.SSGData{FilePath: configSSGFilePath}
}

func useSSG() {
	configSSGData.Load()
	autoreloader.SSGPageBuilder = &ssg.PageBuilder{Config: configSSGData}
	err := autoreloader.Run()
	if err != nil {
		fmt.Println(err)
	}
}
