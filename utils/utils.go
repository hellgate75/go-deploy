package utils

import (
	"github.com/hellgate75/go-deploy/io"
	"path/filepath"
	"runtime"
	"strings"
)

// Gets the Shared libraries extension included by dot, related to current O/S
func GetShareLibExt() string {
	if runtime.GOOS == "windows" {
		return ".dll"
	}
	return ".so"
}

// Gets the Shared libraries extension included by dot, related to current O/S
func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

func FixFolder(folder string, baseDir string, defaultdirName string) string {
	if folder == "" {
		var path string = io.GetPathSeparator() + defaultdirName
		if defaultdirName == "" {
			path = ""
		}
		folder = baseDir + path
	} else if len(folder) < 3 || folder[:2] == "./" || folder[:2] == ".\\" ||
		folder[:3] == "../" || folder[:3] == "..\\" {
		folder = baseDir + io.GetPathSeparator() + folder
	}

	return filepath.Clean(folder)
}

func SliceContains(list []interface{}, item interface{}) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func SliceAppend(list []interface{}, secondList []interface{}) []interface{} {
	return append(list, secondList...)
}

func SliceUnique(list []interface{}) []interface{} {
	var out []interface{} = make([]interface{}, 0)
	var seen map[interface{}]bool = make(map[interface{}]bool)
	for _, item := range list {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}

func StringSliceContains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func StringSliceAppend(list []string, secondList []string) []string {
	return append(list, secondList...)
}

func StringSliceUnique(list []string) []string {
	var out []string = make([]string, 0)
	var seen map[string]bool = make(map[string]bool)
	for _, item := range list {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}

func StringSliceTrim(list []string) []string {
	var out []string = make([]string, 0)
	for _, item := range list {
		if strings.TrimSpace(item) == "" {
			continue
		}
		out = append(out, item)
	}
	return out
}
