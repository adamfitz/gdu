package main

import (
	"fmt"
	"os"

	"gdu/help"
	"gdu/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:\tgdu <directory_path>")
		fmt.Println("Help:\tgdu -h")
		return
	}
	if os.Args[1] != "-h" {
		rootPath := os.Args[1]

		if len(os.Args) > 2 {
			parser.DirResult(rootPath, os.Args[2])
		} else {
			parser.DirResult(rootPath, "")
		}
	} else {
		help.GduHelp()
	}

}
