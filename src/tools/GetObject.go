package tools

import (
	"io"
	"log"
	"net/http"
	"os"

	"triple-s/src/vars"
)

func GetObject(w http.ResponseWriter, r *http.Request) {
	dirPath := vars.DirFlag
	bucketName := r.PathValue("bucketname")
	objectName := r.PathValue("objectname")
	contentType, err := vars.GetObjectContentType(bucketName, objectName)
	if err != nil {
		log.Print("couldn't read the object metadata")
		vars.PrintXMLError(w, http.StatusNotFound, "No such object in metadata")
		return
	}

	if !ExistanceOfBucket(bucketName) || !ExistanceOfObject(bucketName, objectName) {
		vars.PrintXMLError(w, http.StatusNotFound, "No such object in bucket/No such bucket")
		return
	}

	file, err := os.Open(*dirPath + "/" + bucketName + "/" + objectName + vars.GetFileExtension(contentType))
	if err != nil {
		vars.PrintXMLError(w, http.StatusNotFound, "File not found")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", contentType)

	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("couldn't copy file: %v", err)
		vars.PrintXMLError(w, http.StatusInternalServerError, "Failed to send file")
		return
	}

	w.WriteHeader(http.StatusOK)
}
