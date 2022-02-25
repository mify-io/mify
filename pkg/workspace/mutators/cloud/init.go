package cloud

import (
	"fmt"

	"github.com/mify-io/mify/pkg/workspace/mutators"
)

func Init(mutContext *mutators.MutatorContext) error {
	if err := UpdateCloudPublicity(mutContext); err != nil {
		return fmt.Errorf("failed to update cloud publicity: %w", err)
	}

	return nil
}
