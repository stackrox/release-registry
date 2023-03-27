package lifecycle_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/conversions"
	"github.com/stackrox/release-registry/tests/e2e"
	"github.com/stackrox/release-registry/tests/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Use case: Support on-call engineer as the reviewer of the nightly marks a tag as successful.
// Steps:
// 1. Create a new Release.
// 2. List all Releases with prefix.
// 3. List all QualityMilestoneDefinitions.
// 4. Approve Release for QualityMilestone.
// 5. List all Releases with prefix at QualityMilestone.
//
//nolint:funlen
func TestReleasesCanBeCreatedAndApproved(t *testing.T) {
	qualityMilestoneDefinitionName := "Nightly has passed"
	user := "roxbot@redhat.com"

	expectedRelease := &models.Release{
		Tag:     "3.74.x-nightly-20230323",
		Commit:  "78057dba490233f41b4602f2b2e88775ab7fd4c9",
		Creator: "roxbot@redhat.com",
		Metadata: []models.ReleaseMetadata{
			{Key: "Link", Value: "https://github.com/stackrox/stackrox/releases/tag/3.74.x-nightly-20230323"},
		},
	}

	basePath, err := e2e.GetFixturesPath()
	assert.NoError(t, err)

	dbPath := fmt.Sprintf("%s/%s", basePath, "prepopulated-with-qualitymilestonedefinitions.sqlite")
	e2e.SetupE2ETest(t, dbPath)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)

	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(e2e.BufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)

	//nolint:errcheck
	defer conn.Close()

	// 1. Create a new Release
	releaseClient := v1.NewReleaseServiceClient(conn)
	createReleaseResponse, err := releaseClient.Create(
		ctx, conversions.NewCreateReleaseRequestFromRelease(expectedRelease),
	)
	assert.NoError(t, err)

	actualCreatedRelease := conversions.NewReleaseFromCreateReleaseResponse(createReleaseResponse)
	utils.AssertReleasesAreEqual(t, expectedRelease, actualCreatedRelease, false, true)

	// 2. List all Releases with prefix
	prefix := expectedRelease.Tag[:4]
	releaseList, err := releaseClient.List(ctx, &v1.ReleaseServiceListRequest{
		Prefix:  &prefix,
		Preload: true,
	})
	assert.NoError(t, err)
	assert.Len(t, releaseList.GetReleases(), 1)

	actualListedRelease := conversions.NewReleaseFromGetReleaseResponse(releaseList.GetReleases()[0])
	utils.AssertReleasesAreEqual(t, actualCreatedRelease, actualListedRelease, true, true)

	// 3. List all QualityMilestoneDefinitions
	expectedMetadataKeys := []string{"Image", "BuildURL"}
	qmdClient := v1.NewQualityMilestoneDefinitionServiceClient(conn)
	qmdList, err := qmdClient.List(ctx, &v1.QualityMilestoneDefinitionServiceListRequest{})
	assert.NoError(t, err)
	assert.Len(t, qmdList.GetQualityMilestoneDefinitions(), 2)
	assert.Equal(t, qualityMilestoneDefinitionName, qmdList.GetQualityMilestoneDefinitions()[0].GetName())
	assert.Equal(t, expectedMetadataKeys, qmdList.GetQualityMilestoneDefinitions()[0].GetExpectedMetadataKeys())

	// 4. Approve Release for QualityMilestone
	qualityMilestoneMetadata := []*v1.QualityMilestoneMetadata{
		{Key: "Image", Value: ""},
		{Key: "BuildURL", Value: ""},
	}
	qm, err := releaseClient.Approve(ctx, &v1.ReleaseServiceApproveRequest{
		Tag:                            expectedRelease.Tag,
		QualityMilestoneDefinitionName: qualityMilestoneDefinitionName,
		Approver:                       user,
		Metadata:                       qualityMilestoneMetadata,
	})
	assert.NoError(t, err)
	assert.Equal(t, user, qm.GetApprover())
	assert.Equal(t, qualityMilestoneMetadata[0].Key, qm.GetMetadata()[0].GetKey())
	assert.Equal(t, qualityMilestoneMetadata[0].Value, qm.GetMetadata()[0].GetValue())
	assert.Equal(t, qualityMilestoneMetadata[1].Key, qm.GetMetadata()[1].GetKey())
	assert.Equal(t, qualityMilestoneMetadata[1].Value, qm.GetMetadata()[1].GetValue())
	assert.Equal(t, qualityMilestoneDefinitionName, qm.GetQualityMilestoneDefinitionName())
	assert.Equal(t, expectedRelease.Tag, qm.GetTag())

	// 5. List all Releases with prefix at QualityMilestone
	releaseListAtQualityMilestone, err := releaseClient.List(ctx, &v1.ReleaseServiceListRequest{
		Prefix:               &prefix,
		Preload:              true,
		QualityMilestoneName: &qualityMilestoneDefinitionName,
	})
	assert.NoError(t, err)
	assert.Len(t, releaseList.GetReleases(), 1)
	actualReleaseFromListAtQualityMilestone := releaseListAtQualityMilestone.GetReleases()[0]

	utils.AssertReleasesAreEqual(
		t,
		actualListedRelease,
		conversions.NewReleaseFromGetReleaseResponse(actualReleaseFromListAtQualityMilestone),
		true,
		false,
	)
}

