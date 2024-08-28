package main

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/alecthomas/kong"
)

var CLI struct {
	File                  string   `name:"file" help:"Git patch file to modify." type:"path"`
	PathSubstituteMatcher []string `name:"psum" help:"Path substitution matcher." type:"string"`
	PathSubstituteReplace []string `name:"psur" help:"Path substitution substitute." type:"string"`
}

func handle_error(err error) {
	fmt.Println(fmt.Errorf("%v+", err))
	os.Exit(1)
}

func read_git_patch_file(handle *os.File, file string) string {
	if len(file) > 0 {
		fileHandle, err := os.Open(file)
		if err != nil {
			handle_error(err)
		}
		data, err := io.ReadAll(fileHandle)
		if err != nil {
			handle_error(err)
		}
		return string(data)
	}

	file_info, err := handle.Stat()
	if err != nil {
		handle_error(err)
	}
	if file_info.Mode()&os.ModeCharDevice == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			handle_error(err)
		}
		return string(data)
	}
	return ""
}

var path_points = [...]string{
	`^(diff --git a/)`,
	`^(diff --git a/.+? b/)`,
	`^(--- a/)`,
	`^(\+\+\+ b/)`,
}

func path_rename(data string, matchers []string, substitutions []string) string {
	for _, path_point := range path_points {
		for j, matcher := range matchers {
			substitution := substitutions[j]
			re := regexp.MustCompile("(?m)" + path_point + matcher)
			fmt.Println("(?m)" + path_point + matcher)
			data = re.ReplaceAllString(data, `${1}`+substitution)
			fmt.Println(`${1}` + substitution)
		}
	}
	return data
}

func main() {
	kong.Parse(&CLI)

	data := read_git_patch_file(os.Stdin, CLI.File)

	if len(CLI.PathSubstituteMatcher) > 0 {
		data = path_rename(data, CLI.PathSubstituteMatcher, CLI.PathSubstituteReplace)
	}
	fmt.Println(data)
}
