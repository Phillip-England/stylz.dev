package filehandler

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

func ParseTemplates() (*template.Template, error) {
	templates := template.New("")
	err := filepath.Walk("./html", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err := templates.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func ExecuteTemplate(t *template.Template, name string, data interface{}) template.HTML {
	var templateContent bytes.Buffer
	err := t.ExecuteTemplate(&templateContent, name, data)
	if err != nil {
		panic(err)
	}
	return template.HTML(templateContent.String())
}

func ExecuteMarkdown(filepath string) template.HTML {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileBytes := new(bytes.Buffer)
	fileBytes.ReadFrom(file)
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithGuessLanguage(true),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
	)
	var output bytes.Buffer
	if err := md.Convert(fileBytes.Bytes(), &output); err != nil {
		panic(err)
	}
	outputString := output.String()
	for i := 0; i < len(outputString); i++ {
		chunkSize := 6
		atEnd := i+chunkSize >= len(outputString)
		if atEnd {
			break
		}
		chunck := outputString[i : i+chunkSize]
		if chunck == "</pre>" {
			fmt.Println("hit")
			// insertIndex := i - 1
			// copySvg := `
			// 	<div>
			// 		<svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
			// 			<path stroke="currentColor" stroke-linejoin="round" stroke-width="2" d="M9 8v3a1 1 0 0 1-1 1H5m11 4h2a1 1 0 0 0 1-1V5a1 1 0 0 0-1-1h-7a1 1 0 0 0-1 1v1m4 3v10a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1v-7.13a1 1 0 0 1 .24-.65L7.7 8.35A1 1 0 0 1 8.46 8H13a1 1 0 0 1 1 1Z"/>
			// 		</svg>
			// 	</div>
			// `
			// outputString = outputString[:insertIndex] + copySvg + outputString[insertIndex:]
		}
	}
	return template.HTML(outputString)
}
