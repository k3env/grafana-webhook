package templates

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func findFilesWithExtension(rootDir, ext string) ([]string, error) {
	var files []string

	// Функция, которая будет вызываться для каждого найденного файла/директории
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Проверяем, что это файл и что его расширение совпадает с искомым
		if !info.IsDir() && filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func Load(path string) (templates map[string]*template.Template) {
	templates = map[string]*template.Template{}
	files, err := findFilesWithExtension(path, ".tmpl")

	if err != nil {
		log.Fatal().Err(err).Msg("Error finding templates")
	}

	for _, file := range files {
		contents, err := os.ReadFile(file)
		if err != nil {
			log.Error().Err(err).Str("file", file).Msg("Error reading template")
			continue
		}
		fn := strings.Replace(file, path, "templates", 1)
		tn := strings.ReplaceAll(strings.TrimSuffix(fn, ".tmpl"), string(filepath.Separator), ".")
		t, err := template.New(tn).Parse(string(contents))
		if err != nil {
			log.Error().Err(err).Str("file", file).Msg("Error parsing template")
			continue
		}
		templates[tn] = t
	}

	return templates
}
