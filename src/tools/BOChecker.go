package tools

import (
	"bufio"
	"log"
	"os"

	"triple-s/src/vars"
)

func ExistanceOfBucket(bucketName string) bool {
	bucketNameSize := len(bucketName)

	metaData, err := os.Open(*vars.DirFlag + "/buckets.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer metaData.Close()

	scanner := bufio.NewScanner(metaData)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > bucketNameSize && line[:bucketNameSize] == bucketName {
			return true
		}
	}

	return false
}

func ExistanceOfObject(bucketName string, objectName string) bool {
	objectNameSize := len(objectName)

	metaData, err := os.Open(*vars.DirFlag + "/" + bucketName + "/objects.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer metaData.Close()

	scanner := bufio.NewScanner(metaData)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > objectNameSize && line[:objectNameSize] == objectName {
			return true
		}
	}
	return false
}
