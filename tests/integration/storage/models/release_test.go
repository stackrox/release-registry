package models_test

import (
	"testing"

	"github.com/stackrox/release-registry/pkg/configuration"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/version"
	"github.com/stackrox/release-registry/tests/integration"
	"github.com/stackrox/release-registry/tests/utils"
	"github.com/stretchr/testify/assert"
)

const defaultTag = "3.74.0"
const defaultCommit = "b1d4c6264309de1da809dc85ed0825f817c58d8d"
const defaultCreator = "roxbot@redhat.com"

//nolint:gochecknoglobals
var defaultQualityMilestoneMetadata = []models.QualityMilestoneMetadata{
	{Key: "Abc", Value: "abc"},
	{Key: "Def", Value: "def"},
	{Key: "Ghi", Value: "ghi"},
}

//nolint:gochecknoglobals
var defaultReleaseMetadata = []models.ReleaseMetadata{{Key: "Key1", Value: "Value1"}, {Key: "Key2", Value: "Value2"}}

func setupReleaseTest(t *testing.T) {
	t.Helper()

	err := integration.SetupDB()
	assert.NoError(t, err)
	err = integration.Migrate(
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
		configuration.New().Tenant.EmailDomain,
		defaultTag, defaultCommit, defaultCreator,
		defaultReleaseMetadata,
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
		}, {
			Tag:      "0.1.x-14-gabcdef1234",
			Commit:   "e4280c38e2bbb53cd60444e490ce0ea35f1b339c",
			Creator:  "roxbot@redhat.com",
			Metadata: []models.ReleaseMetadata{{Key: "Key1", Value: "Value1"}, {Key: "Key2", Value: "Value2"}},
		}, {
			Tag:      "0.1.1-rc.1",
			Commit:   "e4280c38e2bbb53cd60444e490ce0ea35f1b339c",
			Creator:  "roxbot@redhat.com",
			Metadata: []models.ReleaseMetadata{{Key: "Key1", Value: "Value1"}, {Key: "Key2", Value: "Value2"}},
		}, {
			Tag:      "0.1.x-nightly-20230508",
			Commit:   "e4280c38e2bbb53cd60444e490ce0ea35f1b339c",
			Creator:  "roxbot@redhat.com",
			Metadata: []models.ReleaseMetadata{{Key: "Key1", Value: "Value1"}, {Key: "Key2", Value: "Value2"}},
		},
	}

	expectedReleaseDatabaseObjects := []models.Release{}

	for _, release := range expectedReleases {
		releaseDBO, err := models.CreateRelease(
			configuration.New().Tenant.EmailDomain,
			release.Tag, release.Commit,
			release.Creator, release.Metadata,
		)
		assert.NoError(t, err)

		expectedReleaseDatabaseObjects = append(expectedReleaseDatabaseObjects, *releaseDBO)
	}

	return expectedReleaseDatabaseObjects
}

func TestCreateRelease(t *testing.T) {
	setupReleaseTest(t)

	release := createFakeRelease(t)
	assert.Equal(t, release.Tag, defaultTag)
	assert.Equal(t, release.Commit, defaultCommit)
	assert.Equal(t, release.Creator, defaultCreator)
	assert.Equal(t, release.Metadata, defaultReleaseMetadata)
}

func TestCreateReleaseInvalidSemVer(t *testing.T) {
	setupReleaseTest(t)

	_, err := models.CreateRelease(
		configuration.New().Tenant.EmailDomain,
		"1.2.3aszihiuh", defaultCommit, defaultCreator, []models.ReleaseMetadata{},
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "tag 1.2.3aszihiuh is not a valid version")
}

func TestCreateReleaseNightlyVersion(t *testing.T) {
	setupReleaseTest(t)

	nightlyTag := "3.74.x-nightly-20230320"
	release, err := models.CreateRelease(
		configuration.New().Tenant.EmailDomain,
		nightlyTag, defaultCommit, defaultCreator, []models.ReleaseMetadata{},
	)
	assert.NoError(t, err)
	assert.Equal(t, nightlyTag, release.Tag)
}

