package gentoomen

import (
	"io/ioutil"
	md "github.com/knieriem/markdown"
	"bytes"
	"os"
    "fmt"
    "strings"
)

func GetPage(val string) string {
	if val == "" {
		val = "index"
    }

    val = strings.Trim(val, "/")
    html, err := getFile(val)
    if err != nil {
        return "<html><body><div><span style=\"color:red;font-weight:bold;\">" + err.String() + "</span></div></body></html>"
    }

    return html
}

func getFile(filename string) (string, os.Error) {
	file, err := os.Open("pages/" + filename + ".md")
	if err != nil {
		return "", err
	}

    cached, err := os.Open("cache/" + filename + ".html")
    if err == nil {
        fileStat, err := file.Stat()
        if err != nil {
            return "", err
        }
        fileTime := fileStat.Mtime_ns

        cachedStat, err := cached.Stat()
        if err != nil {
            return "", err
        }
        cachedTime := cachedStat.Mtime_ns

        file.Close()
        cached.Close()

        if fileTime < cachedTime {
            b, err := ioutil.ReadFile("cache/" + filename + ".html")

            if err != nil {
                return "", err
            }

            fmt.Printf("Loaded cached file ``%s''\n", filename + ".html")

            return string(b), nil
        }
    }

    file.Close()
    cached.Close()

	b, err := ioutil.ReadFile("pages/" + filename + ".md")
	if err != nil {
        return "", err
	}
	doc := md.Parse(string(b), md.Extensions{Smart: true})

    var buf bytes.Buffer
	w := bytes.NewBuffer(buf.Bytes())
	doc.WriteHtml(w)

    html, err := Template(w.String(), "page.html")
    if err != nil {
        return "", err
    }

    ioutil.WriteFile("cache/" + filename + ".html", []byte(html), 0644)

    fmt.Printf("Parsed the file ``%s''\n", filename + ".md")

	return html, nil
}

func Template(content, template string) (string, os.Error) {
    templBytes, err := ioutil.ReadFile("templates/" + template)
    if err != nil {
        return "", err
    }
    html := fmt.Sprintf(string(templBytes), content)

    return html, nil
}
