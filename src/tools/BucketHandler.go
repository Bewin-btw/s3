package tools

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"triple-s/src/vars"
)

// templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// working progress without error handling
// лучше перенести создание директории и файлов в отделбную функцию, которая вызывается на инициализации - DONE
func BucketHandler(w http.ResponseWriter, r *http.Request) {
	if !PathValidation(r.PathValue("bucketname")) || vars.ValidPath.FindStringSubmatch(r.URL.Path) == nil {
		log.Print("invalid path")
		vars.PrintXMLError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	dirPath := vars.DirFlag
	bucketName := r.PathValue("bucketname")

	if ExistanceOfBucket(bucketName) {
		vars.PrintXMLError(w, http.StatusConflict, "Bucket name should be unique")
		return
	}
	// content ,err :=io.ReadAll(r.Body)
	if err := os.Mkdir(*dirPath+"/"+bucketName, 0o700); err != nil {
		log.Fatal(err)
	} else {
		f, err := os.Create(*dirPath + "/" + bucketName + "/objects.csv")
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.WriteString("ObjectKey,Size,ContentType,LastModified\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	UpdateBucketsCSV(dirPath, bucketName, 1)
	w.WriteHeader(http.StatusOK)

	xmlResponse := fmt.Sprintf(`
		<BucketInfo>
			<Name>%s</Name>
			<CreationTime>%s</CreationTime>
			<LastModifiedTime>%s</LastModifiedTime>
			<Status>%s</Status>
		</BucketInfo>`, bucketName, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), "active")

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(xmlResponse))
}
