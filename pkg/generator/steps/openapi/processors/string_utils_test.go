package processors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToApiFilename(t *testing.T) {
	assert.Equal(t, "api_path_to_api", toAPIFilename("/path/to/api"))
	assert.Equal(t, "api_path_to_api", toAPIFilename("/path/to/api/"))
	assert.Equal(t, "api_path_to_api_minus", toAPIFilename("/path/to/api-minus/"))
	assert.Equal(t, "api_path_to_api_param_test_other", toAPIFilename("/path/to/api/{param}.{test}/{other}"))

	// reserved filename
	assert.Equal(t, "api_api_v1_test_", toAPIFilename("/api/v1/test"))
}
