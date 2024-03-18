package ssg

import (
	"backend/builder"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/net/websocket"
	"net/http"
)

var (
	clients      = make(map[*websocket.Conn]bool) // connected clients
	clientsMutex = sync.Mutex{}
	broadcast    = make(chan string) // broadcast channel
	builderMutex = sync.Mutex{}      // mutex for the builder process
	isBuilding   = false             // flag to prevent checking for file changes when building
	buildEndTime time.Time           // time at the end of the build, used for comparison

	JsPath              = "frontend/dev/js"
	gohtmlPath          = "frontend/dev/pages"
	requiredDirectories = []string{JsPath, gohtmlPath}
)

func Run() error {
	http.Handle("/ws", websocket.Handler(handleConnections))

	go handleMessages()
	go checkFileChanges()

	return nil
}

func handleConnections(ws *websocket.Conn) {
	clientsMutex.Lock()
	clients[ws] = true
	clientsMutex.Unlock()

	defer ws.Close()
	for {
		// We're just keeping the connection open
		time.Sleep(500 * time.Millisecond)
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		clientsMutex.Lock()
		for client := range clients {
			err := websocket.Message.Send(client, msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		clientsMutex.Unlock()
	}
}

func checkFileChanges() {
	var lastModTime time.Time
	paths := requiredDirectories

	for {
		//loopTime := time.Now()             // record the start of this check
		time.Sleep(100 * time.Millisecond) // the sleep can be adjusted according to your requirement

		builderMutex.Lock()
		if isBuilding {
			builderMutex.Unlock()
			continue
		}

		modTime, err := getLastModifiedTimestamp(paths)
		if err != nil {
			log.Println("Error checking file modification time:", err)
			broadcast <- "Error: " + err.Error()
			builderMutex.Unlock()
			continue
		}

		if modTime.After(lastModTime) {
			// only trigger if the modification time is after the last builder end time
			isBuilding = true
			builderMutex.Unlock()

			err := builder.Build()
			isBuilding = false
			buildEndTime = time.Now()

			if err != nil {
				log.Println(err)
			}

			log.Println("Triggering Full Reload")
			// broadcast file changed, we need page reload
			broadcast <- "fileChanged"
			lastModTime = modTime
		} else {
			builderMutex.Unlock()
		}
	}
}

func getLastModifiedTimestamp(folderPaths []string) (time.Time, error) {
	var newestModTime time.Time // default is zero time, i.e., an "old enough" date

	for _, folderPath := range folderPaths {
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			continue
		}
		err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.ModTime().After(newestModTime) {
				newestModTime = info.ModTime()
			}
			return nil
		})
		if err != nil {
			return time.Time{}, err
		}
	}

	return newestModTime, nil
}
