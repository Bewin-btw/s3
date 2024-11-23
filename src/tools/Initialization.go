package tools

import (
	"log"
	"os"

	"triple-s/src/vars"
)

func Init() {
	dirPath := vars.DirFlag

	_, err := os.Open(*dirPath)
	if err != nil {
		err := os.Mkdir(*dirPath, 0o700)
		if err != nil {
			err := os.MkdirAll(*dirPath, 0o755)
			if err != nil {
				log.Fatal("problem with creating directory: ", err)
			}
		}
		f, err := os.Create(*dirPath + "/buckets.csv")
		if err != nil {
			os.Mkdir("s3", 0o700)
			os.Mkdir(*dirPath, 0o700)
			f, err = os.Create(*dirPath + "/buckets.csv")
			if err != nil {
				log.Fatal(err)
			}
		}
		_, err = f.WriteString("Name,CreationTime,LastModifiedTime,Status\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}
