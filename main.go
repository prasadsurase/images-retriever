package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "path/filepath"
  "sync"
  "time"

  "github.com/dchest/uniuri"
)

var urlsCount int

func main() {
  // root := "/Users/prasad/Downloads/Takeout_2/Google_Photos/"
  root := "/Users/prasad/Desktop/Lohagad"
  var wg sync.WaitGroup

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
      // fmt.Println("File:", path)
    case mode.IsRegular():
      ext := filepath.Ext(path)
      if ext == ".json" {
        files = append(files, path)
        wg.Add(1)
        go handleJSONFile(&wg, path)
      }
    }
    return nil
  })

  if err != nil {
    fmt.Println("Unable to open directory")
    log.Fatal(err)
  }
  wg.Wait()
  fmt.Println("Files count:", len(files))
  fmt.Println("Urls Retrieved:", urlsCount)
}

// parse the passed file and collect the urls.
func handleJSONFile(wg *sync.WaitGroup, path string) error {
  var data map[string]interface{}
  var fileWG sync.WaitGroup
  ext := filepath.Ext(path)
  fmt.Println("File:", path, ", Extension:", ext)
  fileData, err := ioutil.ReadFile(path)
  if err != nil {
    fmt.Println(err)
    log.Fatal(err)
    return err
  }
  defer wg.Done()
  json.Unmarshal(fileData, &data)
  fmt.Println("************************** JSON data *************************")
  fmt.Println(data)
  fileWG.Add(1)
  go parseData(&fileWG, data)
  fileWG.Wait()
  return nil
}

func parseData(fileWG *sync.WaitGroup, data map[string]interface{}) error {
  var parseWG sync.WaitGroup
  if data["url"] != nil && data["url"] != "" {
    fmt.Println("---------------------------==================---------------------------------")
    url := data["url"].(string)
    fmt.Println(url)
    parseWG.Add(1)
    go saveFile(&parseWG, url)
  }
  fileWG.Done()
  parseWG.Wait()
  return nil
}

func saveFile(parseWG *sync.WaitGroup, url string) error {
  urlsCount++
  fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
  fmt.Println("saveUrl: ", url)
  fmt.Println("start Http Get ``````````````````````````````")
  tr := &http.Transport{
    MaxIdleConns:       100,
    IdleConnTimeout:    30 * time.Second,
    DisableCompression: true,
  }
  client := &http.Client{Transport: tr}
  resp, err := client.Get(url)
  fmt.Println("done Http Get ``````````````````````````````")
  fmt.Println("Response:", resp)
  // if (err != nil) && (resp.StatusCode != 200) {
  // 	fmt.Println("Error:", err)
  // 	log.Fatal(err)
  // 	return err
  // }

  fmt.Println("Resp:", resp)
  defer resp.Body.Close()
  defer parseWG.Done()

  //open a file for writing
  // file, err := os.Create("/Users/prasad/Desktop/golang_images/" + uniuri.New() + ".jpg")
  // fmt.Println("Saved file:", *file)
  // defer file.Close()
  // if err != nil {
  // 	fmt.Println("Error:", err)
  // 	log.Fatal(err)
  // 	return err
  // }
  // fmt.Println("Created New File:", *file)
  // Use io.Copy to just dump the response body to the file. This supports huge files
  // size, err := io.Copy(file, resp.Body)
  data, err := ioutil.ReadAll(resp.Body)
  err = ioutil.WriteFile("/Users/prasad/Desktop/golang_images/"+uniuri.New()+".jpg", data, os.FileMode(0777))
  if err != nil {
    fmt.Println("Error:", err)
    log.Fatal(err)
    return err
  }
  fmt.Println()
  // fmt.Println("Saved file:", *file, ", Size:", size)
  return nil
}
