package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var (
		err error
	)

	commandArgs := CheckCommandLineArgs(os.Args)
	if !commandArgs {
		LogMessage("Usage: searchforfilepatterns.exe --file <path to file> --ext <file extention ex. *.out.log> \n --modifieddate <file modified format: 2022-02-15> --pattern <string pattern to seach for>")
		return
	}

	pathToFiles := os.Args[2]
	fileExtenstion := os.Args[4]
	modifiedDate := os.Args[6]
	filePattern := os.Args[8]

	directoryFiles, _ := WalkMatch(pathToFiles, fileExtenstion, modifiedDate)
	if err != nil {
		panic(err)
	}

	for _, directoryfile := range directoryFiles {
		matchingStrings := SearchFilesForStringPattern(directoryfile, filePattern)
		if err != nil {
			panic(err)
		}
		if len(matchingStrings) > 0 {
			LogMessage(fmt.Sprintf("Found matches with the extension provided in the modified date range for directory %v. The file name is %v", pathToFiles, path.Base(directoryfile)))
			LogMessage("Matches are listed below.")
			for _, lines := range matchingStrings {
				LogMessage(lines)
			}
		}
	}
}

func LogMessage(message string) {
	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	fullMessage := timeStamp + ":  " + message
	fmt.Println(fullMessage)
}

func CheckCommandLineArgs(args []string) bool {
	if len(args) == 9 {
		if strings.ToLower(args[1]) == "--file" {
			if strings.ToLower(args[3]) == "--ext" {
				if strings.ToLower(args[5]) == "--modifieddate" {
					if strings.ToLower(args[7]) == "--pattern" {
						return true
					}
				}
			}
		}
	}
	return false
}

func WalkMatch(root, pattern string, fileModDate string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			if GetModTimeForFile(path, fileModDate) {
				matches = append(matches, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func GetModTimeForFile(path string, fileModDate string) bool {
	getfileStats, _ := os.Stat(path)
	getFilesLastModDate := getfileStats.ModTime().Format("2006-01-02")
	if getFilesLastModDate >= fileModDate {
		return true
	} else {
		return false
	}
}

func SearchFilesForStringPattern(path string, pattern string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	lineCounter := 0
	for scanner.Scan() {
		lineCounter++
		if strings.Contains(scanner.Text(), pattern) {
			text = append(text, fmt.Sprintf("Line %v: %v", lineCounter, scanner.Text()))
		}
	}
	file.Close()

	return text
}
