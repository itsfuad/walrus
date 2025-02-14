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
		os.Exit(-1)
	}

	filePath := os.Args[1]

	//must have .wal file
	if len(filePath) < 5 || filePath[len(filePath)-4:] != ".wal" {
		colors.RED.Println("Error: file must have .wal extension")
		os.Exit(-1)
	}

	analyzer.Analyze(filePath, false, false)
}