// Use case: Cloud Service Upgrader finds the latest version.
// Steps:
// 1. List all Releases at QualityMilestone.
// 2. FindLatest Release at QualityMilestone with prefix.
// 3. Approve another QualityMilestone.
// 4. Get Release including approved QualityMilestones.
//
//nolint:funlen
func TestFindingLatestRelease(t *testing.T) {
	qualityMilestoneDefinitionName := "Nightly has passed"
	prefix := "3.74"
	user := "lastname@redhat.com"

	basePath, err := e2e.GetFixturesPath()
	assert.NoError(t, err)

	dbPath := fmt.Sprintf("%s/%s", basePath, "prepopulated-with-approved-releases.sqlite")
	e2e.SetupE2ETest(t, dbPath)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)

	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(e2e.BufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)

	//nolint:errcheck
	defer conn.Close()

	// 1. List all Releases with Prefix at QualityMilestone
	releaseClient := v1.NewReleaseServiceClient(conn)
	releaseList, err := releaseClient.List(ctx, &v1.ReleaseServiceListRequest{
		QualityMilestoneName: &qualityMilestoneDefinitionName,
		Preload:              true,
	})
	assert.NoError(t, err)
	assert.Len(t, releaseList.GetReleases(), 2)

	// 2. FindLatest Release at QualityMilestone with prefix
	latestResponse, err := releaseClient.FindLatest(ctx, &v1.ReleaseServiceFindLatestRequest{
		Prefix:               &prefix,
		QualityMilestoneName: &qualityMilestoneDefinitionName,
		Preload:              true,
	})
	assert.NoError(t, err)
	actualPrefixedRelease, err := models.GetRelease("3.74.x-nightly-20230323", true, false)
	assert.NoError(t, err)
	utils.AssertReleasesAreEqual(
		t,
		actualPrefixedRelease,
		conversions.NewReleaseFromFindLatestReponse(latestResponse),
		true,
		true,
	)

	// 3. Approve another QualityMilestone
	qm, err := releaseClient.Approve(ctx, &v1.ReleaseServiceApproveRequest{
		Tag:                            latestResponse.GetTag(),
		QualityMilestoneDefinitionName: "Canary successful",
		Approver:                       user,
		Metadata: []*v1.QualityMilestoneMetadata{
			{Key: "DeploymentURL", Value: "this is a url"},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, user, qm.GetApprover())
	assert.Equal(t, "Canary successful", qm.GetQualityMilestoneDefinitionName())

	// 4. Get Release including approved QualityMilestones
	actualRelease, err := releaseClient.Get(ctx, &v1.ReleaseServiceGetRequest{
		Tag:     "3.74.x-nightly-20230323",
		Preload: true,
	})
	assert.NoError(t, err)
	assert.Len(t, actualRelease.GetQualityMilestones(), 2)
	assert.Equal(t, actualRelease.GetQualityMilestones()[0].GetName(), qualityMilestoneDefinitionName)
	assert.Equal(t, actualRelease.GetQualityMilestones()[1].GetName(), "Canary successful")
}
