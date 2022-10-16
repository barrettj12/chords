package client

import (
	"bytes"
	"encoding/json"
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
	endpoint := *c.serverURL
	endpoint.Path = API_ARTISTS
	resp, err := http.Get(endpoint.String())
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	artists := []string{}
	err = json.Unmarshal(data, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (c *Client) GetSongs(artist, id, query *string) ([]dblayer.SongMeta, error) {
	endpoint := *c.serverURL
	endpoint.Path = API_SONGS

	// Query params
	v := url.Values{}
	if artist != nil {
		v.Add("artist", *artist)
	}
	if id != nil {
		v.Add("id", *id)
	}
	if query != nil {
		v.Add("query", *query)
	}
	endpoint.RawQuery = v.Encode()

	resp, err := http.Get(endpoint.String())
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	songs := []dblayer.SongMeta{}
	err = json.Unmarshal(data, &songs)
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

	endpoint := *c.serverURL
	endpoint.Path = API_SONGS

	resp, err := http.Post(endpoint.String(), "application/json", bytes.NewReader(data))
	if err != nil {
		return dblayer.SongMeta{}, err
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return dblayer.SongMeta{}, err
	}

	respSong := dblayer.SongMeta{}
	err = json.Unmarshal(data, &respSong)
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

	endpoint := *c.serverURL
	endpoint.Path = API_SONGS
	v := url.Values{}
	v.Add("id", id)
	endpoint.RawQuery = v.Encode()

	resp, err := httpPut(endpoint.String(), "application/json", bytes.NewReader(data))
	if err != nil {
		return dblayer.SongMeta{}, err
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return dblayer.SongMeta{}, err
	}

	respSong := dblayer.SongMeta{}
	err = json.Unmarshal(data, &respSong)
	if err != nil {
		return dblayer.SongMeta{}, err
	}

	return respSong, nil
}

func (c *Client) DeleteSong(id string) error {
	endpoint := *c.serverURL
	endpoint.Path = API_SONGS
	v := url.Values{}
	v.Add("id", id)
	endpoint.RawQuery = v.Encode()

	_, err := httpDelete(endpoint.String())
	return err
}

func (c *Client) GetChords(id string) ([]byte, error) {
	endpoint := *c.serverURL
	endpoint.Path = API_CHORDS
	v := url.Values{}
	v.Add("id", id)
	endpoint.RawQuery = v.Encode()

	resp, err := http.Get(endpoint.String())
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) UpdateChords(id string, chords []byte) ([]byte, error) {
	endpoint := *c.serverURL
	endpoint.Path = API_CHORDS
	v := url.Values{}
	v.Add("id", id)
	endpoint.RawQuery = v.Encode()

	resp, err := httpPut(endpoint.String(), "text/plain", bytes.NewReader(chords))
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

// HELPER METHODS

func httpPut(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
}

func httpDelete(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}
