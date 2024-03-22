package autoreloader

import (
	"log"
	"os"
	"path/filepath"
	"project/ssg"
	"project/ssg/production"
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

	// TODO: move this hardcoded path into the config (also look into making the config a singleton)
	//devPath = "frontend/out"
	devPath                          = "out/Debug/frontend/public"
	requiredDirectoriesForFileChange = []string{devPath}
	// TODO: move this hardcoded path into the config (also look into making the config a singleton)
	requiredDirectoriesForSSGBuild = []string{"frontend/src/gohtml", "frontend/src/css", "frontend/src/js"}
	SSGPageBuilder                 *ssg.PageBuilder
)

func Run() error {
	if production.Enabled {
		log.Println("Warning: autoreloader is disabled in production")
		return nil
	}
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
	var lastModTimeForFileChange time.Time // for files under directories in requiredDirectoriesForFileChange
	var lastModTimeForSSGBuild time.Time   // for files under directories in requiredDirectoriesForSSGBuild
	pathsForFileChange := requiredDirectoriesForFileChange
	pathsForSSGBuild := requiredDirectoriesForSSGBuild

	for {
		time.Sleep(100 * time.Millisecond)

		builderMutex.Lock()
		if isBuilding {
			builderMutex.Unlock()
			continue
		}

		modTime, err := getLastModifiedTimestamp(pathsForFileChange)
		if err != nil {
			log.Println("Error checking file modification time for file change:", err)
			broadcast <- "Error: " + err.Error()
			builderMutex.Unlock()
			continue
		}

		if modTime.After(lastModTimeForFileChange) {
			log.Println("Triggering Full Reload for file change")
			broadcast <- "fileChanged"
			lastModTimeForFileChange = modTime
		}

		modTime, err = getLastModifiedTimestamp(pathsForSSGBuild)
		if err != nil {
			log.Println("Error checking file modification time for SSG build:", err)
			broadcast <- "Error: " + err.Error()
			builderMutex.Unlock()
			continue
		}

		if modTime.After(lastModTimeForSSGBuild) {
			isBuilding = true
			builderMutex.Unlock()

			if SSGPageBuilder != nil {
				err := SSGPageBuilder.Build()
				if err != nil {
					log.Println(err)
				}
				isBuilding = false
				buildEndTime = time.Now()
			}

			log.Println("Triggering Full Reload After SSG Build")
			broadcast <- "fileChanged"
			lastModTimeForSSGBuild = modTime
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
