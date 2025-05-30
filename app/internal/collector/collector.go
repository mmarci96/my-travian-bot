package collector

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type ResourceValues struct {
	Wood     int
	Clay     int
	Iron     int
	Wheat    int
	FreeCrop int
}

// ParseResourceHTML parses the HTML and extracts the resource values.
func ParseResourceHTML(htmlStr string) (*ResourceValues, error) {
	fmt.Println(htmlStr)
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	res := &ResourceValues{}
	extractResources(doc, res)
	return res, nil
}

// extractResources recursively walks through HTML nodes and extracts values based on IDs.
func extractResources(n *html.Node, res *ResourceValues) {
	if n.Type == html.ElementNode && n.Data == "div" {
		var id string
		for _, attr := range n.Attr {
			if attr.Key == "id" {
				id = attr.Val
			}
		}

		switch id {
		case "l1", "l2", "l3", "l4", "stockBarFreeCrop":
			raw := getTextContent(n)
			cleaned := cleanUnicode(raw)
			value, err := strconv.Atoi(cleaned)
			if err == nil {
				switch id {
				case "l1":
					res.Wood = value
				case "l2":
					res.Clay = value
				case "l3":
					res.Iron = value
				case "l4":
					res.Wheat = value
				case "stockBarFreeCrop":
					res.FreeCrop = value
				}
			}
		}
	}

	// Recurse
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractResources(c, res)
	}
}

// getTextContent extracts all text within a node recursively.
func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(getTextContent(c))
	}
	return sb.String()
}

// cleanUnicode strips directional formatting and HTML numeric entities.
func cleanUnicode(s string) string {
	s = strings.ReplaceAll(s, "\u202d", "")
	s = strings.ReplaceAll(s, "\u202c", "")
	s = strings.ReplaceAll(s, "&#x202d;", "")
	s = strings.ReplaceAll(s, "&#x202c;", "")
	s = strings.ReplaceAll(s, "&amp;#x202d;", "")
	s = strings.ReplaceAll(s, "&amp;#x202c;", "")
	return strings.TrimSpace(s)
}
