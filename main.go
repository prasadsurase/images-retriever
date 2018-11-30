package main

import (
  "encoding/json"
  "fmt"
  "io"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "path/filepath"
  "strings"
  "sync"
  "time"

  "github.com/dchest/uniuri"
)

var imgUrlsCount int
var videoUrlsCount int
var failedRespUrlsCount int

// accept a folder path that consists of json files. The json files contain the 'url' to the images/videos. Use these
// urls to download the images/videos.
func main() {
  root := "/Users/prasad/Downloads/Takeout/Google_Photos/"
  // root := "/Users/prasad.surase/Desktop/Lohagad"
  var wg sync.WaitGroup

  //All json files found in the specified folder.
  files := []string{}
  // filepath.Walk is the function which lists all the nested directories and the files in those directories.
  err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
    // error incase the provided folder path doesnt exist
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
      // if the file has json extension, add it to the list of files and send it for parsing.
      if ext == ".json" {
        files = append(files, path)
        wg.Add(1)
        go handleJSONFile(&wg, path)
        time.Sleep(500 * time.Millisecond)
      }
    }
    return nil
  })

  if err != nil {
    fmt.Println("Unable to open directory")
    log.Fatal(err)
  }
  wg.Wait()
  // display summary
  fmt.Println("Files count:", len(files))
  fmt.Println("Urls Retrieved:", imgUrlsCount)
  fmt.Println("Vides Urls:", videoUrlsCount)
  fmt.Println("Failed Response Urls: ", failedRespUrlsCount)
}

// parse the passed file and collect the urls.
func handleJSONFile(wg *sync.WaitGroup, path string) error {
  var data map[string]interface{}
  var fileWG sync.WaitGroup
  ext := filepath.Ext(path)
  fmt.Println("File:", path, ", Extension:", ext)

  // parse the file using the provided path
  fileData, err := ioutil.ReadFile(path)
  if err != nil {
    fmt.Println(err)
    log.Fatal(err)
    return err
  }
  defer wg.Done()
  json.Unmarshal(fileData, &data)
  // fmt.Println("************************** JSON data *************************")
  // fmt.Println(data)
  fileWG.Add(1)

  // call goroutine to parse the data(incase the data is too large)
  go parseData(&fileWG, data)
  fileWG.Wait()
  return nil
}


// parse the json data. As of now, we only take the 'url' key-value from parsed data. Incase the data is array, we need
// to gather all the 'url' key-value pairs.
func parseData(fileWG *sync.WaitGroup, data map[string]interface{}) error {
  var parseWG sync.WaitGroup
  if data["url"] != nil && data["url"] != "" {
    url := data["url"].(string)
    parseWG.Add(1)

    // call go routine to download the image/video using the url
    go saveFile(&parseWG, url)
    time.Sleep(600 * time.Millisecond)
  }
  fileWG.Done()
  parseWG.Wait()
  return nil
}

// download the image/video using the provided url
func saveFile(parseWG *sync.WaitGroup, url string) error {
  if strings.Contains(url, "video-downloads") {
    fmt.Println("Video url: ", url)
    videoUrlsCount++
    return nil
  }
  tr := &http.Transport{
    MaxIdleConns:       10,
    IdleConnTimeout:    30 * time.Second,
    DisableCompression: true,
  }
  client := &http.Client{Transport: tr}
  // fmt.Println("Url: ", url)

  //get the data using url
  resp, err := client.Get(url)
  // fmt.Println("Resp:", resp)

  if resp == nil {
    fmt.Println("Response was nil for url:", url)
    failedRespUrlsCount++
    return nil
  }

  if (err != nil) || (resp.StatusCode != http.StatusOK) {
    fmt.Println("Error:", err)
    fmt.Println("Url:", url)
    log.Fatal(err)
    return nil
  }

  defer resp.Body.Close()
  defer parseWG.Done()

  // open file for writing and dump response data in the newly created file.
  file, err := os.Create("/Users/prasad/Desktop/golang_images/" + uniuri.New() + ".jpg")
  // fmt.Println("Saved file:", *file)
  if err != nil {
    fmt.Println("Error:", err)
    log.Fatal(err)
    return nil
  }
  defer file.Close()

  size, err := io.Copy(file, resp.Body)
  if err != nil {
    fmt.Println("Error:", err)
    log.Fatal(err)
    return nil
  }
  imgUrlsCount++
  time.Sleep(50 * time.Millisecond)
  fmt.Println("Saved file:", file.Name(), ", Size:", size)
  return nil
}
