package config

import (
	"encoding/json"
	"log"
	"os"
)

type SSGData struct {
	FilePath       string
	PagesPath      string `json:"pagesPath"`
	LayoutsPath    string `json:"layoutsPath"`
	CSSPath        string `json:"cssPath"`
	JSPath         string `json:"jsPath"`
	CSSOutPath     string `json:"cssOutPath"`
	JSOutPath      string `json:"jsOutPath"`
	ComponentsPath string `json:"componentsPath"`
	OutLivePath    string `json:"outLivePath"`
	OutReleasePath string `json:"outReleasePath"`
}

type WebserverData struct {
	FilePath    string
	Development struct {
		Host    string `json:"host"`
		WebPath string `json:"webPath"`
	} `json:"development"`
	Production struct {
		Host    string `json:"host"`
		WebPath string `json:"webPath"`
		TLS     struct {
			CertPath string `json:"certPath"`
			KeyPath  string `json:"keyPath"`
		} `json:"tls"`
	} `json:"production"`
}

func (data *SSGData) Load() {
	configFile, err := os.ReadFile(data.FilePath)
	if err != nil {
		log.Fatalln("Could not read config file: ", err)
	}

	err = json.Unmarshal(configFile, data)
	if err != nil {
		log.Fatalln("Could not unmarshal config file: ", err)
	}
}

func (data *WebserverData) Load() {
	configFile, err := os.ReadFile(data.FilePath)
	if err != nil {
		log.Fatalln("Could not read config file: ", err)
	}

	err = json.Unmarshal(configFile, data)
	if err != nil {
		log.Fatalln("Could not unmarshal config file: ", err)
	}
}
