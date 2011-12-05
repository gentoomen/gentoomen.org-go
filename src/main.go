package main

import (
	"web"
    "./gentoomen"
    "io/ioutil"
)

func getStyle(context *web.Context) {
    b, err := ioutil.ReadFile("style.css")
    context.SetHeader("Content-Type", "text/css", true)
    if err != nil {
        context.WriteString(err.String())
    }

    context.Write(b)
}

func main() {
    web.Get("/style.css", getStyle)
    web.Get("/(.*)", gentoomen.GetPage)
    //web.Run("0.0.0.0:8080")
    web.RunFcgi("0.0.0.0:6580")
}
