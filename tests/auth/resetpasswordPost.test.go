package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "http://localhost:8000/auth/resetpassword"
	method := "POST"

	payload := strings.NewReader(`{`+"
	"+`
	"confirm":"6b6b6b6b6b6b6b6b6b6b6b6b6b39e4592986fb69508e88d3a0517efbdb7446034d",`+"
"+`
	"password":"admin1",`+"
"+`
	"password2":"admin1"`+"
"+`
}`)

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, payload)

  if err != nil {
    fmt.Println(err)
    return
  }
  req.Header.Add("Content-Type", "application/json")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
}