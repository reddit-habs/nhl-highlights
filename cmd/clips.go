package main

import (
	"sort"
	"strconv"

	"github.com/sbstp/nhl-highlights/generate"
	"github.com/sbstp/nhl-highlights/nhlapi"
	"github.com/volatiletech/null/v8"
)

func scanClips(api nhlapi.Client, gameID int64) ([]*generate.Clip, error) {
	var clips []*generate.Clip
	content, err := api.Content(gameID)
	if err != nil {
		return nil, err
	}

	for _, item := range content.Highlights.Scoreboard.Items {
		video := getBestMp4Playback(item.Playbacks)

		var eventID null.Int64
		for _, kw := range item.Keywords {
			if kw.Type == "statsEventId" {
				eventID = null.Int64From(stringToInt64(kw.Value))
			}
		}

		clips = append(clips, &generate.Clip{
			ID:          stringToInt64(item.ID),
			GameID:      gameID,
			MediaURL:    video,
			EventID:     eventID.Ptr(),
			Title:       item.Title,
			Blurb:       item.Blurb,
			Description: item.Description,
		})
	}

	sort.Slice(clips, func(i, j int) bool {
		a := clips[i]
		b := clips[j]
		if a.EventID == nil {
			return false
		}
		if b.EventID == nil {
			return true
		}
		return *a.EventID < *b.EventID
	})

	return clips, nil
}

func stringToInt64(s string) int64 {
	x, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return x
}
