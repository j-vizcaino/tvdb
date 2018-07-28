package tvdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"net/url"
)

var baseURL = url.URL{
	Scheme: "https",
	Host:   "api.thetvdb.com",
}

type QueryOption func(*http.Request)

func WithLanguage(language string) QueryOption {
	return func (req *http.Request) {
		if len(language) > 0 {
			req.Header.Set("Accept-Language", language)
		}
	}
}

type Client interface {
	Token() string
	Options() ClientOptions
	SearchSeriesByName(string, ...QueryOption) ([]SeriesSearchResult, error)
	SeriesByID(int, ...QueryOption) (*Series, error)
}

type ClientOptions struct {
	APIKey   string
	UserKey  string
	Username string
}

type client struct {
	token       string
	tokenDate   time.Time
	httpClient  *http.Client
	options     ClientOptions
}

func NewClient(options ClientOptions) (Client, error) {
	c := &client{
		httpClient: &http.Client{},
		options:    options,
	}
	err := c.login()
	if err != nil {
		return nil, fmt.Errorf("login failed, %s", err)
	}
	return c, nil
}

func (c *client) URL(path string) url.URL {
	ret := baseURL
	ret.Path = path
	return ret
}

func (c *client) Languages() ([]Language, error) {
	var data LanguageData
	url := c.URL("/languages")

	if err := c.get(url, &data); err != nil {
		return data.Data, err
	}
	if len(data.Error) > 0 {
		return data.Data, fmt.Errorf(data.Error)
	}
	return data.Data, nil
}

func (c *client) Token() string {
	return c.token
}

func (c *client) Options() ClientOptions {
	return c.options
}

func (c *client) SeriesByID(id int, options ...QueryOption) (*Series, error) {
	var data SeriesData
	url := c.URL(fmt.Sprintf("/series/%d", id))
	if err := c.get(url, &data, options...); err != nil {
		return nil, err
	}
	return data.Data, nil
}

func (c *client) SearchSeriesByName(name string, options ...QueryOption) ([]SeriesSearchResult, error) {

	url := c.URL("/search/series")
	q := url.Query()
	q.Set("name", name)
	url.RawQuery = q.Encode()

	var result SeriesSearchResults
	err := c.get(url, &result, options...)
	return result.Data, err
}

func (c *client) login() error {
	loginData := map[string]string{
		"apiKey":   c.options.APIKey,
		"userKey":  c.options.UserKey,
		"username": c.options.Username,
	}
	postData, err := json.Marshal(loginData)
	if err != nil {
		return err
	}
	loginURL := c.URL("/login")
	res, err := c.httpClient.Post(loginURL.String(), "application/json", bytes.NewReader(postData))
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()

	var token Token
	err = decoder.Decode(&token); if err != nil {
		return err
	}
	if len(token.Error) > 0 {
		return fmt.Errorf("%s", token.Error)
	}
	c.token = token.Value
	return nil
}

func (c *client) get(url url.URL, out interface{}, options ...QueryOption) error {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return err
	}

	for _, optFunc := range options {
		optFunc(req)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	return decoder.Decode(out)
}
