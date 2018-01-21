package main

import (
	esa "esa_cli"
	"fmt"
	"os"
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
