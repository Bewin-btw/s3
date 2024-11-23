package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"triple-s/src/tools"
	"triple-s/src/vars"
)

func main() {
	// http.HandleFunc("/view/", viewHandler)
	flag.Parse()

	if d, err := tools.IsWithinBaseDir(vars.BaseDir, *vars.DirFlag); !d {
		fmt.Println(err)
		log.Fatal("attempt to exit base dir")
	}

	if *vars.PortFlag < 1024 {
		log.Fatal("invalid port")
	}

	tools.Init()

	http.HandleFunc("PUT /{bucketname}", tools.BucketHandler)
	http.HandleFunc("DELETE /{bucketname}", tools.DeleteBucketHandler)
	http.HandleFunc("PUT /{bucketname}/{objectname}", tools.ObjectHandler)
	http.HandleFunc("DELETE /{bucketname}/{objectname}", tools.DeleteObjectHandler)
	http.HandleFunc("GET /{bucketname}/{objectname}", tools.GetObject)
	http.HandleFunc("/", tools.BadPathHandler)
	http.HandleFunc("GET /{$}", tools.GetAll)

	if *vars.HelpFlag {
		vars.HelpFunc()
		return
	}
	log.Println("Server started listening on port:", *vars.PortFlag)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*vars.PortFlag), nil))
}
