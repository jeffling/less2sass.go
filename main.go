package main

import (
  "os"
  "log"
  "strings"
  "io/ioutil"
  "path/filepath"
)

// array of regex to execute in order
// {{re, replace string}}
var lessToSassReplacePairs ReplacePairs = ReplacePairs{
  []ReplacePair {
    ReplacePair{`@(?!import|media|keyframes|-)`, `$`},
    ReplacePair{`\.([\w\-]*)\s*\((.*)\)\s*\{`, `@mixin \1\(\2\)\n{`},
    ReplacePair{`\.([\w\-]*\(.*\)\s*;)`, `@include \1`},
    ReplacePair{`~"(.*)"`, `#{"\1"}`},
    ReplacePair{`spin`, `adjust-hue`},
  },
}

func transformLessToSass(content []byte) []byte {
  return Replacer(lessToSassReplacePairs).Replace(content)
}

func addSuffixIfMissing(s *string, suffix string) {
  if !strings.HasSuffix(*s, suffix) {
    *s = *s + suffix
  }
}

func replaceExt(path string, newExt string) string {
  return strings.TrimSuffix(path, filepath.Ext(path)) + newExt
}

func parseSrc(path string, info os.FileInfo, err error) error {
  if err != nil {
    return err
  }

  if !info.IsDir() && filepath.Ext(path) == ".less" {
    content, err := ioutil.ReadFile(path)
    if err != nil {
      println("There was an error reading file", path)
      return err
    }

    // write file into destination directory
    destPath := os.Args[len(os.Args) - 1]
    addSuffixIfMissing(&destPath, "/")
    
    newFilePath := destPath + replaceExt(path, ".scss")

    err = ioutil.WriteFile(newFilePath, transformLessToSass(content), info.Mode())

    if err != nil {
      println("there was an error writing to", newFilePath)
      return err
    }
  }
  return err
} 

func main() {
  lastArgIndex := IntMax(len(os.Args) - 1, 1)

  if lastArgIndex < 2 {
    log.Fatal("USAGE: less2sass <srcFile or srcDirectory> ... <destDirectory>")
  }
  
  // check if the last argument is correct
  destDir, err := os.Open(os.Args[lastArgIndex]); 
  if err != nil {
    log.Fatal("The last argument should be the destination directory")
  }
  defer destDir.Close()

  if stat, err := destDir.Stat(); err != nil || !stat.IsDir() {
    log.Fatal("The destination argument should point to a directory")
  }

  // walk through the source files/directories
  for _, filePath := range os.Args[1:lastArgIndex] {
    if filepath.Walk(filePath, parseSrc) != nil {
      log.Fatal("Could not walk through source directories")
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

