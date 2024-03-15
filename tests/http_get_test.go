package tests

import (
  "crypto/tls"
  "io/ioutil"
  "log"
  "net/http"
  "testing"
)

func TestInsecureHttps(t *testing.T) {
  // 创建一个http.Client实例
  client := &http.Client{
    Transport: &http.Transport{
      // 配置TLS来信任所有证书
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    },
  }

  // 使用自定义的client发起请求
  var url = ""
  resp, err := client.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()

  // 读取响应内容
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }

  log.Println(string(body))
}
