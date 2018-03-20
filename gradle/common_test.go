package gradle

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanStringSlice(t *testing.T) {
	require.Equal(t, true, reflect.DeepEqual(cleanStringSlice([]string{"", ""}), []string{}))
	require.Equal(t, true, reflect.DeepEqual(cleanStringSlice([]string{"", "test"}), []string{"test"}))
	require.Equal(t, true, reflect.DeepEqual(cleanStringSlice([]string{"", "   space "}), []string{"space"}))
}
