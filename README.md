### README

Intention of the program is to download all the images and videos in the least time possible using go concurrency.
The specified 'root' directory and it's nested directory consists images, videos. Each image/video has an associated metadata JSON file. Incase the image/video is missing, only its associated JSON file is present. We should find all the JSON files, parse them and collect the urls. These urls would be used to retrieve the resource. We can add an check that the resourse doesnt not exist in the current file system