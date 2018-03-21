package gradle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractArtifactName(t *testing.T) {
	proj := Project{
		location: "/root_dir",
		monoRepo: false,
	}
	// non-monorepo
	got, err := proj.extractArtifactName(
		"/root_dir/mymodule/build/reports/myartifact.html")

	require.NoError(t, err)
	require.Equal(t, "mymodule-myartifact.html", got)

	// monorepo
	proj = Project{
		location: "/root_dir",
		monoRepo: true,
	}
	got, err = proj.extractArtifactName(
		"/root_dir/mymodule/build/reports/myartifact.html")

	require.NoError(t, err)
	require.Equal(t, "root_dir-mymodule-myartifact.html", got)

	// in root
	proj = Project{
		location: "/",
		monoRepo: false,
	}
	got, err = proj.extractArtifactName(
		"/mymodule/build/reports/myartifact.html")

	require.NoError(t, err)
	require.Equal(t, "mymodule-myartifact.html", got)

	// in root monorepo
	proj = Project{
		location: "/",
		monoRepo: true,
	}
	got, err = proj.extractArtifactName(
		"/mymodule/build/reports/myartifact.html")

	require.NoError(t, err)
	require.Equal(t, "mymodule-myartifact.html", got)
}
