package main

import (
  "os"
  "strings"
  "io/ioutil"
  "path/filepath"
  "github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

// array of regex to execute in order
// {{re, replace string}}
var transformRegexArray [][]string = [][]string{
  {`@(?!import|media|keyframes|-)`, `$`},
  {`\.([\w\-]*)\s*\((.*)\)\s*\{`, `@mixin \1\(\2\)\n{`},
  {`\.([\w\-]*\(.*\)\s*;)`, `@include \1`},
  {`~"(.*)"`, `#{"\1"}`},
  {`spin`, `adjust-hue`},
}

func transformLessToSass(content []byte) []byte {
  newContent := content
  for _, regexArray := range transformRegexArray {
    re := pcre.MustCompile(regexArray[0], 0)
    newContent = re.ReplaceAll(newContent, []byte(regexArray[1]), 0)
  }
  return newContent
}

func parseSrc(path string, info os.FileInfo, err error) error {
  if !info.IsDir() && filepath.Ext(path) == ".less" {
    content, err := ioutil.ReadFile(path)
    if err != nil {
      println("There was an error reading file", path)
      return nil
    }
    lessContent := content

    // write file into destination directory
    destPath := os.Args[len(os.Args) - 1]
    if !strings.HasSuffix(destPath, "/") {
      destPath = destPath + "/"
    }
    newFilePath := strings.Join([]string{destPath, strings.TrimSuffix(path, filepath.Ext(path)), ".scss"}, "")

    err = ioutil.WriteFile(newFilePath, transformLessToSass(lessContent), info.Mode())

    if err != nil {
      println("there was an error writing the sass file")
      return err
    }
  }
  return nil
} 

func main() {
  lastArgIndex := IntMax(len(os.Args) - 1, 1)

  if lastArgIndex < 2 {
    println("USAGE: less2sass <srcFile or srcDirectory> ... <destDirectory>")
    return
  }
  
  // check if the last argument is correct
  destDir, err := os.Open(os.Args[lastArgIndex]); 
  if err != nil {
    println("The last argument should be the destination directory")
    return
  }

  if stat, err := destDir.Stat(); err != nil || !stat.IsDir() {
    println("The destination argument should point to a directory")
    return
  }

  // walk through the source files/directories
  for _, filePath := range os.Args[1:lastArgIndex] {
    if filepath.Walk(filePath, parseSrc) != nil {
      println ("The ")
    }
  }
}

func IntMax(a int, b int) int {
  if (a > b) {
    return a
  } else {
    return b
  }
}