func TestUpdateRelease(t *testing.T) {
	setupReleaseTest(t)

	release := createFakeRelease(t)
	updatedMetadata := []models.ReleaseMetadata{
		{Key: "Key1", Value: "Value1"},
		{Key: "Key2", Value: "Value2"},
		{Key: "Key3", Value: "Value3"},
	}
	updatedRelease, err := models.UpdateRelease(release.Tag, updatedMetadata, false)
	assert.NoError(t, err)
	assert.Len(t, updatedRelease.Metadata, len(updatedMetadata))

	actualUpdatedRelease, err := models.GetRelease(release.Tag, true, false)
	assert.NoError(t, err)
	assert.Len(t, actualUpdatedRelease.Metadata, len(updatedMetadata))
	utils.AssertReleasesAreEqual(t, updatedRelease, actualUpdatedRelease, true, true)
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
}

func TestRejectUnknownReleaseError(t *testing.T) {
	setupReleaseTest(t)

	_, err := models.RejectRelease("1.1.1", false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "record not found")
}

func TestRejectRejectedReleaseError(t *testing.T) {
	setupReleaseTest(t)

	release := createFakeRelease(t)
	_, err := models.RejectRelease(release.Tag, false)
	assert.NoError(t, err)

	_, err = models.RejectRelease(release.Tag, false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "release not found or already rejected")
}

func TestGetReleaseByTag(t *testing.T) {
	setupReleaseTest(t)

	metadata := []models.ReleaseMetadata{{Key: "Key1", Value: "Value1"}, {Key: "Key2", Value: "Value2"}}
	originalRelease, err := models.CreateRelease(
		configuration.New().Tenant.EmailDomain,
		defaultTag, defaultCommit, defaultCreator,
		metadata,
	)
	assert.NoError(t, err)

	// Get a release without preloading metadata
	retrievedRelease, err := models.GetRelease(defaultTag, false, false)
	assert.NoError(t, err)
	utils.AssertReleasesAreEqual(t, originalRelease, retrievedRelease, true, false)
	assert.Nil(t, retrievedRelease.Metadata)

	// Get a release with preloading metadata
	retrievedRelease, err = models.GetRelease(defaultTag, true, false)
	assert.NoError(t, err)
	utils.AssertReleasesAreEqual(t, originalRelease, retrievedRelease, true, true)

	// Get an unknown release returns an error
	_, err = models.GetRelease("1.1.1", false, false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "record not found")
}

func TestGetReleaseByTagWithQualityMilestones(t *testing.T) {
	setupReleaseTest(t)

	config := configuration.New()

	release := createFakeRelease(t)
	qmd := createFakeQualityMilestoneDefinition(t)
	metadata := []models.QualityMilestoneMetadata{
		{Key: "Abc", Value: "abc"},
		{Key: "Def", Value: "def"},
		{Key: "Ghi", Value: "ghi"},
	}

	_, err := models.ApproveQualityMilestone(
		config.Tenant.EmailDomain,
		release.Tag, qmd.Name, "roxbot@redhat.com",
		metadata,
	)
	assert.NoError(t, err)

	releases, err := models.ListAllReleasesAtQualityMilestone(qmd.Name, []version.Kind{}, true, false)
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	utils.AssertReleasesAreEqual(t, &releases[0], &release, true, true)

	actualRelease, err := models.GetRelease(defaultTag, true, false)
	assert.NoError(t, err)
	assert.Equal(t, actualRelease.Tag, release.Tag)
	assert.Len(t, actualRelease.Metadata, 2)
	assert.Len(t, actualRelease.QualityMilestones, 1)
}

func TestListAllReleasesWithWithoutRejected(t *testing.T) {
	setupReleaseTest(t)

	expectedReleases := createMultipleFakeReleases(t)
	_, err := models.RejectRelease(expectedReleases[0].Tag, false)
	assert.NoError(t, err)

	releasesWithoutRejected, err := models.ListAllReleases([]version.Kind{}, true, false)
	assert.NoError(t, err)
	assert.Len(t, releasesWithoutRejected, 5)

	// expectedReleases[0] is rejected, don't compare that one.
	utils.AssertReleasesAreEqual(t, &expectedReleases[1], &releasesWithoutRejected[0], false, true)
	utils.AssertReleasesAreEqual(t, &expectedReleases[2], &releasesWithoutRejected[1], false, true)

	releasesWithRejected, err := models.ListAllReleases([]version.Kind{}, true, true)
	assert.NoError(t, err)
	assert.Len(t, releasesWithRejected, 6)

	for i := range expectedReleases {
		utils.AssertReleasesAreEqual(t, &expectedReleases[i], &releasesWithRejected[i], false, true)
	}
}

