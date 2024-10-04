package main

import (
	"code-reviewer/internal/services/filewalk"
	"flag"
	"fmt"
	"log"
	"strings"
)

func printHelp() {
	fmt.Println("Usage: code-reviewer [options]")
	fmt.Println("Options:")
	fmt.Println("  -scandir string")
	fmt.Println("        Directory to scan (default \".\")")
	fmt.Println("  -excludes string")
	fmt.Println("        Comma-separated list of directories, files, extensions, or regex patterns to exclude")
	fmt.Println("  -extensions string")
	fmt.Println("        Comma-separated list of extensions used")
	fmt.Println("  -keywords string")
	fmt.Println("        Comma-separated list of keywords to filter files")
}

func main() {

	scandir := flag.String("scandir", ".", "Directory to scan")
	excludes := flag.String("excludes", "", "Comma-separated list of directories, files, extensions, or regex patterns to exclude")
	extensions := flag.String("extensions", "", "Comma-separated list of extensions used")
	keywords := flag.String("keywords", "", "Comma-separated list of keywords to filter files")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()
	if *help {
		printHelp()
		return
	}
	var excludePatterns []string
	if *excludes != "" {
		excludePatterns = strings.Split(*excludes, ",")
	}
	var extensionList []string
	if *extensions != "" {
		extensionList = strings.Split(*extensions, ",")
	}

	var keywordList []string
	if *keywords != "" {
		keywordList = strings.Split(*keywords, ",")
	}

	filepaths, err := filewalk.FileFinder(*scandir, excludePatterns, extensionList, keywordList)
	if err != nil {
		log.Fatalf("Error while finding files: %v\n", err)
	}
	fmt.Printf("Found files:\n%s\n\n", strings.Join(filepaths, "\n"))
	fmt.Printf("Total number of files found: %d\n\n", len(filepaths))
	fmt.Println("Starting Review Code")
	// reviewer.Review(filepaths)
}
