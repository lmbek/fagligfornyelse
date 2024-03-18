package builder

import (
	"bufio"
	"bytes"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	PagesPath      = "frontend/dev/pages"
	ComponentsPath = "frontend/dev/components"
	LayoutsPath    = "frontend/dev/layouts"
	OutputPath     = "frontend/live/public"
)

func Build() error {
	return generateHTMLFiles(IsDevelopmentMode())
}

func gatherFiles(path string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(p) == ".gohtml" {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func chooseLayout(pageFile string) (string, error) {
	file, err := os.Open(pageFile)
	if err != nil {
		return "", err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	firstLine, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	layoutPrefix := "// layout:"
	if strings.HasPrefix(strings.TrimSpace(firstLine), layoutPrefix) {
		// Extract layout name
		layoutName := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(firstLine), layoutPrefix))
		return layoutName, nil
	}

	return "default", nil
}

func generateHTMLFiles(IsDevelopment bool) error {
	tmpl := template.New("").Delims("{{", "}}")

	layoutFiles, err := gatherFiles(LayoutsPath)
	if err != nil {
		return err
	}
	for _, fpath := range layoutFiles {
		tmpl, err = tmpl.ParseFiles(fpath)
		if err != nil {
			return err
		}
	}

	componentFiles, err := gatherFiles(ComponentsPath)
	if err != nil {
		return err
	}
	for _, fpath := range componentFiles {
		tmpl, err = tmpl.ParseFiles(fpath)
		if err != nil {
			return err
		}
	}

	pageFiles, err := gatherFiles(PagesPath)
	if err != nil {
		return err
	}

	for _, pageFile := range pageFiles {
		chosenLayout, err := chooseLayout(pageFile)
		if err != nil {
			return err
		}

		buf := &bytes.Buffer{}

		pageTmpl, err := tmpl.Clone()
		if err != nil {
			return err
		}

		pageTmpl, err = pageTmpl.ParseFiles(pageFile)
		if err != nil {
			return err
		}

		err = pageTmpl.ExecuteTemplate(buf, chosenLayout, struct{ IsDevelopment bool }{IsDevelopment})
		if err != nil {
			return err
		}

		relativePath := strings.TrimPrefix(filepath.ToSlash(pageFile), filepath.ToSlash(PagesPath+"/"))

		htmlFilename := strings.TrimSuffix(relativePath, filepath.Ext(relativePath)) + ".html"
		outPath := filepath.Join(OutputPath, htmlFilename)

		if err := os.MkdirAll(filepath.Dir(outPath), fs.ModePerm); err != nil {
			return err
		}

		if err := os.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
			return err
		}
	}

	return nil
}