func TestListAllReleasesWithPrefix(t *testing.T) {
	setupReleaseTest(t)

	expectedReleaseDatabaseObjects := createMultipleFakeReleases(t)

	// Expect only 2 releases due to the third release having the wrong prefix
	actualReleases, err := models.ListAllReleasesWithPrefix("1.0", []version.Kind{}, true, false)
	assert.NoError(t, err)

	assert.Len(t, actualReleases, 2)
	utils.AssertReleasesAreEqual(t, &expectedReleaseDatabaseObjects[0], &actualReleases[0], false, true)
	utils.AssertReleasesAreEqual(t, &expectedReleaseDatabaseObjects[1], &actualReleases[1], false, true)
}

func TestListAllReleasesAtQualityMilestone(t *testing.T) {
	setupReleaseTest(t)

	config := configuration.New()

	_, err := models.CreateRelease(
		config.Tenant.EmailDomain, "2.0.0", "b1d4c6264309de1da809dc85ed0825f817c58d8d", "roxbot@redhat.com",
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
	_, err = models.ApproveQualityMilestone(
		config.Tenant.EmailDomain,
		release.Tag, qmd.Name, "roxbot@redhat.com",
		metadata,
	)
	assert.NoError(t, err)

	// Expect only one release, due to other one not approved for QualityMilestone
	releases, err := models.ListAllReleasesAtQualityMilestone(qmd.Name, []version.Kind{}, true, false)
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	utils.AssertReleasesAreEqual(t, &release, &releases[0], true, false)
}

func TestListAllReleasesAtQualityMilestoneWithPrefix(t *testing.T) {
	setupReleaseTest(t)

	config := configuration.New()
	qmd := createFakeQualityMilestoneDefinition(t)

	prefixedRelease, err := models.CreateRelease(
		config.Tenant.EmailDomain, "1.0.0", "b1d4c6264309de1da809dc85ed0825f817c58d8d", "roxbot@redhat.com",
		[]models.ReleaseMetadata{},
	)
	assert.NoError(t, err)

	_, err = models.CreateRelease(
		config.Tenant.EmailDomain, "2.0.0", "b1d4c6264309de1da809dc85ed0825f817c58d8d", "roxbot@redhat.com",
		[]models.ReleaseMetadata{},
	)
	assert.NoError(t, err)

	_, err = models.ApproveQualityMilestone(
		config.Tenant.EmailDomain, "1.0.0", qmd.Name,
		"roxbot@redhat.com", defaultQualityMilestoneMetadata,
	)
	assert.NoError(t, err)
	_, err = models.ApproveQualityMilestone(
		config.Tenant.EmailDomain, "2.0.0", qmd.Name,
		"roxbot@redhat.com", defaultQualityMilestoneMetadata,
	)
	assert.NoError(t, err)

	// Expect only 1 release due to 2.0.0 not having the correct prefix
	releases, err := models.ListAllReleasesWithPrefixAtQualityMilestone(
		prefixedRelease.Tag, qmd.Name,
		[]version.Kind{}, false, false,
	)
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	utils.AssertReleasesAreEqual(t, prefixedRelease, &releases[0], true, false)
	assert.Equal(t, releases[0].Tag, prefixedRelease.Tag)
}

func TestFindLatestRelease(t *testing.T) {
	setupReleaseTest(t)

	expectedReleases := createMultipleFakeReleases(t)

	latest, err := models.FindLatestRelease([]version.Kind{}, false, false)
	assert.NoError(t, err)
	utils.AssertReleasesAreEqual(t, &expectedReleases[2], latest, false, false)
}

func TestFindLatestReleasesNoReleases(t *testing.T) {
	setupReleaseTest(t)

	_, err := models.FindLatestRelease([]version.Kind{}, false, false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "no releases found")
}

func TestFindLatestReleasesNightlies(t *testing.T) {
	setupReleaseTest(t)

	config := configuration.New()
	lastNightly, err := models.CreateRelease(
		config.Tenant.EmailDomain,
		"3.74.x-nightly-20230321",
		defaultCommit, defaultCreator, []models.ReleaseMetadata{},
	)
	assert.NoError(t, err)

	_, err = models.CreateRelease(
		config.Tenant.EmailDomain,
		"3.74.x-nightly-20230319",
		defaultCommit, defaultCreator, []models.ReleaseMetadata{},
	)
	assert.NoError(t, err)

	_, err = models.CreateRelease(
		config.Tenant.EmailDomain,
		"3.74.x-nightly-20230320",
		defaultCommit, defaultCreator, []models.ReleaseMetadata{},
	)
	assert.NoError(t, err)

	latestRetrieved, err := models.FindLatestRelease([]version.Kind{}, false, false)
	assert.NoError(t, err)
	utils.AssertReleasesAreEqual(t, lastNightly, latestRetrieved, true, true)
}

func TestFindLatestReleaseWithPrefix(t *testing.T) {
	setupReleaseTest(t)

	expectedReleases := createMultipleFakeReleases(t)

	latest, err := models.FindLatestReleaseWithPrefix("1.0", []version.Kind{}, true, false)
	assert.NoError(t, err)
	utils.AssertReleasesAreEqual(t, &expectedReleases[1], latest, false, true)
}

func TestFindLatestReleaseAtQualityMilestone(t *testing.T) {
	setupReleaseTest(t)

	expectedReleases := createMultipleFakeReleases(t)
	qmd := createFakeQualityMilestoneDefinition(t)

	_, err := models.ApproveQualityMilestone(
		configuration.New().Tenant.EmailDomain,
		expectedReleases[0].Tag, qmd.Name,
		"roxbot@redhat.com", defaultQualityMilestoneMetadata,
	)
	assert.NoError(t, err)

	latest, err := models.FindLatestReleaseAtQualityMilestone("QM1", []version.Kind{}, true, false)
	assert.NoError(t, err)
	utils.AssertReleasesAreEqual(t, &expectedReleases[0], latest, false, true)
}

func TestFindLatestRelaseWithPrefixAtQualityMilestone(t *testing.T) {
	setupReleaseTest(t)

	expectedReleases := createMultipleFakeReleases(t)
	qmd := createFakeQualityMilestoneDefinition(t)

	// Approve both 1.x releases, expect 1.0.1 to be latest
	_, err := models.ApproveQualityMilestone(
		configuration.New().Tenant.EmailDomain,
		expectedReleases[0].Tag, qmd.Name,
		"roxbot@redhat.com", defaultQualityMilestoneMetadata,
	)
	assert.NoError(t, err)

	_, err = models.ApproveQualityMilestone(
		configuration.New().Tenant.EmailDomain,
		expectedReleases[1].Tag, qmd.Name,
		"roxbot@redhat.com", defaultQualityMilestoneMetadata,
	)
	assert.NoError(t, err)

	latest, err := models.FindLatestRelaseWithPrefixAtQualityMilestone("1.0", "QM1", []version.Kind{}, true, false)
	assert.NoError(t, err)
	utils.AssertReleasesAreEqual(t, &expectedReleases[1], latest, false, true)
}

func TestReleaseKindExclusion(t *testing.T) {
	setupReleaseTest(t)
	expectedReleases := createMultipleFakeReleases(t)

	allReleases, err := models.ListAllReleases([]version.Kind{}, false, false)
	assert.NoError(t, err)
	assert.Len(t, allReleases, 6)

	ignoredKinds := []version.Kind{version.ReleaseKind}
	withExcludedReleases, err := models.ListAllReleases(ignoredKinds, false, false)
	assert.NoError(t, err)
	assert.Len(t, withExcludedReleases, 3)

	for i := range withExcludedReleases {
		utils.AssertReleasesAreEqual(t, &expectedReleases[i+3], &withExcludedReleases[i], false, false)
	}

	ignoredKinds = []version.Kind{version.RCKind, version.NightlyKind}
	withExcludedReleases, err = models.ListAllReleases(ignoredKinds, false, false)
	assert.NoError(t, err)
	assert.Len(t, withExcludedReleases, 4)

	for i := range withExcludedReleases {
		utils.AssertReleasesAreEqual(t, &expectedReleases[i], &withExcludedReleases[i], false, false)
	}
}
