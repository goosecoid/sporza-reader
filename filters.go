package main

import (
	"log"
	"os"
	"strings"
)

func getFilters() []string {
	f, err := os.ReadFile("filters.txt")
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(f), ",")
}
