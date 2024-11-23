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

func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	dirPath := vars.DirFlag
	bucketName := r.PathValue("bucketname")
	objectName := r.PathValue("objectname")

	if !ExistanceOfBucket(bucketName) || !ExistanceOfObject(bucketName, objectName) {
		vars.PrintXMLError(w, http.StatusNotFound, "No such object in bucket/No such bucket")
		return
	}

	contentType, err := vars.GetObjectContentType(bucketName, objectName)
	if err != nil {
		log.Print("couldn't read the object metadata")
		vars.PrintXMLError(w, http.StatusNotFound, "No such object in metadata")
		return
	}

	err = os.Remove(*dirPath + "/" + bucketName + "/" + objectName + vars.GetFileExtension(contentType))
	if err != nil {
		log.Print("deletion unsuccessful")
		vars.PrintXMLError(w, http.StatusInternalServerError, "Object wasn't deleted")
	}

	if err := updateObjectMeta(bucketName, objectName); err != nil {
		log.Print("object metadata updating incompleted")
	}

	UpdateBucketsCSV(dirPath, bucketName, 2)
	w.WriteHeader(http.StatusNoContent)
}

func updateObjectMeta(bucketName string, objectName string) error {
	var res [][]string
	f, err := os.Open(*vars.DirFlag + "/" + bucketName + "/objects.csv")
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
		if len(line) > 0 && i != 0 && line[0] == objectName {
			continue
		}
		res = append(res, line)
	}

	file, err := os.Create(*vars.DirFlag + "/" + bucketName + "/" + "objects.csv")
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
