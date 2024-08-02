package provider

import (
	"context"
	"github.com/patricekahlhorn/xdcc-search/search"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSunXdccProvider_Search(t *testing.T) {

	searchTerm := "searchTerm"

	validResponse := "{\"botrec\":[\"100000.0kB\\/s\",\"200000.0kB\\/s\"],\"network\":[\"some.irc.network.net\",\"another.irc.network.net\"],\"bot\":[\"Bot1\",\"Bot2\"],\"channel\":[\"Channel1\",\"Channel2\"],\"packnum\":[\"#123\",\"#456\"],\"gets\":[\"5x\",\"6x\"],\"fsize\":[\"[5.8G]\",\"[10M]\"],\"fname\":[\"File1\",\"File2\"]}"
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("sterm") != searchTerm {
					t.Errorf("Invalid search query', got: %s", r.URL.Query().Get("sterm"))
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(validResponse))
			},
		),
	)
	defer server.Close()

	SunXdccUrl = server.URL

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    search.Results
		wantErr bool
	}{
		{
			"Parse SunXDCC Response",
			args{context.Background()},
			search.Results{
				&search.Result{
					BotRecord: "100000.0kB/s",
					Network:   "some.irc.network.net",
					Bot:       "Bot1",
					Channel:   "Channel1",
					Pack:      "#123",
					Gets:      5,
					FileSize:  8000000000,
					FileName:  "File1",
					Provider:  "sunxdcc.com",
				},
				&search.Result{
					BotRecord: "200000.0kB/s",
					Network:   "another.irc.network.net",
					Bot:       "Bot2",
					Channel:   "Channel2",
					Pack:      "#456",
					Gets:      6,
					FileSize:  10000000,
					FileName:  "File2",
					Provider:  "sunxdcc.com",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := SunXdccProvider{}
			got, err := p.Search(searchTerm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for k, v := range got {
				if *v != *tt.want[k] {
					t.Errorf("Search() got = %v, want %v", *v, *tt.want[k])
				}
			}
		})
	}
}
