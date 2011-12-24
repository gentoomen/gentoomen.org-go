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
	//context.SetHeader("Content-Type", "application/xhtml+xml; charset=utf-8", true) // xhtml MUST use application/xhtml+xml

	if val == "" {
		val = "index"
	}

	val = strings.Trim(val, "/")
	html, err := getFile(val)
	if err != nil {
		return "<html><head><title>Error</title></head><body><div style=\"text-align:center;\"><span style=\"color:red;font-weight:bold;\">" + err.String() + "</span></div></body></html>"
	}

	return html
}

func getFile(filename string) (string, os.Error) {
	if checkCache(filename) {
		b, err := ioutil.ReadFile("cache/" + filename + ".html")

		if err != nil {
			return "", err
		}

		fmt.Printf("Loaded cached file ``%s''\n", filename + ".html")

		return string(b), nil
	}

	b, err := ioutil.ReadFile("pages/" + filename + ".md")
	if err != nil {
		return "", err
	}
	doc := md.Parse(string(b), md.Extensions{Smart: true})

	var buf bytes.Buffer
	doc.WriteHtml(&buf)

	html, err := Template(getLinks(), getProjects(), buf.String(), "page.html")
	if err != nil {
		return "", err
	}

	ioutil.WriteFile("cache/" + filename + ".html", []byte(html), 0644)

	fmt.Printf("Parsed the file ``%s''\n", filename + ".md")

	return html, nil
}

func checkCache(filename string) bool {
	file := "pages/" + filename + ".md"
	cached := "cache/" + filename + ".html"
	links := "links.txt"
	projects := "projects.txt"

	cacheTime := getModifiedTime(cached)

	return cacheTime > 0 && cacheTime > getModifiedTime(file) && cacheTime > getModifiedTime(links) && cacheTime > getModifiedTime(projects)
}

func getModifiedTime(filename string) int64 {
	file, err := os.Open(filename)
	if err != nil {
		return -1
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return -1
	}

	return fileStat.Mtime_ns
}

func Template(links, projects, content, template string) (string, os.Error) {
	templBytes, err := ioutil.ReadFile("templates/" + template)
	if err != nil {
		return "", err
	}
	html := fmt.Sprintf(string(templBytes), links, projects, content)

	return html, nil
}

func getLinks() string {
	template := "<li><a href=\"%s\">%s</a></li>\n"

	return getUrls("links.txt", template)
}

func getProjects() string {
	template := "<li><a href=\"%s\">%s</a></li>\n"

	return getUrls("projects.txt", template)
}

func getUrls(file, template string) string {
	linksb, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading ``%s'': %s\n", file, err.String())
		return fmt.Sprintf(template, "/", "Home")
	}

	var buf bytes.Buffer
	links := strings.Split(strings.TrimSpace(string(linksb)), "\n")
	for _, line := range links {
		arr := strings.SplitN(line, ":", 2)
		if len(arr) == 2 {
			buf.WriteString(fmt.Sprintf(template, strings.TrimSpace(arr[1]), strings.TrimSpace(arr[0])))
		} else {
			fmt.Fprintf(os.Stderr, "Error reading ``%s'', invalid line: %s\n", file, line)
		}
	}

	return buf.String()
}
