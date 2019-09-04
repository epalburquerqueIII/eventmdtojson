package parse

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	s "strings"
)

var filePattern = regexp.MustCompile(`\.md$`)

type frontmatter struct {
	title string
}

type mdFile struct {
	filename string
}

// Parsed represents a single parsed file
type Parsed struct {
	Title         string `json:"title"`
	Date          string `json:"date"`
	Description   string `json:"description"`
	Tipo          string `json:"tipo"`
	Image         string `json:"image"`
	Imageslide    string `json:"imageslide"`
	Author        string `json:"author"`
	Identificator string `json:"identificador"`
	Categorias    string `json:"categorias"`
	Tags          string `json:"tags"`
	Body          string `json:"body"`
}

func readMDFiles(dir string) ([]mdFile, error) {
	mdFiles := []mdFile{}
	d, err := os.Open(dir)

	defer d.Close()
	files, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, n := range files {
		if filePattern.Match([]byte(n)) {
			_, err := ioutil.ReadFile(filepath.Join(dir, n)) //ioutil para entrada y salida de datos
			if err != nil {
				log.Printf("Cannot read file %s/%s", dir, n)
				continue
			}
			newFile := mdFile{
				filename: n,
				//bytes:    f
			}
			mdFiles = append(mdFiles, newFile) //almacena los datos de los ficheros en un array
		}
	}
	return mdFiles, err
}
func analizalinea(linea string) (string, string) {

	var campo, dato string

	a := s.Index(linea, ":")
	if a == -1 {
		dato = strings.TrimSpace(linea[a+1:])
		dato = dato[1 : len(dato)-1]
	}
	if len(linea) > a+3 {
		if a > 0 {
			campo = linea[:a]
			dato = strings.TrimSpace(linea[a+1:])
			dato = dato[1 : len(dato)-1]
		}
	}
	return campo, dato
}

//Files parses a directory of markdown files and converts them into Event
// types
func Files(dir string) ([]Parsed, error) {
	var lineas []string
	var title, fecha string

	var body string
	var inicio, fin bool = false, false
	var campo, dato, descripcion, tipos, imagen, slider, author, identificador, categorias, tags string
	events := []Parsed{}
	// Find event files in specified dir
	eventFiles, err := readMDFiles(dir)
	if err != nil {
		return nil, err
	}
	// Sort the files by the date in the title
	//	eventFiles, err = sortFilesChronological(eventFiles)

	for _, fichero := range eventFiles {
		file, err := os.Open(dir + "/" + fichero.filename) // abre el archivo
		scanner := bufio.NewScanner(file)                  // esto deberia escanearlo y analizar las lineas
		for scanner.Scan() {
			lineas = append(lineas, scanner.Text()) //va linea or linea añadiendo el contenido
		}
		if err != nil {
			panic(err.Error())

		}
	}
	for _, linea := range lineas {

		if inicio {
			campo, dato = analizalinea(linea)
			if linea == "---" {
				fin = true
				inicio = false
			} else {
				switch s.ToLower(campo) {
				case "title":
					title = dato

				case "date":
					fecha = dato

				case "description":
					descripcion = dato

				case "type":
					tipos = dato

				case "image":
					imagen = dato

				case "imageslider":
					slider = dato

				case "author":
					author = dato

				case "identifier":
					identificador = dato

				case "categories":
					categorias = dato

				case "tags":
					tags = dato
				}
			}
		}

		if fin {
			body = body + linea
		} else {
			if linea == "---" {
				inicio = true
			}
		}

	}
	event := Parsed{
		Title:         title,
		Date:          fecha,
		Description:   descripcion,
		Tipo:          tipos,
		Image:         imagen,
		Imageslide:    slider,
		Author:        author,
		Identificator: identificador,
		Categorias:    categorias,
		Tags:          tags,
		Body:          body,
	}
	events = append(events, event)

	return events, err
}
