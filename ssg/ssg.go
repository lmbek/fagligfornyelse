package ssg

// not used anymore:
//"github.com/tdewolff/minify/v2"
//"github.com/tdewolff/minify/v2/html"
//"github.com/tdewolff/minify/v2/js"

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"project/ssg/config"
	"project/ssg/custom_minifier/htmlmin"
	"project/ssg/helper"
	"project/ssg/production"
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

func (pageBuilder *PageBuilder) generateJSFiles(deployPath string, files []File) error {
	// we should copy (map) all javascript files from /frontend/out/js to out/Debug/frontend/public/js

	// in go 1.23 we can do this:
	// err = os.CopyFS(destDir, os.DirFS(srcDir))

	fmt.Println("generating js")
	// but until then, we do:
	src := pageBuilder.Config.JSPath
	out := deployPath

	err := helper.CopyDir(out, src)
	if err != nil {
		return err
	}

	return nil
}

func (pageBuilder *PageBuilder) generateCSSFiles(deployPath string, files []File) error {
	// in go 1.23 we can do this:
	// err = os.CopyFS(destDir, os.DirFS(srcDir))

	// we should copy (map) all javascript files from /frontend/out/js to out/Debug/frontend/public/js

	// in go 1.23 we can do this:
	// err = os.CopyFS(destDir, os.DirFS(srcDir))

	fmt.Println("generating css")
	// but until then, we do:
	src := pageBuilder.Config.CSSPath
	out := deployPath

	err := helper.CopyFS(out, os.DirFS(src))
	if err != nil {
		return err
	}
	return nil
}

func (pageBuilder *PageBuilder) Build() error {
	if production.Enabled {
		log.Println("Warning: Page builder is disabled in production")
		return nil
	}

	// read all css files inside directory pages (look for all subdirectories also)
	//cssFiles, err := pageBuilder.getFilesFromDirPath(pageBuilder.Config.CSSPath)
	//if err != nil {
	//	return err
	//}

	// read all js files inside directory pages (look for all subdirectories also)
	//jsFiles, err := pageBuilder.getFilesFromDirPath(pageBuilder.Config.JSPath)
	//if err != nil {
	//	return err
	//}

	// as we only copy files, and dont do anything with them at the moment with Go, we dont need to loop the files
	var cssFiles, jsFiles []File

	// read all html files inside directory pages (look for all subdirectories also)
	htmlFiles, err := pageBuilder.getFilesFromDirPath(pageBuilder.Config.PagesPath)
	if err != nil {
		return err
	}

	// generate files in Debug or Release directories
	if production.Enabled {
		// create an empty files at the destination (from files variable) (production/release) by using go build tag
		fmt.Println("Generating production build (Release directory)")
		err = pageBuilder.generateCSSFiles(pageBuilder.Config.CSSOutPath, cssFiles)
		if err != nil {
			return err
		}

		err = pageBuilder.generateJSFiles(pageBuilder.Config.JSOutPath, jsFiles)
		if err != nil {
			return err
		}
		return pageBuilder.generateHTMLFiles(pageBuilder.Config.OutReleasePath, htmlFiles)
	} else {
		// create an empty files at the destination (from files variable) (development/live) by using go build tag
		fmt.Println("Generating dev build (Debug directory)")
		err = pageBuilder.generateCSSFiles(pageBuilder.Config.CSSOutPath, cssFiles)
		if err != nil {
			return err
		}

		err = pageBuilder.generateJSFiles(pageBuilder.Config.JSOutPath, jsFiles)
		if err != nil {
			return err
		}

		return pageBuilder.generateHTMLFiles(pageBuilder.Config.OutLivePath, htmlFiles)
	}
}

func minifyHTML(html string) ([]byte, error) {
	return htmlmin.Minify([]byte(html), &htmlmin.Options{
		MinifyScripts: false,
		MinifyStyles:  false,
		UnquoteAttrs:  false,
	})
}

func (pageBuilder *PageBuilder) generateHTMLFiles(deployPath string, files []File) error {

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

		//fmt.Println(result)

		html, err := minifyHTML(buffer.String())
		//if err != nil {
		//	return err
		//}
		/*
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
		*/
		result = strings.ReplaceAll(string(html), "\n", "")

		//fmt.Println(result)

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
