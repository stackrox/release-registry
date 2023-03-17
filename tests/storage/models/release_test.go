package models_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/tests"
	"github.com/stretchr/testify/assert"
)

const defaultTag = "3.74.0"
const defaultCommit = "b1d4c6264309de1da809dc85ed0825f817c58d8d"
const defaultCreator = "roxbot@redhat.com"

func setupReleaseTest(t *testing.T) {
	t.Helper()

	err := tests.SetupDB()
	assert.NoError(t, err)
	err = tests.Migrate(
		&models.Release{},
		&models.ReleaseMetadata{},
		&models.QualityMilestoneMetadata{},
		&models.QualityMilestoneDefinition{},
		&models.QualityMilestone{},
	)
	assert.NoError(t, err)
}

func createFakeRelease(t *testing.T) models.Release {
	t.Helper()

	release, err := models.CreateRelease(
		*configuration.New(),
		defaultTag, defaultCommit, defaultCreator,
		[]models.ReleaseMetadata{},
	)
	assert.NoError(t, err)

	return *release
}

func createMultipleFakeReleases(t *testing.T) []models.Release {
	t.Helper()

	expectedReleases := []models.Release{
		{
			Tag:      "1.0.0",
			Commit:   "b1d4c6264309de1da809dc85ed0825f817c58d8d",
			Creator:  "roxbot@redhat.com",
			Metadata: []models.ReleaseMetadata{{Key: "Key1", Value: "Value1"}, {Key: "Key2", Value: "Value2"}},
		}, {
			Tag:      "1.0.1",
			Commit:   "c289b8587a56462d7d64682053171ab69f5c5202",
			Creator:  "roxbot@redhat.com",
			Metadata: []models.ReleaseMetadata{{Key: "Key1", Value: "Value1"}, {Key: "Key2", Value: "Value2"}},
		}, {
			Tag:      "2.0.0",
			Commit:   "e4280c38e2bbb53cd60444e490ce0ea35f1b339c",
			Creator:  "roxbot@redhat.com",
			Metadata: []models.ReleaseMetadata{{Key: "Key1", Value: "Value1"}, {Key: "Key2", Value: "Value2"}},
		},
	}

	expectedReleaseDatabaseObjects := []models.Release{}

	for _, release := range expectedReleases {
		releaseDBO, err := models.CreateRelease(
			*configuration.New(),
			release.Tag, release.Commit,
			release.Creator, release.Metadata,
		)
		assert.NoError(t, err)

		expectedReleaseDatabaseObjects = append(expectedReleaseDatabaseObjects, *releaseDBO)
	}

	return expectedReleaseDatabaseObjects
}

func assertReleasesAreEqual(t *testing.T, expected, actual models.Release, checkPreload bool) {
	t.Helper()
	assert.Equal(t, expected.Commit, actual.Commit)
	assert.Equal(t, expected.Tag, actual.Tag)
	assert.Equal(t, expected.Creator, actual.Creator)
	assert.Equal(t, expected.CreatedAt.UnixNano(), actual.CreatedAt.UnixNano())

	if checkPreload {
		assert.Equal(t, expected.Metadata[0].Key, actual.Metadata[0].Key)
		assert.Equal(t, expected.Metadata[0].Value, actual.Metadata[0].Value)
		assert.Equal(t, expected.Metadata[1].Key, actual.Metadata[1].Key)
		assert.Equal(t, expected.Metadata[1].Value, actual.Metadata[1].Value)
	} else {
		assert.Equal(t, []models.ReleaseMetadata(nil), actual.Metadata)
	}
}

func TestCreateRelease(t *testing.T) {
	setupReleaseTest(t)

	release := createFakeRelease(t)
	assert.Equal(t, release.Tag, defaultTag)
	assert.Equal(t, release.Commit, defaultCommit)
	assert.Equal(t, release.Creator, defaultCreator)
	assert.Equal(t, release.Metadata, []models.ReleaseMetadata{})

	_, err := models.CreateRelease(
		*configuration.New(),
		"1.2.3.4.5.6", defaultCommit, defaultCreator, []models.ReleaseMetadata{},
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "tag is not a valid SemVer")
}

func TestRejectRelease(t *testing.T) {
	setupReleaseTest(t)

	release := createFakeRelease(t)
	updatedRelease, err := models.RejectRelease(release.Tag, false)
	assert.NoError(t, err)
	assert.Equal(t, release.Tag, updatedRelease.Tag)
	assert.Equal(t, updatedRelease.Rejected, true)

	_, err = models.GetRelease(release.Tag, false, false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "record not found")

	hiddenRelease, err := models.GetRelease(release.Tag, false, true)
	assert.NoError(t, err)
	assert.Equal(t, release.Tag, hiddenRelease.Tag)

	// Reject unknown release results in error
	_, err = models.RejectRelease("unknown tag", false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "record not found")

	// Reject rejected release results in error
	_, err = models.RejectRelease(release.Tag, false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "release not found or already rejected")
}

func TestGetReleaseByTag(t *testing.T) {
	setupReleaseTest(t)

	metadata := []models.ReleaseMetadata{{Key: "Key1", Value: "Value1"}, {Key: "Key2", Value: "Value2"}}
	originalRelease, err := models.CreateRelease(*configuration.New(), defaultTag, defaultCommit, defaultCreator, metadata)
	assert.NoError(t, err)

	// Get a release without preloading metadata
	retrievedRelease, err := models.GetRelease(defaultTag, false, false)
	assert.NoError(t, err)
	assertReleasesAreEqual(t, *originalRelease, *retrievedRelease, false)

	// Get a release with preloading metadata
	retrievedRelease, err = models.GetRelease(defaultTag, true, false)
	assert.NoError(t, err)
	assertReleasesAreEqual(t, *originalRelease, *retrievedRelease, true)

	// Get an unknown release returns an error
	_, err = models.GetRelease("unknown tag", false, false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "record not found")
}

