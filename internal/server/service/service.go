package service

import (
	"strings"

	"golang.org/x/net/html"
)

type Service struct {
}

func (s *Service) ParseString(doc *html.Node) (url map[string]*html.Node, err error) {
	var processAllProduct func(*html.Node)
	processAllProduct = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "html" {
			url = processNode(n)

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processAllProduct(c)
		}
	}
	processAllProduct(doc)

	return url, nil
}

var urls = make(map[string]*html.Node, 0)

func processNode(n *html.Node) map[string]*html.Node {

	switch n.Data {

	case "a":
		for _, a := range n.Attr {
			if a.Key == "href" {
				if strings.Contains(a.Val, "wiki") {
					if string(a.Val[0]) == "/" {
						urls[a.Val] = n
						urls["https://ru.wikipedia.org"+a.Val] = n
					} else {
						urls[a.Val] = n
					}
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processNode(c)
	}

	return urls
}
