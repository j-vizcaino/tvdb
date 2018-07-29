package tvdb

type Token struct {
	Value string `json:"token,omitempty"`
	Error string `json:"Error,omitempty"`
}

type QueryrErrors struct {
	// Invalid filters passed to route
	InvalidFilters []string `json:"invalidFilters,omitempty"`
	// Invalid language or translation missing
	InvalidLanguage string `json:"invalidLanguage,omitempty"`
	// Invalid query params passed to route
	InvalidQueryParams []string `json:"invalidQueryParams,omitempty"`
}

type SeriesSearchResults struct {
	Data []SeriesSearchResult `json:"data"`
}

type SeriesSearchResult struct {
	Aliases    []string `json:"aliases"`
	Banner     string   `json:"banner"`
	FirstAired string   `json:"firstAired,omitempty"`
	ID         int      `json:"id,omitempty"`
	Network    string   `json:"network,omitempty"`
	Overview   string   `json:"overview,omitempty"`
	SeriesName string   `json:"seriesName,omitempty"`
	Slug       string   `json:"slug,omitempty"`
	Status     string   `json:"status,omitempty"`
}

type SeriesData struct {
	Data   *Series       `json:"data,omitempty"`
	Errors *QueryrErrors `json:"errors,omitempty"`
}

type Series struct {
	Added           string    `json:"added,omitempty"`
	AirsDayOfWeek   string    `json:"airsDayOfWeek,omitempty"`
	AirsTime        string    `json:"airsTime,omitempty"`
	Aliases         []string  `json:"aliases,omitempty"`
	Banner          string    `json:"banner,omitempty"`
	FirstAired      string    `json:"firstAired,omitempty"`
	Genre           []string  `json:"genre,omitempty"`
	ID              int       `json:"id,omitempty"`
	ImdbID          string    `json:"imdbId,omitempty"`
	LastUpdated     int		  `json:"lastUpdated,omitempty"`
	Network         string    `json:"network,omitempty"`
	NetworkID       string    `json:"networkId,omitempty"`
	Overview        string    `json:"overview,omitempty"`
	Rating          string    `json:"rating,omitempty"`
	Runtime         string    `json:"runtime,omitempty"`
	SeriesID        string    `json:"seriesId,omitempty"`
	SeriesName      string    `json:"seriesName,omitempty"`
	SiteRating      float32   `json:"siteRating,omitempty"`
	SiteRatingCount int       `json:"siteRatingCount,omitempty"`
	Slug            string    `json:"slug,omitempty"`
	Status          string    `json:"status,omitempty"`
	Zap2itID        string    `json:"zap2itId,omitempty"`
}

type SeriesEpisodesData struct {
	Data   []Episode    `json:"data,omitempty"`
	Errors QueryrErrors `json:"errors,omitempty"`
	Pages  Pages        `json:"links,omitempty"`
}

type Episode struct {
	AbsoluteNumber     int       `json:"absoluteNumber,omitempty"`
	AiredEpisodeNumber int       `json:"airedEpisodeNumber,omitempty"`
	AiredSeason        int       `json:"airedSeason,omitempty"`
	AirsAfterSeason    int       `json:"airsAfterSeason,omitempty"`
	AirsBeforeEpisode  int       `json:"airsBeforeEpisode,omitempty"`
	AirsBeforeSeason   int       `json:"airsBeforeSeason,omitempty"`
	Director           string    `json:"director,omitempty"`
	Directors          []string  `json:"directors,omitempty"`
	DVDChapter         int       `json:"dvdChapter,omitempty"`
	DVDDiscID          string    `json:"dvdDiscid,omitempty"`
	DVDEpisodeNumber   int       `json:"dvdEpisodeNumber,omitempty"`
	DVDSeason          int       `json:"dvdSeason,omitempty"`
	EpisodeName        string    `json:"episodeName,omitempty"`
	Filename           string    `json:"filename,omitempty"`
	FirstAired         string    `json:"firstAired,omitempty"`
	GuestStars         []string  `json:"guestStars,omitempty"`
	ID                 int       `json:"id,omitempty"`
	ImdbID             string    `json:"imdbId,omitempty"`
	LastUpdated        int       `json:"lastUpdated,omitempty"`
	LastUpdatedBy      int       `json:"lastUpdatedBy,omitempty"`
	Overview           string    `json:"overview,omitempty"`
	ProductionCode     string    `json:"productionCode,omitempty"`
	SeriesID           int       `json:"seriesId,omitempty"`
	ShowURL            string    `json:"showUrl,omitempty"`
	SiteRating         float32   `json:"siteRating,omitempty"`
	SiteRatingCount    int       `json:"siteRatingCount,omitempty"`
	ThumbAdded         string    `json:"thumbAdded,omitempty"`
	ThumbAuthor        int       `json:"thumbAuthor,omitempty"`
	ThumbHeight        string    `json:"thumbHeight,omitempty"`
	ThumbWidth         string    `json:"thumbWidth,omitempty"`
	Writers            []string  `json:"writers,omitempty"`
}

type Pages struct {
	First    int `json:"first,omitempty"`
	Last     int `json:"last,omitempty"`
	Next     int `json:"next,omitempty"`
	Previous int `json:"previous,omitempty"`
}

type LanguageData struct {
	Data  []Language `json:"data,omitempty"`
	Error string     `json:"Error,omitempty"`
}

type Language struct {
	Abbreviation string `json:"abbreviation,omitempty"`
	EnglishName  string `json:"englishName,omitempty"`
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
}
