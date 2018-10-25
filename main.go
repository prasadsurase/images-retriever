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

	"github.com/dchest/uniuri"
)

func main() {
	root := "/Users/prasad/Downloads/Takeout_2/Google_Photos/"

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

// parse the passed file and collect the urls.
func handleJSONFile(path *string) error {
	var data map[string]interface{}
	ext := filepath.Ext(*path)
	fmt.Println("File:", *path, ", Extension:", ext)
	fileData, err := ioutil.ReadFile(*path)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return err
	}
	json.Unmarshal(fileData, &data)
	fmt.Println("************************** JSON data *************************")
	fmt.Println(data)
	go parseData(data)
	return nil
}

func parseData(data map[string]interface{}) error {
	if data["url"] != nil && data["url"] != "" {
		fmt.Println("---------------------------==================---------------------------------")
		url := data["url"].(string)
		fmt.Println(url)
		go saveFile(url)
	}
	return nil
}

func saveFile(url string) error {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("saveFile: ", url)
	resp, err := http.Get(url)
	fmt.Println("Response:", resp)
	if err != nil {
		fmt.Println("Error:", err)
		log.Fatal(err)
		return err
	}

	fmt.Println("Resp:", resp)
	defer resp.Body.Close()

	//open a file for writing
	file, err := os.Create("/Users/prasad/images/" + uniuri.New() + ".jpg")
	fmt.Println("Saved file:", file)
	defer file.Close()
	if err != nil {
		fmt.Println("Error:", err)
		log.Fatal(err)
		return err
	}
	fmt.Println("New File:", file)
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		log.Fatal(err)
		return err
	}
	fmt.Println("Saved file:", file)
	return nil
}
