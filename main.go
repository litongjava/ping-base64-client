package main

import (
  "bytes"
  "flag"
  "fmt"
  "io"
  "io/ioutil"
  "log"
  "mime/multipart"
  "net/http"
  "os"
  "path/filepath"
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
  if len(*url) == 0 {
    fmt.Println("please specified remote server url")
    return
  }
  if len(*filePath) == 0 {
    fmt.Println("please specified local file path")
    return
  }
  log.Println("uploading...")
  start := time.Now().Unix()
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

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
  end := time.Now().Unix()
  fmt.Println(end-start, "s")
  fmt.Println("done")
}
