package ssg

import (
	"bytes"
	"fmt"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"project/ssg/config"
	"project/ssg/production"
	"regexp"
	"strings"
)

type IPageBuilder interface {
	Build() error
}

type PageBuilder struct {
	Config config.SSGData
}

type File struct {
	name string
	path string
}

func (pageBuilder *PageBuilder) Build() error {
	if production.Enabled {
		log.Println("Warning: Page builder is disabled in production")
		return nil
	}

	// read all pages inside directory pages (look for all subdirectories also)
	files, err := pageBuilder.getFilesFromDirPath(pageBuilder.Config.PagesPath)
	if err != nil {
		return err
	}

	// generate files in Debug or Release directories
	if production.Enabled {
		// create an empty files at the destination (from files variable) (production/release) by using go build tag
		fmt.Println("Generating production build (Release directory)")
		err := pageBuilder.generateFiles(pageBuilder.Config.OutReleasePath, files)
		if err != nil {
			return err
		}
	} else {
		// create an empty files at the destination (from files variable) (development/live) by using go build tag
		fmt.Println("Generating dev build (Debug directory)")
		err := pageBuilder.generateFiles(pageBuilder.Config.OutLivePath, files)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pageBuilder *PageBuilder) generateFiles(deployPath string, files []File) error {
	// Read the layout and component files.
	layoutFiles, err := pageBuilder.getFilesFromDirPath(pageBuilder.Config.LayoutsPath)
	if err != nil {

		return err
	}
	componentFiles, err := pageBuilder.getFilesFromDirPath(pageBuilder.Config.ComponentsPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		// collect all template file paths
		var templateFilePaths []string

		// then add the components
		for _, componentFile := range componentFiles {
			templateFilePaths = append(templateFilePaths, filepath.Join(componentFile.path, componentFile.name))
		}

		// add the layout (must be first, as they are collecting the others
		for _, layoutFile := range layoutFiles {
			templateFilePaths = append(templateFilePaths, filepath.Join(layoutFile.path, layoutFile.name))
		}

		// then add the pages
		templateFilePaths = append(templateFilePaths, filepath.Join(file.path, file.name))

		// Create an empty template
		tpl := template.New(file.name).Delims("{{", "}}")

		for _, templatePath := range templateFilePaths {
			// Parse each file into the template
			_, err := tpl.ParseFiles(templatePath)
			if err != nil {
				fmt.Println("Could not parse the template (probably syntax): ", err)
				fmt.Println("Example of what an error could be: space in a template comment")
			}
		}

		// execute and save the template in a buffer
		var buffer bytes.Buffer
		err = tpl.Execute(&buffer, nil)
		if err != nil {
			fmt.Println("template execution (at "+path.Join(file.path, file.name)+") error: ", err)
			//return err
		}

		var result string

		// minify the html
		m := minify.New()

		// Set up the desired HTML minifier parameters
		mHtml := &html.Minifier{
			KeepDefaultAttrVals: true,
			KeepWhitespace:      true,
			KeepDocumentTags:    true,
		}

		m.AddFunc("text/html", mHtml.Minify) // use method as func
		m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

		result, err = m.String("text/html", string(buffer.Bytes()))
		if err != nil {
			fmt.Println(err)
		}

		result = strings.ReplaceAll(result, "\n", "")

		// generate the file path (will be the result index.html from index.gohtml - or *.html from *.gohtml)
		file.path = strings.TrimPrefix(file.path, pageBuilder.Config.PagesPath)

		// make the directories to the new files
		err = os.MkdirAll(path.Join(deployPath, file.path), os.ModePerm)
		if err != nil {
			return err
		}

		// define the file name extension for upcoming new file
		if path.Ext(file.name) == ".gohtml" {
			file.name = strings.TrimSuffix(file.name, path.Ext(file.name)) + ".html"
		}

		// create the files
		newFile, err := os.Create(path.Join(deployPath, path.Join(file.path, file.name)))
		if err != nil {
			return err
		}
		defer newFile.Close()

		// Execute the template and write the resulting HTML into the file.
		_, err = newFile.Write([]byte(result))
		if err != nil {
			return err
		}
	}

	return nil
}

// getFilesFromDirPath retrieves a list of files from the specified directory path,
// including files from all subdirectories recursively.
func (pageBuilder *PageBuilder) getFilesFromDirPath(dirPath string) ([]File, error) {
	var files []File // Initialize an empty slice to store File objects.

	// Read the main directory.
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return files, err // Return an empty file list and the encountered error.
	}

	// Traverse through files and subdirectories.
	for _, file := range dir {
		if file.IsDir() {
			// If the current item is a directory, recursively call getFilesFromDirPath.
			subdirectoryPath := path.Join(dirPath, file.Name())
			subdirectoryFiles, err := pageBuilder.getFilesFromDirPath(subdirectoryPath)
			if err != nil {
				// Log the error and continue to the next iteration.
				log.Printf("Error reading subdirectory %s: %v", subdirectoryPath, err)
				continue
			}
			// Append files from the subdirectory to the main file list.
			files = append(files, subdirectoryFiles...)
		} else {
			// If the current item is a file, append it to the file list.
			files = append(files, File{name: file.Name(), path: dirPath})
		}
	}
	// Return the list of files and no error.
	return files, nil
}
