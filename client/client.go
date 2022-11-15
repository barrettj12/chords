package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/barrettj12/chords/backend/dblayer"
)

// Client makes it easy to access API methods
type Client struct {
	serverURL *url.URL
	authKey   string
}

const (
	API_ARTISTS = "/api/v0/artists"
	API_SONGS   = "/api/v0/songs"
	API_CHORDS  = "/api/v0/chords"
)

func NewClient(serverURL, authKey string) (*Client, error) {
	parsed, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	return &Client{parsed, authKey}, nil
}

func (c *Client) GetArtists() ([]string, error) {
	resp, err := c.request(requestParams{
		method: http.MethodGet,
		path:   API_ARTISTS,
	})
	if err != nil {
		return nil, err
	}

	artists := []string{}
	err = json.Unmarshal(resp, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (c *Client) GetSongs(artist, id, query *string) ([]dblayer.SongMeta, error) {
	resp, err := c.request(requestParams{
		method: http.MethodGet,
		path:   API_SONGS,
		queryParams: map[string]*string{
			"artist": artist,
			"id":     id,
			"query":  query,
		},
	})
	if err != nil {
		return nil, err
	}

	songs := []dblayer.SongMeta{}
	err = json.Unmarshal(resp, &songs)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (c *Client) NewSong(song dblayer.SongMeta) (dblayer.SongMeta, error) {
	data, err := json.Marshal(song)
	if err != nil {
		return dblayer.SongMeta{}, err
	}

	resp, err := c.request(requestParams{
		method:      http.MethodPost,
		path:        API_SONGS,
		auth:        true,
		body:        data,
		contentType: "application/json",
	})
	if err != nil {
		return dblayer.SongMeta{}, err
	}

	respSong := dblayer.SongMeta{}
	err = json.Unmarshal(resp, &respSong)
	if err != nil {
		return dblayer.SongMeta{}, err
	}

	return respSong, nil
}

func (c *Client) UpdateSong(id string, song dblayer.SongMeta) (dblayer.SongMeta, error) {
	data, err := json.Marshal(song)
	if err != nil {
		return dblayer.SongMeta{}, err
	}

	resp, err := c.request(requestParams{
		method: http.MethodPut,
		path:   API_SONGS,
		queryParams: map[string]*string{
			"id": &id,
		},
		auth:        true,
		body:        data,
		contentType: "application/json",
	})
	if err != nil {
		return dblayer.SongMeta{}, err
	}

	respSong := dblayer.SongMeta{}
	err = json.Unmarshal(resp, &respSong)
	if err != nil {
		return dblayer.SongMeta{}, err
	}

	return respSong, nil
}

func (c *Client) DeleteSong(id string) error {
	_, err := c.request(requestParams{
		method: http.MethodDelete,
		path:   API_SONGS,
		queryParams: map[string]*string{
			"id": &id,
		},
		auth: true,
	})
	return err
}

func (c *Client) GetChords(id string) ([]byte, error) {
	return c.request(requestParams{
		method: http.MethodGet,
		path:   API_CHORDS,
		queryParams: map[string]*string{
			"id": &id,
		},
	})
}

func (c *Client) UpdateChords(id string, chords []byte) ([]byte, error) {
	return c.request(requestParams{
		method: http.MethodPut,
		path:   API_CHORDS,
		queryParams: map[string]*string{
			"id": &id,
		},
		auth:        true,
		body:        chords,
		contentType: "text/plain",
	})
}

// HELPER METHODS

// Common logic for making HTTP requests
func (c *Client) request(rp requestParams) ([]byte, error) {
	// Prepare request URL
	endpoint := *c.serverURL
	endpoint.Path = rp.path

	// Add query params
	v := url.Values{}
	for key, val := range rp.queryParams {
		if val != nil {
			v.Add(key, *val)
		}
	}
	endpoint.RawQuery = v.Encode()

	// Prepare request
	req, err := http.NewRequest(rp.method, endpoint.String(), bytes.NewReader(rp.body))
	if err != nil {
		return nil, err
	}

	// Add headers
	if rp.contentType != "" {
		req.Header.Set("Content-Type", rp.contentType)
	}
	if rp.auth {
		req.Header.Set("Authorization", c.authKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// For 4xx/5xx response codes, we want to error
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("response has status %q", resp.Status)
	}

	// Read body and return
	return io.ReadAll(resp.Body)
}

// Parameters for an API request
type requestParams struct {
	method      string
	path        string
	queryParams map[string]*string
	auth        bool
	body        []byte
	contentType string
}
