package provider

import (
	"encoding/json"
	"github.com/patricekahlhorn/xdcc-search/search"
	"io"
	"net/http"
	"strconv"
)

type XdccRocksProvider struct{}

type XdccRocksResponse struct {
	Results []struct {
		Serverhost string `json:"serverhost"`
		Channels   []struct {
			Channelname string `json:"channelname"`
			Bots        []struct {
				Botname string `json:"botname"`
				Files   []struct {
					Packnumber   int `json:"packnumber"`
					Numdownloads int `json:"numdownloads"`
					File         struct {
						Filename        string  `json:"filename"`
						Filesizeinbytes float64 `json:"filesizeinbytes"`
					} `json:"file"`
				} `json:"files"`
			} `json:"bots"`
		} `json:"channels"`
	} `json:"results"`
	Maxpages int `json:"maxpages"`
}

var XdccRocksUrl = "https://xdcc.rocks/search/?getpages=true&searchword="

func (p *XdccRocksProvider) Search(query string) (search.Results, error) {

	data, err := p.fetchData(query)

	if err != nil {
		return search.Results{}, nil
	}

	if data.Maxpages > 1 {
		for i := 2; i < data.Maxpages+1; i++ {
			res, _ := p.fetchData(query + "&page=" + strconv.Itoa(i))
			data.Results = append(data.Results, res.Results...)
		}
	}

	searchResults := make(search.Results, 0, len(data.Results))

	for _, v := range data.Results {
		searchResults = append(searchResults, &search.Result{
			Network:  v.Serverhost,
			Bot:      v.Channels[0].Bots[0].Botname,
			Channel:  v.Channels[0].Channelname,
			Pack:     strconv.Itoa(v.Channels[0].Bots[0].Files[0].Packnumber),
			Gets:     v.Channels[0].Bots[0].Files[0].Numdownloads,
			FileSize: uint64(v.Channels[0].Bots[0].Files[0].File.Filesizeinbytes),
			FileName: v.Channels[0].Bots[0].Files[0].File.Filename,
			Provider: "xdcc.rocks",
		})
	}

	return searchResults, nil
}

func (p *XdccRocksProvider) fetchData(query string) (*XdccRocksResponse, error) {
	res, err := http.Get(XdccRocksUrl + query)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var r *XdccRocksResponse

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return r, err

}
