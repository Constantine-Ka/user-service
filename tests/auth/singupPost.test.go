package main


import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "http://localhost:8000/auth/sing-up"
	method := "POST"

	payload := strings.NewReader(`{`+"
	"+`
	"login":"admin",`+"
"+`
	"firstName": "Константин",`+"
"+`
	"email":"kanivec3@gmail.com",`+"
"+`
	"password":"admin",`+"
"+`
	"secondName":"Vit",`+"
"+`
	"lastName":"Ka",`+"
"+`
	"imagePath":"https://pkg.go.dev/static/shared/icon/content_copy_gm_grey_24dp.svg",`+"
"+`
	"gender":0,`+"
"+`
	"birthday":788281200,`+"
"+`
	"description":"Обо мне очень хорошо отзываются"`+"
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
