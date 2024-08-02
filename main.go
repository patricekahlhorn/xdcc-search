package main

import (
	"fmt"
	"github.com/patricekahlhorn/xdcc-search/client"
	"github.com/patricekahlhorn/xdcc-search/search/provider"
	"time"
)

func main() {
	c := client.NewClient(
		client.WithTimeout(5000*time.Millisecond),
		client.WithProviders([]provider.Provider{
			&provider.SunXdccProvider{},
			&provider.XdccRocksProvider{},
			&provider.XdccEuProvider{},
		}),
	)

	keywords := []string{
		"ubuntu-20.04-desktop-amd64.iso",
	}

	res, err := c.Search(keywords)

	if res != nil {
		res = res.RemoveDuplicates().Sort()
	}

	for _, r := range res {
		fmt.Printf("%+v\n", r)
	}

	if err != nil {
		fmt.Println(err)
	}
}
