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
	"time"

	"gopkg.in/russross/blackfriday.v2"
)

var filePattern = regexp.MustCompile(`\.md$`)

type frontmatter struct {
	title string
}

type mdFile struct {
	filename string
	bytes    []byte
}

// Parsed represents a single parsed file
type Parsed struct {
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Tipo        string    `json:"tipo"`
	Image       string    `json:"image"`
	Imageslide  string    `json:"imageslide"`
	Categories  string    `json:"categories"`
	Tags        string    `json:"tags"`
	Body        string    `json:"body"`
}

func readMDFiles(dir string) ([]mdFile, error) {
	mdFiles := []mdFile{}
	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	files, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, n := range files {
		if filePattern.Match([]byte(n)) {
			f, err := ioutil.ReadFile(filepath.Join(dir, n)) //ioutil para entrada y salida de datos
			if err != nil {
				log.Printf("Cannot read file %s/%s", dir, n)
				continue
			}
			// Extract slug and date from filename
			//			filenameParts := filePattern.FindAllStringSubmatch(n, -1)
			//dateStr := filenameParts[0][2]
			//slugStr := filenameParts[0][1]
			//d, err := time.Parse("2006-01-02", dateStr)
			newFile := mdFile{
				filename: n,
				bytes:    f}
			mdFiles = append(mdFiles, newFile) //almacena los datos de los ficheros en un array
		}
	}
	return mdFiles, nil
}

/* func extractYAMLFrontmatter(body []byte) (map[string]string, string, error) {
	frontmatterPattern := regexp.MustCompile(`---\n(.*: .*\n)+---`)
	bodyString := frontmatterPattern.ReplaceAllString(string(body), "")
	frontmatterString := frontmatterPattern.Find(body)
	plainYAMLString := strings.Replace(string(frontmatterString), "---", "", 2)
	parsedYAML := make(map[string]string)
	err := yaml.Unmarshal([]byte(plainYAMLString), &parsedYAML)
	if err != nil {
		return nil, bodyString, err
	}
	return parsedYAML, bodyString, nil
} */

/* func sortFilesChronological(f []mdFile) ([]mdFile, error) {
	fSorted := make([]mdFile, len(f))
	copy(fSorted, f)
	sort.Slice(fSorted, func(i, j int) bool { return fSorted[i].date.After(fSorted[j].date) })
	return fSorted, nil
}
*/
func parseBodyHTML(b []byte) []byte {
	// Custom img tag
	imgTagPattern := regexp.MustCompile(`(?im)\%img\[(.*)\]\((.*)\)`)
	b = imgTagPattern.ReplaceAll(b, []byte("<img src=\"$2\" alt=\"$1\" /><div class=\"img-caption\">$1</div>"))
	// Render standard markdown
	bodyHTML := blackfriday.Run(b)
	return bodyHTML
}

func analizalinea(linea string) (string, string) {

	var campo, dato string

	a := s.Index(linea, ":")
	if a > 0 {
		campo = linea[:a]
		dato = strings.TrimSpace(linea[a+2:])
		dato = dato[2 : len(dato)-1]
	}

	return campo, dato
}

//Files parses a directory of markdown files and converts them into Event
// types
func Files(dir string) ([]Parsed, error) {
	var lineas []string
	var descripcion []string
	var inicio bool = false
	var campo, dato string
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

		for _, linea := range lineas {
			if inicio {
				campo, dato = analizalinea(linea)
				if campo == "title" {
					descripcion = append(descripcion, dato)
				} else {
					if campo == "imagen" {
					}
				}
			}
			if linea == "---" {
				inicio = true
			}
		}
		if err != nil {
			return nil, err
		}
		defer file.Close()
		event := Parsed{
			Title: fichero.filename,
			//Date:        time.Time,
			Description: fichero.filename,
			Tipo:        fichero.filename,
			Image:       fichero.filename,
			Imageslide:  fichero.filename,
			Categories:  fichero.filename,
			Tags:        fichero.filename,
			Body:        fichero.filename,
		}
		events = append(events, event)
	}
	return nil, nil
}
