package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
)

const kListenPort = 8099

func main() {
	// 接口回调测试
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

		log.Printf("/house/price method:%s,header:%v,body:%s \n", request.Method, request.Header, string(body))

		_, err = writer.Write(resData)
		if err != nil {
			log.Println(err.Error())
		}
		//_, _ = fmt.Fprintf(writer, string(resData))
	})

	// 服务接口方式测试
	// GET 输入：city，输出：price
	http.HandleFunc("/house/price2", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			log.Println("not support http post")
			return
		}

		city := request.URL.Query().Get("city")
		if city == "" {
			log.Println("invalid city")
			return
		}

		log.Println(fmt.Sprintf("/house/price2 city:%s", city))

		price := 52624
		switch city {
		case "北京":
			price = 57192
		case "深圳":
			price = 56569
		case "杭州":
			price = 28148
		}

		res := map[string]interface{}{
			"price": price,
		}
		data, err := json.Marshal(res)
		if err != nil {
			log.Println(err.Error())
			return
		}

		_, err = writer.Write(data)
		if err != nil {
			log.Println(err.Error())
		}
	})

	log.Print("server start on :", kListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", kListenPort), nil))
}
