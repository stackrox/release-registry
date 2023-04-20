package validate_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/utils/validate"
	"github.com/stretchr/testify/assert"
)

func TestValidateActorHasValidEmail(t *testing.T) {
	config := configuration.New("../../../example")
	config.Tenant.EmailDomain = "@redhat.com"

	// Invalid email
	err := validate.IsValidActorEmail(config.Tenant.EmailDomain, "roxbot@stackrox.com")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "actor roxbot@stackrox.com has invalid email domain, expected @redhat.com")

	// Valid email
	err = validate.IsValidActorEmail(config.Tenant.EmailDomain, "roxbot@redhat.com")
	assert.NoError(t, err)
}
