package integration

import (
	"testing"

	"github.com/mify-io/mify/internal/mify"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	basePath := t.TempDir()
	ctx := mify.NewContext()
	require.NoError(t, mify.CreateWorkspace(ctx, basePath, "workspace1"))
}
