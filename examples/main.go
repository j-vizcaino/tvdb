package main

import (
	"fmt"
	"github.com/j-vizcaino/tvdb"
)

func main() {
	// Creates a new TVDB client, performs login
	clt, err := tvdb.NewClient(tvdb.ClientOptions{
		APIKey:   "your-api-key",
		Language: "en",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	series, err := clt.SearchSeriesByName("The Simpsons")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Found %d serie(s) matching 'The Simpsons':\n", len(series))
	for _, s := range series {
		fmt.Printf("* %s (id: %d)\n", s.SeriesName, s.ID)
	}

	// NOTE: 71663 is The Simpsons series ID

	// Fetch all the episodes from the series
	fmt.Println("\nFetching all episodes...")
	episodes, err := clt.EpisodesBySeriesID(71663)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Series has %d episodes\n", len(episodes))
	for _, ep := range episodes {
		fmt.Printf("- S%02dE%02d: %s\n", ep.AiredSeason, ep.AiredEpisodeNumber, ep.EpisodeName)
	}

	fmt.Println("\nFetching all seasons specials...")
	// Fetch all the episodes from season specials (season 0)
	episodes, err = clt.EpisodesBySeriesID(71663, tvdb.WithAiredSeasonNumber(0))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Series has %d season specials\n", len(episodes))
	for _, ep := range episodes {
		fmt.Printf("- S%02dE%02d: %s\n", ep.AiredSeason, ep.AiredEpisodeNumber, ep.EpisodeName)
	}

	// Fetch all the episodes from season 8, episode 1
	episodes, err = clt.EpisodesBySeriesID(71663, tvdb.WithAiredSeasonNumber(8), tvdb.WithAiredEpisodeNumber(1))
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(episodes) < 1 {
		fmt.Printf("Unable to find season 8 episode 1")
		return
	}
	ep := episodes[0]
	fmt.Printf("\n== S%02dE%02d - %s\nOverview: %s\n", ep.AiredSeason, ep.AiredEpisodeNumber, ep.EpisodeName, ep.Overview)
}
