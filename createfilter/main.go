package createfilter

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

// Usage: $ ./create_filter --path test_data.csv

// DefaultNbMois is the default number of the most recent months during which the effectif of the company must reach the threshold.
const DefaultNbMois = 100

// DefaultMinEffectif is the default effectif threshold, expressed in number of employees.
const DefaultMinEffectif = 10

// DefaultNbIgnoredRecords is the default number of most recent months to ignore from the effectif file.
const DefaultNbIgnoredRecords = 2

// Implementation of the prepare-import command.
func main() {

	var path = flag.String("path", "", "Chemin d'accès au fichier effectif")
	var nbMois = flag.Int(
		"nbMois",
		DefaultNbMois,
		"Nombre de mois observés (avec effectif connu) pour déterminer si l'entreprise dépasse 10 salariés",
	)
	var minEffectif = flag.Int(
		"minEffectif",
		DefaultMinEffectif,
		"Si une entreprise atteint ou dépasse 'minEffectif' dans les 'nbMois' derniers mois, elle est inclue dans le périmètre du filtre.",
	)
	var nIgnoredRecords = flag.Int(
		"nIgnoredRecords",
		DefaultNbIgnoredRecords,
		"Nombre de colonnes à ignorer à la fin du fichier effectif",
	)
	flag.Parse()

	err := CreateFilter(os.Stdout, *path, *nbMois, *minEffectif, *nIgnoredRecords)
	if err != nil {
		log.Panic(err)
	}
}

// CreateFilter generates a "filter" from an "effectif" file.
func CreateFilter(writer io.Writer, effectifFileName string, nbMois, minEffectif int, nIgnoredRecords int) error {
	last := guessLastNMissing(effectifFileName, nIgnoredRecords)
	f, err := os.Open(effectifFileName)
	defer f.Close()
	if err != nil {
		return err
	}
	r := initializeEffectifReader(f)
	outputPerimeter(r, writer, nbMois, minEffectif, nIgnoredRecords+last)
	return nil
}

func initializeEffectifReader(f *os.File) *csv.Reader {
	r := csv.NewReader(bufio.NewReader(f))
	r.LazyQuotes = true
	r.Comma = ';'
	return r
}

func outputPerimeter(r *csv.Reader, w io.Writer, nbMois, minEffectif, nIgnoredRecords int) {
	perimeter := []string{}
	_, err := r.Read() // en tête
	if err != nil {
		log.Panic(err)
	}
	for {
		record, err := r.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}
		shouldKeep := len(record[1]) == 14 &&
			isInsidePerimeter(record[0:len(record)-nIgnoredRecords], nbMois, minEffectif)

		if shouldKeep {
			perimeter = append(perimeter, record[1])
		}
	}
	for _, siret := range perimeter {
		fmt.Fprintln(w, siret[0:9])
	}
}

func isInsidePerimeter(record []string, nbMois, minEffectif int) bool {
	for i := len(record) - 1; i >= len(record)-nbMois && i >= 0; i-- {
		if record[i] == "" {
			continue
		}
		reg, err := regexp.Compile("[^0-9]")
		if err != nil {
			log.Fatal(err)
		}
		effectif, err := strconv.Atoi(reg.ReplaceAllString(record[i], ""))
		if err != nil {
			fmt.Println(record)
			log.Panic(err)
		}
		if effectif >= minEffectif {
			return true
		}
	}
	return false
}

func guessLastNMissing(path string, nIgnoredRecords int) int {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Panic(err)
	}
	r := initializeEffectifReader(f)
	_, err = r.Read() // en tête
	if err != nil {
		log.Panic(err)
	}
	return guessLastNMissingFromReader(r, nIgnoredRecords)
}

func guessLastNMissingFromReader(r *csv.Reader, nIgnoredRecords int) int {
	lastNonMissing := -1
	var recordLength int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}
		recordLength = len(record)
		for i := len(record) - 1 - nIgnoredRecords; i > lastNonMissing; i-- {
			if record[i] != "" {
				lastNonMissing = i
				continue
			}
		}
	}
	return recordLength - 1 - nIgnoredRecords - lastNonMissing // index
}
