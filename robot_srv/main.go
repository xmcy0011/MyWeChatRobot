package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const kListenPort = 8099

func main() {
	http.HandleFunc("/house/price", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "POST" {
			_, _ = fmt.Fprintf(writer, "not support http get")
			return
		}
		if request.ContentLength > 10*1024*1024 {
			log.Printf("content length to large,> 10MB")
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Println(err.Error())
			return
		}

		log.Printf("http get/post,method:%s,header:%v,body:%s \n", request.Method, request.Header, string(body))
		_, _ = fmt.Fprintf(writer, "success")
	})
	log.Print("server start on :", kListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", kListenPort), nil))
}
