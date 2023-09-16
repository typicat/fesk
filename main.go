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

var oring, lax, harr int32
var fluga, spinn int32

func scrape() []Rapport {
	rap := make([]Rapport, 0)
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
				}
				if crap.Metod == "Spinn" {
					spinn++
				}
			case 8:
				crap.Langd = el.Text + "cm"
				strings.TrimSpace(crap.Langd)
			case 10:
				crap.Plats = el.Text
				strings.TrimSpace(crap.Plats)
			}
		})
		rap = append(rap, crap)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Fetching data... ")
	})
	c.Visit("https://kagealven.com/fangstrapporter-aktuella/")
	return rap
}

func summarize() {
	fmt.Println("Summarizing data...")
	fmt.Printf("Fluga %d - Spinn %d\n", fluga, spinn)
	fmt.Printf("Öringar %d - Laxar %d - Harrar %d\n", oring, lax, harr)
}

func main() {
	fmt.Printf("::: \033[34mfesk 0.1\033[0m - csv catch report from Kågeälven\n")
	forYear := time.Now()
	curYear := forYear.Year()
	outputFile := "rapport-" + fmt.Sprintf("%d", curYear) + ".csv"

	rap := scrape()
	fmt.Println("writing to file")

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
	fmt.Println("\033[32mDone!\033[0m")
	summarize()
}
