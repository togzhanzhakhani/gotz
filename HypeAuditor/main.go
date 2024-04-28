package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("hypeauditor.com"),
	)

	file, err := os.Create("results.csv")
	if err != nil {
		log.Fatalf("Could not create file: %s", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Rank", "Username", "Name", "Category", "Followers", "Country", "Eng. Auth.", "Eng. Avg."})

	c.OnHTML(".row[data-v-b11c405a]", func(e *colly.HTMLElement) {
		rank := e.ChildText(".row-cell.rank")
		username := e.ChildText(".contributor__name")
		name := e.ChildText(".contributor__title")
		category := e.ChildTexts(".tag__content") 
		followers := e.ChildText(".row-cell.subscribers")
		country := e.ChildText(".row-cell.audience")
		engAuth := e.ChildText(".row-cell.authentic")
		engAvg := e.ChildText(".row-cell.engagement")

		joinedCategory := ""
		if len(category) > 0 {
			for _, cat := range category {
				joinedCategory += cat + "; "
			}
		}

		writer.Write([]string{rank, username, name, joinedCategory, followers, country, engAuth, engAvg})
	})

	c.Visit("https://hypeauditor.com/top-instagram-all-russia/")

	fmt.Println("Скрапинг завершен")
}
