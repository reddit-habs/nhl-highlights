package repository

import (
	"testing"

	"github.com/sbstp/nhl-highlights/addrof"
	"github.com/stretchr/testify/require"
)

func TestGames(t *testing.T) {
	r, err := New(":memory:")
	require.NoError(t, err)

	// Insert
	err = r.UpsertGame(&Game{
		GameID: 20210629,
		Date:   "2021-06-29",
		Type:   "R",
		Away:   "MTL",
		Home:   "TBL",
		Season: "20212022",
	})
	require.NoError(t, err)

	// Select missing
	games, err := r.GetGamesMissingContent(false)
	require.NoError(t, err)
	require.Equal(t, 1, len(games))

	// Select
	game, err := r.GetGame(games[0].GameID)
	require.Equal(t, games[0], game)

	// Update
	updated := &Game{
		GameID:   20210629,
		Date:     "2022-06-29",
		Type:     "P",
		Away:     "SJS",
		Home:     "LAK",
		Season:   "20212023",
		Recap:    addrof.String("recap"),
		Extended: addrof.String("extended"),
	}
	err = r.UpsertGame(updated)
	require.NoError(t, err)

	// Select
	game2, err := r.GetGame(games[0].GameID)
	require.NoError(t, err)
	require.Equal(t, updated, game2)
}
