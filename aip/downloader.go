package aip

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type Chart struct {
	Name string
	Icao string
	Url  string
}

const baseurl = "https://aip.dfs.de/BasicVFR/pages/"

func DownloadCharts(airports []AirportConfig) {
  failed := []AirportConfig{}
  for _, airport := range airports {
    url := getAIPUrl(airport.M)
    fmt.Printf("Processing: %v\n", airport)
    fmt.Printf("AIP URL: %s%s\n", baseurl, url)
    charts := getChartOverview(baseurl+url, airport.Icao)
    if len(charts) == 0 {
      failed = append(failed, airport)
    }
    for _, chart := range charts {
      b64 := downloadChart(chart)
      saveAsPdf(chart, b64)
    }
  }
  fmt.Printf("Failed to download: %v\n", failed)
}

func downloadChart(chart Chart) string {
	r, err := http.Get(chart.Url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	node, _ := html.Parse(r.Body)
	main := findNode(node, "main")
  img := findNode(main, "img")
	src := getAttr(img, "src")
  b64 := strings.Replace(src, "data:image/png;base64,", "", 1)
  return b64
}

func getChartOverview(curl string, icao string) []Chart {
	var result []Chart
	r, err := http.Get(curl)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	htmlNode, _ := html.Parse(r.Body)
	ul := findNode(htmlNode, "ul")
  if ul == nil {
    return result
  }
	var list []string
	for child := ul.FirstChild; child != nil; child = child.NextSibling {
		a := findNode(child, "a")
		if a != nil {
			link := getAttr(a, "href")
			list = append(list, link)
		}
	}

	length := len(list)
	for i, item := range list {

		r, _ := regexp.Compile("[a-zA-Z0-9]*\\.html$")
    bUrl := r.ReplaceAllString(curl, "")
		if i == length-2 {
			result = append(result, Chart{fmt.Sprintf("%s_airport_GND.pdf", icao), icao, bUrl + item})
		} else if i < length-2 {
			name := fmt.Sprintf("%s_airport_VAD.pdf", icao)
			if length > 3 {
				name = fmt.Sprintf("%s_airport_VAD %d.pdf", icao, i+1)
			}
			result = append(result, Chart{name, icao, bUrl + item})
		}
	}
	return result
}

func getAIPUrl(code string) string {
	var url = baseurl + code + ".html"
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	htmlNode, _ := html.Parse(r.Body)
	var aipUrl string
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "meta" {
			for _, a := range node.Attr {
				if a.Key == "content" && strings.Contains(a.Val, "url=") {
					r, _ := regexp.Compile("url=(.*)")
					m := r.FindStringSubmatch(a.Val)
					if len(m) > 0 {
						aipUrl = m[1]
					}
					return
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(htmlNode)
	return aipUrl
}

func findNode(node *html.Node, nodeType string) *html.Node {
	var foundNode *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == nodeType {
			foundNode = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(node)
	return foundNode
}

func getAttr(node *html.Node, attr string) string {
	for _, a := range node.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}
