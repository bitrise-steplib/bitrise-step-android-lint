package gradle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractArtifactName(t *testing.T) {
	// non-monorepo
	got, err := extractArtifactName(
		Project{
			location: "root_dir",
			monoRepo: false,
		},
		"root_dir/mymodule/build/reports/myartifact.html")

	require.NoError(t, err)
	require.Equal(t, "mymodule-myartifact.html", got)

	// monorepo
	got, err = extractArtifactName(
		Project{
			location: "root_dir",
			monoRepo: true,
		},
		"root_dir/mymodule/build/reports/myartifact.html")

	require.NoError(t, err)
	require.Equal(t, "root_dir-mymodule-myartifact.html", got)
}
