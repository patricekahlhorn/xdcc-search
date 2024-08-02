package provider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/patricekahlhorn/xdcc-search/search"
	"net/http"
	"regexp"
	"strings"
)

type XdccEuProvider struct{}

var xdccEuUrl = "https://www.xdcc.eu/search.php"

func (p *XdccEuProvider) Search(query string) (search.Results, error) {
	res, err := http.Get(xdccEuUrl + "?searchkey=" + query)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var searchResults search.Results

	doc.Find("tr").Each(func(k int, s *goquery.Selection) {
		// ignore header
		if k == 0 {
			return
		}
		children := s.Children()

		network, _ := children.Eq(1).Find("a").Attr("href")
		var re = regexp.MustCompile(`(?m).*?\/\/(.*?)\/.*`)
		network = re.ReplaceAllString(network, "$1")

		searchResults = append(searchResults, &search.Result{
			Network:  network,
			Bot:      children.Eq(2).Text(),
			Channel:  strings.Trim(children.Eq(1).Text(), " "),
			Pack:     children.Eq(3).Text(),
			Gets:     parseGets(children.Eq(4).Text()),
			FileSize: parseFileSize(children.Eq(5).Text()),
			FileName: strings.ReplaceAll(children.Eq(6).Text(), " ", "_"),
			Provider: "xdcc.eu",
		})
	})

	return searchResults, nil
}
