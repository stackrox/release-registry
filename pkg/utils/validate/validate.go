package validate

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/stackrox/release-registry/pkg/utils/version"
)

// IsValidVersion checks if a tag is considered a valid version.
func IsValidVersion(tag string) error {
	if err := version.Validate(tag); err != nil {
		return errors.Wrapf(err, "tag %s is not a valid version", tag)
	}

	return nil
}

// IsValidCommit checks if a commit is considered valid.
func IsValidCommit(commit string) error {
	if !IsValidString(`^[0-9a-f]{40}`, commit) {
		return fmt.Errorf("commit %s is not a valid long Git SHA", commit)
	}

	return nil
}

// IsValidActorEmail checks if the approver has the expected email domain.
func IsValidActorEmail(validDomain, actor string) error {
	if !strings.HasSuffix(actor, validDomain) {
		return fmt.Errorf("actor %s has invalid email domain, expected %s", actor, validDomain)
	}

	return nil
}

// IsNotEmpty checks that a word is not empty.
func IsNotEmpty(word string) error {
	if len(word) == 0 {
		return errors.New("string is empty")
	}

	return nil
}
