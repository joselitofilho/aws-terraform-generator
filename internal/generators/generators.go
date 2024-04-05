package generators

import (
	"strings"
	"text/template"

	templategenerators "github.com/diagram-code-generator/template/pkg/generators"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

// NewGenerator initialises a new instance of templategenerators.TemplateGenerator with additional template functions
// provided as a template.FuncMap.
func NewGenerator() *templategenerators.TemplateGenerator {
	return templategenerators.NewTemplateGenerator(
		templategenerators.WithExtraFuncs(template.FuncMap{
			"getFileByName":  func(files map[string]config.File, name string) config.File { return files[name] },
			"getFileImports": func(files map[string]config.File, name string) []string { return files[name].Imports },
		}),
	)
}

// MustGenerateFile generates a single file using the provided template. It logs any errors encountered during the
// generation process, except for the ErrUnsupportedFileType error, which is ignored.
func MustGenerateFile(tg *templategenerators.TemplateGenerator,
	templatesMap map[string]string, fileName, fileTmpl, outputFile string, data any,
) {
	err := tg.GenerateFile(templatesMap, fileName, fileTmpl, outputFile, data)
	if err != nil {
		fmtcolor.Yellow.Println(err)
	}
}

// MustGenerateFiles generates multiple files at once using the provided templates. It logs any errors encountered
// during the generation process.
func MustGenerateFiles(tg *templategenerators.TemplateGenerator,
	defaultTemplatesMap map[string]string, filesMap map[string]File, data any, output string,
) {
	templatesMap := map[string]string{}
	for k, file := range filesMap {
		templatesMap[k] = file.Tmpl
	}

	err := tg.GenerateFiles(defaultTemplatesMap, templatesMap, data, output)
	if err != nil {
		fmtcolor.Yellow.Println(err)
	}
}

// CreateFilesMap creates a map of file configurations from a slice of config.File structs. Each element in the slice
// corresponds to a key-value pair in the map, where the key is the file name and the value is a File struct containing
// template and import information.
func CreateFilesMap(files []config.File) map[string]File {
	filesConf := map[string]File{}
	for i := range files {
		filesConf[files[i].Name] = File{
			Tmpl:    files[i].Tmpl,
			Imports: files[i].Imports,
		}
	}

	return filesConf
}

// CreateTemplatesMap creates a map of templates from a slice of config.FilenameTemplateMap structs. Each struct
// represents a map where keys are file names and values are corresponding templates. It merges these maps into a
// single map.
func CreateTemplatesMap(filenameTemplatesList []config.FilenameTemplateMap) map[string]string {
	templatesMap := map[string]string{}

	for i := range filenameTemplatesList {
		for filename, tmpl := range filenameTemplatesList[i] {
			templatesMap[filename] = tmpl
		}
	}

	return templatesMap
}

// FilterTemplatesMap filters a map of filenames to templates based on a given filter string. It iterates over each
// key-value pair in the input map and adds pairs to a new map if the filename contains the filter string.
func FilterTemplatesMap(filter string, templatesMap map[string]string) map[string]string {
	filtred := map[string]string{}

	for filename, tmpl := range templatesMap {
		if strings.Contains(filename, filter) {
			filtred[filename] = tmpl
		}
	}

	return filtred
}
