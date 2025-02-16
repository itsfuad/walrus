package main

import (
	//Standard packages
	"os"

	//Walrus packages
	"walrus/compiler/analyzer"
	"walrus/compiler/colors"
)

func main() {

	if len(os.Args) < 2 {
		colors.GREEN.Println("Usage: walrus <file>")
		return
	}

	filePath := os.Args[1]

	r, err := analyzer.Analyze(filePath, true, false, false)
	if len(r) > 0 {
		r.DisplayAll()
	}

	if err != nil {
		colors.RED.Println("Error analyzing file: ", err)
	}
}
