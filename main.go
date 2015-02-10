package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// array of regex to execute in order
// {{re, replace string}}
var lessToSassReplacePairs ReplacePairs = ReplacePairs{
	ReplacePair{`@(?!import|media|keyframes|-)`, `$`},
	ReplacePair{`\.([\w\-]*)\s*\((.*)\)\s*\{`, `@mixin \1\(\2\)\n{`},
	ReplacePair{`\.([\w\-]*\(.*\)\s*;)`, `@include \1`},
	ReplacePair{`~"(.*)"`, `#{"\1"}`},
	ReplacePair{`spin`, `adjust-hue`},
}

func transformLessToSass(content []byte) []byte {
	return Replacer(lessToSassReplacePairs).Replace(content)
}

func parseSrc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() && filepath.Ext(path) == ".less" {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println("There was an error reading file", path)
			return err
		}

		// write file into destination directory
		destPath := os.Args[len(os.Args)-1]
		addSuffixIfMissing(&destPath, "/")

		newFilePath := destPath + replaceExt(path, ".scss")

		err = ioutil.WriteFile(newFilePath, transformLessToSass(content), info.Mode())
		if err != nil {
			log.Println("there was an error writing to", newFilePath)
			return err
		}
	}
	return err
}

func main() {
	lastArgIndex := IntMax(len(os.Args)-1, 1)

	if lastArgIndex < 2 {
		log.Fatal("USAGE: less2sass <srcFile or srcDirectory> ... <destDirectory>")
	}

	// check if the last argument is correct
	if _, err := dirExists(os.Args[lastArgIndex]); err != nil {
		log.Fatal("The last argument should be the destination directory")
	}

	// walk through the source files/directories
	for _, filePath := range os.Args[1:lastArgIndex] {
		if err := filepath.Walk(filePath, parseSrc); err != nil {
			log.Println("Could not walk through source directory", filePath)
			log.Println(err)
		}
	}
}
