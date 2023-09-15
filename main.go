package main

import (
	"fmt"
	"os"
	"strings"

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

func scrape() []Rapport {
	rap := make([]Rapport, 0)
	c := colly.NewCollector()
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		crap := Rapport{}
		e.ForEach("td", func(_ int, el *colly.HTMLElement) {
			switch el.Index {
			case 0:
				crap.Namn = el.Text + ","
				strings.TrimSpace(crap.Namn)
			case 1:
				crap.Datum = el.Text + ","
				strings.TrimSpace(crap.Datum)
			case 2:
				crap.Art = el.Text + ","
				strings.TrimSpace(crap.Art)
			case 7:
				crap.Metod = el.Text + ","
				strings.TrimSpace(crap.Metod)
			case 8:
				crap.Langd = el.Text + "cm,"
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

func main() {
	fmt.Printf("::: \033[34mfesk 0.1\033[0m - csv catch report from Kågeälven\n")
	rap := scrape()
	fmt.Println("writing to file")
	f, err := os.Create("rapport.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	for _, v := range rap {
		fmt.Fprintln(f, v.Datum, v.Namn, v.Art, v.Langd, v.Metod, v.Plats)
	}
	fmt.Println("\033[32mDone!\033[0m")

}
