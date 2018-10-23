package main

import (
  "fmt"
  "log"
  "os"
  "path/filepath"
)

func main() {
  root := "/Users/prasad/Downloads/Takeout_2/Google_Photos/"
  // files, err := ioutil.ReadDir(root)
  // if err != nil {
  // 	fmt.Println("Unable to open directory")
  // 	log.Fatal(err)
  // }

  // for _, file := range files {
  // 	fmt.Println(file.Name())
  // }

  files := []string{}
  // filepath.Walk is the function which lists all the nested directories and the files in those directories.
  err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      log.Panicln("Something went wrong in: ", path)
      log.Fatal(err)
      return err
    }

    fi, err := os.Stat(path)
    if err != nil {
      fmt.Println(err)
      log.Fatal(err)
      return err
    }
    switch mode := fi.Mode(); {
    case mode.IsDir():
      fmt.Println("File:", path)
    case mode.IsRegular():
      ext := filepath.Ext(path)
      if ext == ".json" {
        files = append(files, path)
        go handleJSONFile(&path)
      }
    }
    return nil
  })

  if err != nil {
    fmt.Println("Unable to open directory")
    log.Fatal(err)
  }
  fmt.Println("Files count:", len(files))
}

func handleJSONFile(path *string) error {
  ext := filepath.Ext(*path)
  fmt.Println("File:", path, ", Extension:", ext)
  // fmt.Println("File:", *path)
  return nil
}
