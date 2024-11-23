package tools

import (
	"encoding/csv"
	"encoding/xml"
	"log"
	"net/http"
	"os"

	"triple-s/src/vars"
)

type Response struct {
	Dir     string   `xml:"Dir"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

type Bucket struct {
	Name             string `xml:"Name"`
	CreationTime     string `xml:"CreationTime"`
	LastModifiedTime string `xml:"LastModifiedTime"`
	Status           string `xml:"Status"`
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	filePath := *vars.DirFlag + "/buckets.csv"
	f, err := os.Open(filePath)
	if err != nil {
		log.Print("no file with path: " + filePath)
		vars.PrintXMLError(w, http.StatusInternalServerError, "No file with such path")
		return
	} else {
		buckets := Response{*vars.DirFlag, readBuckets(f)}
		x, err := xml.Marshal(buckets)
		if err != nil {
			log.Print(err)
		} else {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			w.Write(x)
		}
	}
}

func readBuckets(f *os.File) []Bucket {
	var res []Bucket

	reader := csv.NewReader(f)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
		return res
	}

	for i, line := range lines {
		if len(line) >= 4 && i != 0 {
			res = append(res, Bucket{
				Name:             line[0],
				CreationTime:     line[1],
				LastModifiedTime: line[2],
				Status:           line[3],
			})
		}
	}
	return res
}
