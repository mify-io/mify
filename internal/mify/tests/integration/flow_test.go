package integration

import (
	"flag"
	"path"
	"testing"

	"github.com/mify-io/mify/internal/mify"
	"github.com/stretchr/testify/require"
)

var approve bool

func init() {
	flag.BoolVar(&approve, "approve", false, "Approve test result")
}

func TestFullFlow1(t *testing.T) {
	approval := NewApprovalContext(t)
	tempDir := t.TempDir()
	basePath := path.Join(tempDir, "workspace1")
	ctx := mify.NewContext(mify.Config{})

	approval.NewSubtest()
	require.NoError(t, mify.CreateWorkspace(ctx, tempDir, "workspace1"))
	approval.EndSubtest(tempDir)

	approval.NewSubtest()
	require.NoError(t, mify.CreateService(ctx, basePath, "go", "service1"))
	approval.EndSubtest(tempDir)

	approval.NewSubtest()
	require.NoError(t, mify.CreateService(ctx, basePath, "go", "service2"))
	approval.EndSubtest(tempDir)

	approval.NewSubtest()
	require.NoError(t, mify.AddClient(ctx, basePath, "service1", "service2"))
	approval.EndSubtest(tempDir)

	approval.NewSubtest()
	require.NoError(t, mify.RemoveClient(ctx, basePath, "service1", "service2"))
	approval.EndSubtest(tempDir)

	approval.NewSubtest()
	require.NoError(t, mify.AddClient(ctx, basePath, "service1", "service2"))
	approval.EndSubtest(tempDir)

	approval.NewSubtest()
	require.NoError(t, mify.CreateFrontend(ctx, basePath, "vue_js", "front"))
	approval.EndSubtest(tempDir)

	approval.NewSubtest()
	require.NoError(t, mify.AddClient(ctx, basePath, "front", "service1"))
	approval.EndSubtest(tempDir)

	if approve {
		approval.Approve()
	} else {
		approval.Verify()
	}
}
