package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func Print(title string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal("Could not print")
	}

	fmt.Println(title+":", string(b))
	fmt.Println()
}
