package main

import (
	"web"
	"io/ioutil"
	md "./markdown"
	"bytes"
	"fmt"
	"strings"
)

func hello(val string) string {
	filename := strings.SplitN(val, "/", 2)
	b,_ := ioutil.ReadFile("./"+filename[1])
	fmt.Println(filename)
	fmt.Println(val)
	doc := md.Parse(string(b), md.Extensions{Smart: true})

	var buf bytes.Buffer

	w := bytes.NewBuffer(buf.Bytes())
	doc.WriteHtml(w)
	return w.String()
}

func main() {
    web.Get("/(.*)", hello)
    web.RunFcgi("0.0.0.0:6850")
}
