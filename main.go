package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/replace", replaceURLContent)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	log.Println("收到根路径请求")
	fmt.Fprint(w, "hello world")
}

func fetchData(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var data []byte
	waitTime := 1.0

	for try := 1; try <= 3; try++ {
		time.Sleep(time.Duration(waitTime) * time.Second)

		resp, err := client.Get(url)
		if err != nil {
			waitTime *= 2
			continue
		}

		data, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			waitTime *= 2
			continue
		}

		return data, nil
	}

	return nil, errors.New("retry limit exceeded")
}

func replaceURLContent(w http.ResponseWriter, r *http.Request) {
	log.Println("收到替换路径请求")

	originalURL := r.URL.Query().Get("url")
	oldValue := r.URL.Query().Get("old_value")
	newValue := r.URL.Query().Get("new_value")

	if originalURL == "" {
		http.Error(w, "缺少参数 'url'", http.StatusBadRequest)
		return
	}

	if oldValue == "" {
		oldValue = "icook.hk"
	}

	if newValue == "" {
		newValue = "cfip.gay"
	}

	log.Printf("原始 URL: %s\n", originalURL)
	log.Printf("旧值: %s\n", oldValue)
	log.Printf("新值: %s\n", newValue)

	rawData, err := fetchData(originalURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("获取错误: %s", err), http.StatusBadRequest)
		return
	}

	decodedData, err := base64.StdEncoding.DecodeString(string(rawData))
	if err != nil {
		http.Error(w, fmt.Sprintf("解码错误: %s", err), http.StatusBadRequest)
		return
	}

	decodedString := string(decodedData)
	replacedData := strings.ReplaceAll(decodedString, oldValue, newValue)

	encodedReplacedData := base64.StdEncoding.EncodeToString([]byte(replacedData))
	w.Header().Set("Content-Type", "text/plain")

	log.Printf("替换后的数据: %s\n", replacedData)

	fmt.Fprint(w, encodedReplacedData)
}
