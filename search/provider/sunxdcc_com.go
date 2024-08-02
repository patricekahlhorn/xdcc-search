package provider

import (
	"encoding/json"
	"errors"
	"github.com/patricekahlhorn/xdcc-search/search"
	"io"
	"net/http"
	"strings"
)

type SunXdccProvider struct{}

type Response struct {
	Botrec  []string `json:"botrec"`
	Network []string `json:"network"`
	Bot     []string `json:"bot"`
	Channel []string `json:"channel"`
	Packnum []string `json:"packnum"`
	Gets    []string `json:"gets"`
	Fsize   []string `json:"fsize"`
	Fname   []string `json:"fname"`
}

var SunXdccUrl = "https://sunxdcc.com/deliver.php"

func (p *SunXdccProvider) Search(query string) (search.Results, error) {
	res, err := http.Get(SunXdccUrl + "?sterm=" + query)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var r *Response

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if err = r.validate(); err != nil {
		return nil, err
	}

	searchResults := make(search.Results, 0, len(r.Fname))

	for k := range r.Fname {
		searchResults = append(searchResults, &search.Result{
			BotRecord: r.Botrec[k],
			Network:   r.Network[k],
			Bot:       r.Bot[k],
			Channel:   r.Channel[k],
			Pack:      r.Packnum[k],
			Gets:      parseGets(r.Gets[k]),
			FileSize:  parseFileSize(r.Fsize[k]),
			FileName:  strings.ReplaceAll(r.Fname[k], " ", "_"),
			Provider:  "sunxdcc.com",
		})
	}

	return searchResults, nil
}

func (r *Response) validate() error {
	sizes := [8]int{
		len(r.Botrec),
		len(r.Network),
		len(r.Bot),
		len(r.Channel),
		len(r.Packnum),
		len(r.Gets),
		len(r.Fsize),
		len(r.Fname),
	}

	length := sizes[0]
	for _, l := range sizes {
		if length != l {
			return errors.New("invalid response")

		}
	}
	return nil
}
