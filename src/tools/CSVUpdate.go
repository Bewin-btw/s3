package tools

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"triple-s/src/vars"
)

// добавить тип операции(обновл/перепись)
func UpdateBucketsCSV(dirPath *string, bucketName string, typee int) {
	f, err := os.OpenFile(*dirPath+"/buckets.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fileInfo, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.Size() == 0 {
		_, err = f.WriteString("Name,CreationTime,LastModifiedTime,Status\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	// update whole bucket meta
	if typee == 1 {
		_, err = f.WriteString(bucketName + "," + time.Now().Format("2006-01-02 15:04:05") + "," + time.Now().Format("2006-01-02 15:04:05") + "," + "active\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	// update time only
	if typee == 2 {
		f.Close()
		f, err = os.Open(*dirPath + "/buckets.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		var lines []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Split(line, ",")
			if len(fields) >= 4 && fields[0] == bucketName {
				fields[2] = time.Now().Format("2006-01-02 15:04:05")
				line = strings.Join(fields, ",")
			}
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(*dirPath + "/buckets.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for _, line := range lines {
			_, err = writer.WriteString(line + "\n")
			if err != nil {
				log.Fatal(err)
			}
			writer.Flush()
		}
	}
}

func UpdateObjectCSV(dirPath *string, bucketName string, objectName string, contentType string, contentLengthHeader string, typee int) {
	f, err := os.OpenFile(*dirPath+"/"+bucketName+"/"+"objects.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fileInfo, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.Size() == 0 {
		_, err = f.WriteString("ObjectKey,Size,ContentType,LastModified\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	if typee == 1 {
		_, err = f.WriteString(objectName + "," + contentLengthHeader + "," + contentType + "," + time.Now().Format("2006-01-02 15:04:05") + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	if typee == 2 {
		f.Close()
		f, err = os.Open(*dirPath + "/" + bucketName + "/" + "objects.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		var lines []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Split(line, ",")
			if len(fields) >= 4 && fields[0] == objectName {
				fields[1] = contentLengthHeader
				if fields[2] != contentType {
					deletePreviousObjects(fields[2], bucketName, objectName)
					fields[2] = contentType
				}
				fields[3] = time.Now().Format("2006-01-02 15:04:05")
				line = strings.Join(fields, ",")
			}

			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(*dirPath + "/" + bucketName + "/" + "objects.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for _, line := range lines {
			_, err = writer.WriteString(line + "\n")
			if err != nil {
				log.Fatal(err)
			}
			writer.Flush()
		}
	}
	UpdateBucketsCSV(dirPath, bucketName, 2)
}

func deletePreviousObjects(contentType string, bucketName string, objectName string) {
	onDeletion := objectName + vars.GetFileExtension(contentType)

	err := os.Remove(*vars.DirFlag + "/" + bucketName + "/" + onDeletion)
	if err != nil {
		log.Print("error on deletion of an object ", err)
	}
}
