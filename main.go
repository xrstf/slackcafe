package main

import (
	"io/ioutil"
	"log"
)

func main() {
	// body, err := fetchMenuImage()
	// if err != nil {
	// 	log.Fatalf("Failed to determine menu image URL: %v", err)
	// }
	// defer body.Close()

	// tmpFile := "foo.png"

	// err = prepareImage(body, tmpFile)
	// if err != nil {
	// 	log.Fatalf("Failed to prepare image: %v", err)
	// }

	// text, err := ocr(tmpFile)
	// if err != nil {
	// 	log.Fatalf("Failed to OCR image: %v", err)
	// }

	text, err := ioutil.ReadFile("ocr.txt")
	if err != nil {
		log.Fatalf("Failed to read test file: %v", err)
	}

	karte, err := parseMenu(string(text))
	if err != nil {
		log.Fatalf("Failed to parse menu: %v", err)
	}

	printMenu(karte)
}

func printMenu(menu []dailyMenu) {
	for idx, m := range menu {
		log.Printf("Day: %d", idx+1)
		log.Printf("MI  %s", m.M1.title)
		log.Printf("    %-50s %s", m.M1.subtitle, m.M1.price)
		log.Printf("MII %s", m.M2.title)
		log.Printf("    %-50s %s", m.M1.subtitle, m.M2.price)
		log.Println("")
	}
}
