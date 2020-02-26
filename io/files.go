package io

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

// Gets the Path Separator as string type
func GetPathSeparator() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return string(os.PathSeparator)
}

// Gets the Path Separator as string type
func GetShareLibExt() string {
	if runtime.GOOS == "windows" {
		return ".dll"
	}
	return ".so"
}

//Verify if a atring file path corresponds to a directory
func IsFolder(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// Get files in a folder (eventually recursively)
func GetFiles(path string, recursive bool) []string {
	var out []string = make([]string, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return out
	}
	for _, file := range files {
		var name = path + GetPathSeparator() + file.Name()
		if !file.IsDir() {
			out = append(out, name)
		} else if recursive {
			var filesX []string = GetFiles(name, recursive)
			for _, fileX := range filesX {
				out = append(out, fileX)
			}
		}
	}
	return out
}

// Get files in a folder (eventually recursively), which name matches with given function execution
func GetMatchedFiles(path string, recursive bool, matcher func(string) bool) []string {
	var out []string = make([]string, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return out
	}
	for _, file := range files {
		var name = path + GetPathSeparator() + file.Name()
		if !file.IsDir() {
			if matcher(name) {
				out = append(out, name)
			}
		} else if recursive {
			var filesX []string = GetMatchedFiles(name, recursive, matcher)
			for _, fileX := range filesX {
				if matcher(name) {
					out = append(out, fileX)
				}
			}
		}
	}
	return out
}

func FindFilesIn(path string, recursive bool, searchText string) []string {
	return GetMatchedFiles(path, recursive, func(name string) bool {
		return strings.Contains(name, searchText)
	})
}
