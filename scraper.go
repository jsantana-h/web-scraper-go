package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	// Importing Colly
	"github.com/gocolly/colly"
)

type PokemonProduct struct {
	url, image, name, price string
}

func main() {
	c := colly.NewCollector()

	var pokemonProducts []PokemonProduct

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	// iterating over the list of HTML product elements
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		// initializing a new PokemonProduct instance
		pokemonProduct := PokemonProduct{}

		// scraping the data of interest
		pokemonProduct.url = e.ChildAttr("a", "href")
		pokemonProduct.image = e.ChildAttr("img", "src")
		pokemonProduct.name = e.ChildText("h2")
		pokemonProduct.price = e.ChildText(".price")

		// adding the product instance with scraped data to the list of products
		pokemonProducts = append(pokemonProducts, pokemonProduct)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	// downloading the target HTML page
	c.Visit("https://scrapeme.live/shop/")

	// opening the CSV file
	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// initializing a file writer
	writer := csv.NewWriter(file)

	// defining the CSV headers
	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}
	// writing the column headers
	writer.Write(headers)

	// adding each Pokemon product to the CSV output file
	for _, pokemonProduct := range pokemonProducts {
		// converting a PokemonProduct to an array of strings
		record := []string{
			pokemonProduct.url,
			pokemonProduct.image,
			pokemonProduct.name,
			pokemonProduct.price,
		}

		// writing a new CSV record
		writer.Write(record)
	}
	defer writer.Flush()

	fmt.Println("Hello, World!")
}
