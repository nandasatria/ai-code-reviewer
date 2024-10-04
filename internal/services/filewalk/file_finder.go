package filewalk

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func compilePatterns(patterns []string) ([]*regexp.Regexp, error) {
	var regexps []*regexp.Regexp
	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern %q: %v", pattern, err)
		}
		regexps = append(regexps, re)
	}
	return regexps, nil
}

func isMatchExclude(path string, excludePatterns []*regexp.Regexp) bool {
	for _, re := range excludePatterns {
		if re.MatchString(path) {
			return true
		}
	}
	return false
}

func FileFinder(dir string, exclude []string, extensions []string, keywords []string) ([]string, error) {
	var result []string
	extensionSet := make(map[string]struct{}, len(extensions))

	for _, ext := range extensions {
		extensionSet[ext] = struct{}{}
	}

	excludeRegexps, err := compilePatterns(exclude)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return fmt.Errorf("error walking the path %q: %v", path, err)
		}

		if info.IsDir() || isMatchExclude(path, excludeRegexps) {
			return nil
		}

		foundKeywords := false
		for _, keyword := range keywords {
			if strings.Contains(path, keyword) {
				foundKeywords = true
				break
			}
		}

		if len(keywords) > 0 && !foundKeywords {
			return nil
		}

		ext := filepath.Ext(path)
		if _, match := extensionSet[ext]; !match && len(extensions) > 0 {
			return nil
		}

		result = append(result, path)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error occurred while scanning the directory %q: %v", dir, err)
	}
	return result, nil
}