func TestGetReleaseByTagWithQualityMilestones(t *testing.T) {
	setupReleaseTest(t)

	config := *configuration.New()

	release := createFakeRelease(t)
	qmd := createFakeQualityMilestoneDefinition(t)
	metadata := []models.QualityMilestoneMetadata{
		{Key: "Abc", Value: "abc"},
		{Key: "Def", Value: "def"},
		{Key: "Ghi", Value: "ghi"},
	}

	_, err := models.ApproveQualityMilestone(config, release.Tag, qmd.Name, "roxbot@redhat.com", metadata)
	assert.NoError(t, err)

	releases, err := models.ListAllReleasesAtQualityMilestone(qmd.Name, true, false)
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, releases[0].Tag, release.Tag)

	actualRelease, err := models.GetRelease(defaultTag, true, false)
	assert.NoError(t, err)
	assert.Equal(t, actualRelease.Tag, release.Tag)
	assert.Equal(t, len(actualRelease.Metadata), 0)

	assert.Equal(t, 1, len(actualRelease.QualityMilestones))
}

func TestListAllReleases(t *testing.T) {
	setupReleaseTest(t)

	expectedReleases := createMultipleFakeReleases(t)

	actualReleases, err := models.ListAllReleases(true, false)
	assert.NoError(t, err)

	assert.Len(t, actualReleases, 3)

	for i, expectedRelease := range expectedReleases {
		assertReleasesAreEqual(t, expectedRelease, actualReleases[i], true)
	}
}

func TestListAllReleasesWithPrefix(t *testing.T) {
	setupReleaseTest(t)

	expectedReleaseDatabaseObjects := createMultipleFakeReleases(t)

	// Expect only 2 releases due to the third release having the wrong prefix
	actualReleases, err := models.ListAllReleasesWithPrefix("1.0", true, false)
	assert.NoError(t, err)

	assert.Len(t, actualReleases, 2)
	assertReleasesAreEqual(t, expectedReleaseDatabaseObjects[0], actualReleases[0], true)
	assertReleasesAreEqual(t, expectedReleaseDatabaseObjects[1], actualReleases[1], true)
}

func TestListAllReleasesAtQualityMilestone(t *testing.T) {
	setupReleaseTest(t)

	config := *configuration.New()

	_, err := models.CreateRelease(
		config, "2.0.0", "b1d4c6264309de1da809dc85ed0825f817c58d8d", "roxbot@redhat.com",
		[]models.ReleaseMetadata{},
	)
	assert.NoError(t, err)

	release := createFakeRelease(t)
	qmd := createFakeQualityMilestoneDefinition(t)
	metadata := []models.QualityMilestoneMetadata{
		{Key: "Abc", Value: "abc"},
		{Key: "Def", Value: "def"},
		{Key: "Ghi", Value: "ghi"},
	}
	_, err = models.ApproveQualityMilestone(config, release.Tag, qmd.Name, "roxbot@redhat.com", metadata)
	assert.NoError(t, err)

	// Expect only one release, due to other one not approved for QualityMilestone
	releases, err := models.ListAllReleasesAtQualityMilestone(qmd.Name, true, false)
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, releases[0].Tag, release.Tag)
}

func TestListAllReleasesAtQualityMilestoneWithPrefix(t *testing.T) {
	setupReleaseTest(t)

	config := *configuration.New()
	qmd := createFakeQualityMilestoneDefinition(t)

	fakeMetadata := []models.QualityMilestoneMetadata{
		{Key: "Abc", Value: "abc"},
		{Key: "Def", Value: "def"},
		{Key: "Ghi", Value: "ghi"},
	}
	prefixedRelease, err := models.CreateRelease(
		config, "1.0.0", "b1d4c6264309de1da809dc85ed0825f817c58d8d", "roxbot@redhat.com",
		[]models.ReleaseMetadata{},
	)
	assert.NoError(t, err)

	_, err = models.CreateRelease(
		config, "2.0.0", "b1d4c6264309de1da809dc85ed0825f817c58d8d", "roxbot@redhat.com",
		[]models.ReleaseMetadata{},
	)
	assert.NoError(t, err)

	_, err = models.ApproveQualityMilestone(config, "1.0.0", qmd.Name, "roxbot@redhat.com", fakeMetadata)
	assert.NoError(t, err)
	_, err = models.ApproveQualityMilestone(config, "2.0.0", qmd.Name, "roxbot@redhat.com", fakeMetadata)
	assert.NoError(t, err)

	// Expect only 1 release due to 2.0.0 not having the correct prefix
	releases, err := models.ListAllReleasesWithPrefixAtQualityMilestone(prefixedRelease.Tag, qmd.Name, false, false)
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, releases[0].Tag, prefixedRelease.Tag)
}

func TestFindLatestRelease(t *testing.T) {
	setupReleaseTest(t)

	expectedReleases := createMultipleFakeReleases(t)

	latest, err := models.FindLatestRelease(false, false)
	assert.NoError(t, err)
	assertReleasesAreEqual(t, expectedReleases[2], *latest, false)

	// TODO: what if no releases ?
}

func TestFindLatestReleaseWithPrefix(t *testing.T) {}

func TestFindLatestReleaseAtQualityMilestone(t *testing.T) {}

func TestFindLatestRelaseWithPrefixAtQualityMilestone(t *testing.T) {}
