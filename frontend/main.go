package main

import (
	"log"
	"net/http"
)

func main() {
	// Start goroutine for the web server
	//	go func() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Println("Listening on http://127.0.0.1:8081...")
	err := http.ListenAndServe("127.0.0.1:8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//	}()
	/*
		dir, cmd := "src", "esbuild.js"
		var modTime time.Time
		for range time.Tick(1 * time.Second) {
			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err == nil && info.ModTime().After(modTime) {
					modTime = info.ModTime()
					log.Println("New file modification at:", modTime)
					err := exec.Command("node", cmd).Run()
					if err != nil {
						log.Println("Command execution error:", err)
					} else {
						log.Println("Rebuild successfully triggered.")
					}
				}
				return nil
			})
		}
	*/
}
