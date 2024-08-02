package search

import (
	"fmt"
	"sort"
)

type Result struct {
	BotRecord string
	Network   string
	Bot       string
	Channel   string
	Pack      string
	Gets      int
	FileSize  uint64
	FileName  string
	Provider  string
}

type Results []*Result

func (r Results) RemoveDuplicates() Results {
	allKeys := make(map[string]bool)
	var list Results
	for _, result := range r {
		if _, value := allKeys[result.ToUrl()]; !value {
			allKeys[result.ToUrl()] = true
			list = append(list, result)
		}
	}
	return list
}

func (r Results) Sort() Results {
	sort.Slice(r, func(i, j int) bool {
		return r[i].FileName < r[j].FileName
	})

	return r
}

func (r Result) ToUrl() string {
	return fmt.Sprintf("irc://%s/%s/%s/%s", r.Network, r.Channel, r.Bot, r.Pack)
}
