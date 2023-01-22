package generate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTemplates(t *testing.T) {
	expectedTemplates := []string{"clips.html", "highlights.html", "home.html"}
	for _, et := range expectedTemplates {
		require.NotNil(t, templates.Lookup(et))
	}
}
