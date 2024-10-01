package internal

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"os"
	"path"
)

type Logic struct {
	verbose bool
	data    map[string]interface{}
}

func NewLogic(verbose bool) *Logic {
	return &Logic{verbose: verbose}
}

func (l *Logic) Execute(dataFiles []string, templateFiles []string) {
	var jsons [][]byte

	for _, i := range dataFiles {
		data, err := os.ReadFile(i)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
		jsons = append(jsons, data)
	}

	mergedJSONs, _ := MergeJSONsToJSON(jsons...)
	//parse back to map
	l.data = JSONToMap(mergedJSONs)

	if l.verbose {
		log.Println(string(mergedJSONs))
	}

	for _, arg := range templateFiles {
		fs, err := os.Stat(arg)
		switch {
		case os.IsNotExist(err):
			log.Printf(`template "%s" does not exist`+"\n", arg)
		case fs.Mode().IsDir():
			l.dirFlow(arg)
		case fs.Mode().IsRegular():
			l.fileFlow(arg)
		default:
			log.Fatalf(`can't handle this "%s"?`, arg)
		}
	}
}

func (l *Logic) fileFlow(arg string) {
	if l.verbose {
		log.Printf("processing %s", arg)
	}
	var data []byte
	var err error
	if data, err = os.ReadFile(arg); err != nil && !json.Valid(data) {
		log.Printf(`could not parse template "%s"`, arg)
		return
	}

	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"toJSON": func(data interface{}) string {
			if jstr, err := json.Marshal(data); err == nil {
				return string(jstr)
			}
			return ""
		},
	}).Parse(string(data)))

	tmpl.Execute(os.Stdout, l.data)
}
func (l *Logic) dirFlow(arg string) {
	entries, err := os.ReadDir(arg)
	if err != nil && err != io.EOF {
		log.Printf(`can't list templates for "%s"`+"\n", arg)
		return
	}
	for _, entry := range entries {
		switch entry.IsDir() {
		case true:
			l.dirFlow(path.Join(arg, entry.Name()))
		default:
			l.fileFlow(entry.Name())
		}
	}
}
