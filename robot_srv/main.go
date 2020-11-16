package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
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
		root, err := simplejson.NewJson(body)
		if err != nil {
			log.Println(err.Error())
			return
		}

		nulValue := root.Get("nlu")
		answerValue, err := nulValue.Get("answer").String()
		if err != nil {
			log.Println(err.Error())
			return
		}

		nulValue.Set("answer", fmt.Sprintf("%s 再加工----------", answerValue))

		resData, err := nulValue.Encode()
		if err != nil {
			log.Println(err.Error())
			return
		}

		log.Printf("http get/post,method:%s,header:%v,body:%s \n", request.Method, request.Header, string(body))

		_, err = writer.Write(resData)
		if err != nil {
			log.Println(err.Error())
		}
		//_, _ = fmt.Fprintf(writer, string(resData))
	})
	log.Print("server start on :", kListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", kListenPort), nil))
}
