// Package tvdb provides a Go client implementation of The TVDB API v2
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

type requestOption func(*http.Request)

// QueryOption provides a query option function.
type QueryOption func(*url.Values)

func withLanguage(language string) requestOption {
	return func (req *http.Request) {
		if len(language) > 0 {
			req.Header.Set("Accept-Language", language)
		}
	}
}

// WithAiredSeasonNumber filters episodes matching the given season number.
//
// This function should be used with EpisodesBySeriesID.
func WithAiredSeasonNumber(seasonNumber int) QueryOption {
	return withQueryIntOption("airedSeason", seasonNumber)
}

// WithAiredEpisodeNumber filters episodes matching the given episode number.
//
// This function should be used with EpisodesBySeriesID.
func WithAiredEpisodeNumber(episodeNumber int) QueryOption {
	return withQueryIntOption("airedEpisode", episodeNumber)
}

// WithDVDSeasonNumber filters episodes matching the given DVD season number.
//
// This function should be used with EpisodesBySeriesID
func WithDVDSeasonNumber(seasonNumber int) QueryOption {
	return withQueryIntOption("dvdSeason", seasonNumber)
}

// WithDVDEpisodeNumber filters episodes matching the given DVD episode number.
//
// This function should be used with EpisodesBySeriesID
func WithDVDEpisodeNumber(episodeNumber int) QueryOption {
	return withQueryIntOption("dvdEpisode", episodeNumber)
}

// WithAbsoluteEpisodeNumber filters episodes matching the given absolute episode number.
//
// This function should be used with EpisodesBySeriesID
func WithAbsoluteEpisodeNumber(episodeNumber int) QueryOption {
	return withQueryIntOption("absoluteNumber", episodeNumber)
}

func withQueryIntOption(name string, value int) QueryOption {
	return func(v *url.Values) {
		v.Set(name, fmt.Sprintf("%d", value))
	}
}

func withQueryOption(name, value string) QueryOption {
	return func(v *url.Values) {
		v.Set(name, value)
	}
}

// ClientOptions represents options for a TVDB client.
//
// Either APIKey or UserKey and Username are mandatory for login.
// Language provides a hint to return result in the given language and can be changed dynamically
// using the WithLanguage() helper.
type ClientOptions struct {
	APIKey   string
	UserKey  string
	Username string
	Language string
}

type client struct {
	token       string
	tokenDate   time.Time
	httpClient  *http.Client
	options     ClientOptions
}

// NewClient creates a new TVDB client.
//
// The function immediately tries to login and returns the handle of the new client.
func NewClient(options ClientOptions) (*client, error) {
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

// URL builds a request URL for the given path and options.
func (c *client) URL(path string, options ...QueryOption) url.URL {
	ret := baseURL
	ret.Path = path

	if len(options) > 0 {
		q := ret.Query()
		for _, opt := range options {
			opt(&q)
		}
		ret.RawQuery = q.Encode()
	}
	return ret
}

// Languages returns a list of languages supported by The TVDB.
func (c *client) Languages() ([]Language, error) {
	var data LanguageData
	fullURL := c.URL("/languages")

	if err := c.get(fullURL, &data); err != nil {
		return data.Data, err
	}
	if len(data.Error) > 0 {
		return data.Data, fmt.Errorf(data.Error)
	}
	return data.Data, nil
}

// Token returns the JWT used for authentication.
func (c *client) Token() string {
	return c.token
}

// Options returns a copy of the options used by the client.
func (c *client) Options() ClientOptions {
	return c.options
}

// WithLanguage updates the client default language and returns a shallow copy of the client handle.
func (c* client) WithLanguage(language string) *client {
	c.options.Language = language
	return c
}

// SieriesByID returns a single  series information given a TVDB ID.
func (c *client) SeriesByID(id int) (*Series, error) {
	var data SeriesData
	fullURL := c.URL(fmt.Sprintf("/series/%d", id))
	if err := c.get(fullURL, &data, withLanguage(c.options.Language)); err != nil {
		return nil, err
	}
	return data.Data, nil
}

// EpisodesBySeriesID returns a list of episodes for the given series ID, optionally filtered by filters.
//
// When no filters are provided, this function returns all the episodes of the series.
// Filters can be used to get, for example, all the episodes from a given season, or a single episode in a season.
func (c *client) EpisodesBySeriesID(seriesID int,filters ...QueryOption) ([]Episode, error) {
	uri := fmt.Sprintf("/series/%d/episodes", seriesID)
	if len(filters) > 0 {
		uri += "/query"
	}

	var data SeriesEpisodesData
	var episodes []Episode

	last := 2
	for page := 1; page < last; page++ {
		filtersWithPage := append(filters, withQueryIntOption("page", page))
		fullURL := c.URL(uri, filtersWithPage...)
		if err := c.get(fullURL, &data, withLanguage(c.options.Language)); err != nil {
			return nil, err
		}
		last = data.Pages.Last

		episodes = append(episodes, data.Data...)
	}
	return episodes, nil
}

// SearchSeriesByName returns a list of series matching name.
//
// Results are returned in the configured client language.
func (c *client) SearchSeriesByName(name string) ([]SeriesSearchResult, error) {
	var result SeriesSearchResults
	fullURL := c.URL("/search/series", withQueryOption("name", name))
	err := c.get(fullURL, &result, withLanguage(c.options.Language))
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

func (c *client) get(url url.URL, out interface{}, options ...requestOption) error {
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
