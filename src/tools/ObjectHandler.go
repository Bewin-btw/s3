package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"triple-s/src/vars"
)

func ObjectHandler(w http.ResponseWriter, r *http.Request) {
	objectName := r.PathValue("objectname")
	bucketName := r.PathValue("bucketname")
	contentType := r.Header.Get("Content-Type")
	ext := vars.GetFileExtension(contentType)
	contentLengthHeader := r.Header.Get("Content-Length")

	fmt.Printf("Content type header: : %s\n", contentType)
	fmt.Printf("Content length header: %s\n", contentLengthHeader)

	if !ExistanceOfBucket(bucketName) {
		vars.PrintXMLError(w, http.StatusNotFound, "No such bucket")
		return
	}

	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		vars.PrintXMLError(w, http.StatusInternalServerError, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	if !contentCheck(bodyData, contentType, contentLengthHeader, w) {
		return
	}

	file, err := os.Create(*vars.DirFlag + "/" + bucketName + "/" + objectName + ext)
	if err != nil {
		vars.PrintXMLError(w, http.StatusInternalServerError, "Failed to create file")
		return
	}
	defer file.Close()

	_, err = file.Write(bodyData)
	if err != nil {
		vars.PrintXMLError(w, http.StatusInternalServerError, "Failed to write to file")
		return
	}

	if !ExistanceOfObject(bucketName, objectName) {
		UpdateObjectCSV(vars.DirFlag, bucketName, objectName, contentType, contentLengthHeader, 1)
	} else {
		UpdateObjectCSV(vars.DirFlag, bucketName, objectName, contentType, contentLengthHeader, 2)
	}

	w.WriteHeader(http.StatusOK)
}

func contentCheck(bodyData []byte, contentType string, contentLengthHeader string, w http.ResponseWriter) bool {
	if contentType == "" {
		vars.PrintXMLError(w, http.StatusBadRequest, "Content-Type header is required")
		return false
	}
	contentLength, err := strconv.Atoi(contentLengthHeader)
	if err != nil {
		vars.PrintXMLError(w, http.StatusBadRequest, "Invalid Content-Length")
		return false
	}

	if len(bodyData) != contentLength {
		vars.PrintXMLError(w, http.StatusBadRequest, "Content-Length does not match the body length")
		return false
	}
	return true
}
