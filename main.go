package main

import (
	"code-reviewer/internal/services/filewalk"
	"code-reviewer/internal/services/reviewer"
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {

	scandir := flag.String("scandir", ".", "Directory to scan")
	excludes := flag.String("exclude", "", "Comma-separated list of directories, files, extensions, or regex patterns to exclude")
	extensions := flag.String("extensions", "", "Comma-separated list of extensions used")
	keywords := flag.String("keywords", "", "Comma-separated list of keywords to filter files")

	flag.Parse()
	excludePatterns := strings.Split(*excludes, ",")
	extensionList := strings.Split(*extensions, ",")
	keywordList := strings.Split(*keywords, ",")

	filepaths, err := filewalk.FileFinder(*scandir, excludePatterns, extensionList, keywordList)
	if err != nil {
		log.Fatalf("Error while finding files: %v\n", err)
	}
	fmt.Printf("Found files:\n%s\n\n", strings.Join(filepaths, "\n"))
	fmt.Printf("Total number of files found: %d\n\n", len(filepaths))
	fmt.Println("Starting Review Code")
	reviewer.Review(filepaths)
}
