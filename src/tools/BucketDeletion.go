package tools

import (
	"bufio"
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strings"

	"triple-s/src/vars"
)

func DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	dirPath := vars.DirFlag
	bucketName := r.PathValue("bucketname")

	if !ExistanceOfBucket(bucketName) {
		vars.PrintXMLError(w, http.StatusNotFound, "No such bucket")
		return
	}

	file, err := os.Open(*dirPath + "/" + bucketName + "/objects.csv")
	if err != nil {
		log.Print("couldn't open objects.csv")
		vars.PrintXMLError(w, http.StatusInternalServerError, "Objects.csv wasn't opened")
		return
	}
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Print("couldn't read objects.csv")
		vars.PrintXMLError(w, http.StatusInternalServerError, "Objects.csv wasn't readed")
		return
	}
	for i, line := range lines {
		if i != 0 && len(line) > 0 {
			log.Print("bucket deletion incomplete because of objects")
			vars.PrintXMLError(w, http.StatusConflict, "Objects inside the bucket")
			return
		}
	}

	err = os.Remove(*dirPath + "/" + bucketName + "/objects.csv")
	if err != nil {
		log.Print("bucket deletion incomplete because of os.Remove: ", err)
		vars.PrintXMLError(w, http.StatusInternalServerError, "Deletion failed")
		return
	}
	err = os.Remove(*dirPath + "/" + bucketName)
	if err != nil {
		log.Print("bucket deletion incomplete because of os.Remove: ", err)
		vars.PrintXMLError(w, http.StatusInternalServerError, "Deletion failed")
		return
	}
	updateBucketMetadata(bucketName)

	w.WriteHeader(http.StatusNoContent)
}

func updateBucketMetadata(bucketName string) error {
	var res [][]string
	f, err := os.Open(*vars.DirFlag + "/buckets.csv")
	if err != nil {
		return err
	}
	defer f.Close()
	reader := csv.NewReader(f)

	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for i, line := range lines {
		if len(line) > 0 && i != 0 && line[0] == bucketName {
			continue
		}
		res = append(res, line)
	}

	file, err := os.Create(*vars.DirFlag + "/buckets.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range res {
		_, err = writer.WriteString(strings.Join(line, ",") + "\n")
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
	}
	return nil
}
