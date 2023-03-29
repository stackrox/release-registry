package validate_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/utils/validate"
	"github.com/stretchr/testify/assert"
)

func TestIsValidString(t *testing.T) {
	assert.True(t, validate.IsValidString(`^[a-z]*$`, "foo"))
	assert.False(t, validate.IsValidString(`^[a-z]*$`, "Bar"))
}
