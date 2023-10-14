package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	colly "github.com/gocolly/colly/v2"
)

type Rapport struct {
	Datum      string
	Namn       string
	Art        string
	Kon        string
	Fettfena   string
	Aterutsatt string
	Utlekt     string
	Metod      string
	Langd      string
	Vikt       string
	Plats      string
	Kommentar  string
}

var (
	oring, lax, harr, fluga, spinn, aterutsatt                      int
	totalOringFluga, totalOringSpinn, totalLaxFluga, totalLaxSpinn int
)

func scrape() []Rapport {
	rapportData := make([]Rapport, 0)
	c := colly.NewCollector()
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		crap := Rapport{}
		e.ForEach("td", func(_ int, el *colly.HTMLElement) {
			switch el.Index {
			case 0:
				crap.Namn = strings.TrimSpace(el.Text)
			case 1:
				crap.Datum = el.Text
			case 2:
				crap.Art = el.Text
				if crap.Art == "Öring" {
					oring++
				}
				if crap.Art == "Lax" {
					lax++
				}
				if crap.Art == "Harr" {
					harr++
				}
			case 3:
				crap.Kon = el.Text
			case 4:
				crap.Fettfena = el.Text
			case 5:
				crap.Aterutsatt = el.Text
				aterutsatt++
			case 6:
				crap.Utlekt = el.Text
			case 7:
				crap.Metod = el.Text
				if crap.Metod == "Fluga" {
					fluga++
					if crap.Art == "Öring" {
						totalOringFluga++
					}
					if crap.Art == "Lax" {
						totalLaxFluga++
					}
				}
				if crap.Metod == "Spinn" {
					spinn++
					if crap.Art == "Öring" {
						totalOringSpinn++
					}
					if crap.Art == "Lax" {
						totalLaxSpinn++
					}
				}
			case 8:
				crap.Langd = strings.TrimSpace(el.Text) + "cm"
			case 9:
				crap.Vikt = strings.TrimSpace(el.Text)
			case 10:
				crap.Plats = el.Text
			case 11:
				crap.Kommentar = strings.TrimSpace(el.Text)
			}
		})
		rapportData = append(rapportData, crap)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("* Hämtar rapporter... ")
	})
	c.Visit("https://kagealven.com/fangstrapporter-aktuella/")
	return rapportData
}

func summarize(year int) {
	totalRapporterade := oring + lax + harr
	fmt.Println("Summering ", year, " - ", totalRapporterade, "st")
  fmt.Printf(" + Öring %d - Lax %d - Harr %d\n", oring, lax, harr)
	fmt.Println(" - Fluga ", fluga, "(Öring ", totalOringFluga, "Lax ", totalLaxFluga, ")")
	fmt.Println(" - Spinn ", spinn, "(Öring ", totalOringSpinn, "Lax ", totalLaxSpinn, ")")
	fmt.Println("Återutsatta: ", aterutsatt)
}

func main() {
	fmt.Printf("::: \033[34mfesk 0.1\033[0m - csv catch report from Kågeälven\n")
	yearToday := time.Now()
	curYear := yearToday.Year()
	outputFile := "rapport-" + fmt.Sprintf("%d", curYear) + ".csv"

	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	defer writer.Flush()
	headers := []string{
		"Datum",
		"Namn",
		"Art",
		"Kön",
		"Fettfena",
		"Återutsatt",
		"Utlekt",
		"Längd",
		"Metod",
		"Vikt",
		"Plats",
		"Kommentar",
	}
	err = writer.Write(headers)
	if err != nil {
		fmt.Println(err)
		return
	}

	rap := scrape()
	data := [][]string{}
	for _, v := range rap {
		data = append(
			data,
			[]string{
				v.Datum,
				v.Namn,
				v.Art,
				v.Kon,
				v.Fettfena,
				v.Aterutsatt,
				v.Utlekt,
				v.Langd,
				v.Metod,
				v.Vikt,
				v.Plats,
				v.Kommentar,
			},
		)
	}

	err = writer.WriteAll(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\033[32m* Done!\033[0m\n\n")
	summarize(curYear)
}
