package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	//url := "https://zwfw-test.dgprj.com/ztbjg-portal/ping-base64-webapi/file/upload-unzip/"
	url := flag.String("url", "", "remote server url")
	//filePath := "E:\code\project\ShuZiGuangDong\ztjgw-personal\ztbjg-yhzx.zip"
	filePath := flag.String("file", "", "local file path")
	//d := "/data/apps/web/ztbjg"
	m := flag.String("m", "", "move to file path")
	d := flag.String("d", "", "extra file path")
	c := flag.String("c", "", "full command")
	flag.Parse()
	log.Println("url:", *url)
	log.Println("filePath:", *filePath)
	log.Println("m:", *m)
	log.Println("d:", *d)
	log.Println("c:", *c)
	start := time.Now().Unix()
	if strings.HasSuffix(*url, "upload-run/") {
		log.Println("upload-run")
		if len(*url) == 0 {
			fmt.Println("please specified remote server url")
			return
		}
		if len(*filePath) == 0 {
			fmt.Println("please specified local file path")
			return
		}
		uploadAndRun(url, filePath, m, d, c)
	} else if strings.HasSuffix(*url, "web/") {
		log.Println("web")
		flag.Parse()
		runRemoteCmd(url, c)
	}
	end := time.Now().Unix()
	fmt.Println(end-start, "s")
	fmt.Println("done")
}

func uploadAndRun(url *string, filePath *string, m *string, d *string, c *string) {
	log.Println("uploading...")

	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, errFile1 := os.Open(*filePath)
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("file", filepath.Base(*filePath))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}

	if len(*m) != 0 {
		_ = writer.WriteField("m", *m)
	}

	if len(*d) != 0 {
		_ = writer.WriteField("d", *d)
	}

	if len(*c) != 0 {
		_ = writer.WriteField("c", *c)
	}

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, *url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "apifox/1.0.0 (https://www.apifox.cn)")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	code := res.StatusCode
	fmt.Println("response status dode", code)
	defer res.Body.Close()
	//goland:noinspection ALL
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func runRemoteCmd(url *string, c *string) {
	log.Println("running...")
	cmd := base64.StdEncoding.EncodeToString([]byte(*c))
	newUrl := *url + cmd
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, newUrl, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	//goland:noinspection ALL
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	decodedBytes := make([]byte, base64.StdEncoding.DecodedLen(len(body)))
	_, err = base64.StdEncoding.Decode(decodedBytes, []byte(body))
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}

	fmt.Println(string(decodedBytes))
}
