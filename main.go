package main

import (
	"fmt"
	"os"

	esa "esa_cli"
)

func main() {
	auth := os.Getenv("ESA_AUTH")
	team := os.Getenv("ESA_TEAM")
	client := esa.NewClient(auth, team)
	categories, _ := client.GetCategories()
	for _, category := range categories.Categories {
		for _, tree := range category.Tree() {
			fmt.Println(tree)
		}
	}

}
