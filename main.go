package main

import (
  "bytes"
  "crypto/tls"
  "encoding/base64"
  "flag"
  "fmt"
  "github.com/spf13/viper"
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

func init() {
  log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}
func main() {
  a := flag.String("a", "", "action")
  e := flag.String("e", "dev", "environment:such dev,test,prod")
  url := flag.String("url", "", "remote server url")
  //filePath := target/malang-pen-api-server-1.0.0.jar
  filePath := flag.String("file", "", "local file path")
  //d := "/data/apps/web/ztbjg"
  m := flag.String("m", "", "move to file path")
  d := flag.String("d", "", "extra file path")
  c := flag.String("c", "", "full command")

  flag.Parse()

  // read config file
  viper.SetConfigFile("ping-base64.toml")
  viper.SetConfigType("toml")
  if err := viper.ReadInConfig(); err != nil {
    log.Println("Error reading config file, %s", err)
  }

  // read config item
  if *url == "" {
    *url = viper.GetString(*e + "." + *a + ".url")
  }

  if *a == "upload-run" {

    if *filePath == "" {
      *filePath = viper.GetString(*e + ".upload-run.file")
    }

    if *m == "" {
      *m = viper.GetString(*e + ".upload-run.m")
    }

    if *d == "" {
      *d = viper.GetString(*e + ".upload-run.d")
    }

    if *c == "" {
      *c = viper.GetString(*e + "." + *a + ".c")
    }

    log.Println("url:", *url)
    log.Println("filePath:", *filePath)
    log.Println("m:", *m)
    log.Println("d:", *d)
    log.Println("c:", *c)
  }

  if *a == "web" {
    if *c == "" {
      *c = viper.GetString(*e + "." + *a + ".c")
    }
    log.Println("url:", *url)
    log.Println("c:", *c)
  }

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
    runRemoteCmd(url, c)
  }
  end := time.Now().Unix()
  fmt.Println(end-start, "s")
  fmt.Println("done")
}

func uploadAndRun(url *string, filePath *string, m *string, d *string, c *string) {
  var file, errFile1 = os.Open(*filePath)
  if errFile1 != nil {
    log.Fatalln(errFile1)
  }
  defer file.Close()
  var fileInfo, err = file.Stat()
  if err != nil {
    log.Fatalln(err)
  }
  var fileSize = fileInfo.Size()

  log.Printf("file is uploading,and file size is %d KB", fileSize/1024)

  method := "POST"
  payload := &bytes.Buffer{}

  // io.Writer
  writer := multipart.NewWriter(payload)
  part1, errFile1 := writer.CreateFormFile("file", filepath.Base(*filePath))
  _, errFile1 = io.Copy(part1, file)
  if errFile1 != nil {
    log.Fatalln(errFile1)
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

  err = writer.Close()
  if err != nil {
    fmt.Println(err)
    return
  }

  req, err := http.NewRequest(method, *url, payload)
  if err != nil {
    log.Fatalln(err)
  }
  req.Header.Set("Content-Type", writer.FormDataContentType())

  client := &http.Client{
    Transport: &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    },
  }
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
    log.Fatalln(err)
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
