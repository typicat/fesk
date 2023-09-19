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
	Datum string
	Namn  string
	Plats string
	Art   string
	Langd string
	Metod string
}

var oring, lax, harr int
var fluga, spinn, totalOringFluga, totalOringSpinn, totalLaxFluga, totalLaxSpinn int

func scrape() []Rapport {
	rapportData := make([]Rapport, 0)
	c := colly.NewCollector()
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		crap := Rapport{}
		e.ForEach("td", func(_ int, el *colly.HTMLElement) {
			switch el.Index {
			case 0:
				crap.Namn = el.Text
				strings.TrimSpace(crap.Namn)
			case 1:
				crap.Datum = el.Text
				strings.TrimSpace(crap.Datum)
			case 2:
				crap.Art = el.Text
				strings.TrimSpace(crap.Art)
				if crap.Art == "Öring" {
					oring++
				}
				if crap.Art == "Lax" {
					lax++
				}
				if crap.Art == "Harr" {
					harr++
				}
			case 7:
				crap.Metod = el.Text
				strings.TrimSpace(crap.Metod)
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
				crap.Langd = el.Text + "cm"
				strings.TrimSpace(crap.Langd)
			case 10:
				crap.Plats = el.Text
				strings.TrimSpace(crap.Plats)
			}
		})
		rapportData = append(rapportData, crap)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Fetching data... ")
	})
	c.Visit("https://kagealven.com/fangstrapporter-aktuella/")
	return rapportData
}

func summarize(year int) {
	fmt.Println("::: Summary", year)
	fmt.Printf("- Öring\tFluga %d | Spinn %d\n", totalOringFluga, totalOringSpinn)
	fmt.Printf("- Lax\tFluga %d | Spinn %d\n\n", totalLaxFluga, totalLaxSpinn)
	fmt.Println("::: Total")
	fmt.Printf("Öringar %d \t Laxar %d \t Harrar %d\n", oring, lax, harr)
	fmt.Printf("Fluga %d \t Spinn %d\n", fluga, spinn)
}

func main() {
	fmt.Printf("::: \033[34mfesk 0.1\033[0m - csv catch report from Kågeälven\n")
	forYear := time.Now()
	curYear := forYear.Year()
	outputFile := "rapport-" + fmt.Sprintf("%d", curYear) + ".csv"

	rap := scrape()
	fmt.Printf("writing to file... ")
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()
	headers := []string{"Datum", "Namn", "Art", "Längd", "Metod", "Plats"}
	err = writer.Write(headers)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := [][]string{}
	for _, v := range rap {
		data = append(data, []string{v.Datum, v.Namn, v.Art, v.Langd, v.Metod, v.Plats})
	}
	err = writer.WriteAll(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\033[32mDone!\033[0m\n\n")
	summarize(curYear)
}
