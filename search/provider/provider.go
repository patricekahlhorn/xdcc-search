package provider

import (
	"github.com/dustin/go-humanize"
	"github.com/patricekahlhorn/xdcc-search/search"
	"regexp"
	"strconv"
	"strings"
)

func parseGets(str string) int {
	gets, _ := strconv.Atoi(strings.ReplaceAll(str, "x", ""))

	return gets
}

func parseFileSize(str string) uint64 {
	var re = regexp.MustCompile(`(?m).*?(\d*[A-Z]).*`)
	str = re.ReplaceAllString(str, "$1")
	bytes, _ := humanize.ParseBytes(str)

	return bytes
}

type Provider interface {
	Search(query string) (search.Results, error)
}
