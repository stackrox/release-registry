// Package conversions converts business objects from storage to service representation and back.
package conversions

import (
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/version"
	"gorm.io/gorm"
)

// NewCreateReleaseRequestFromRelease converts new Releases from storage to storage representation (request).
func NewCreateReleaseRequestFromRelease(release *models.Release) *v1.ReleaseServiceCreateRequest {
	return &v1.ReleaseServiceCreateRequest{
		Tag:      release.Tag,
		Commit:   release.Commit,
		Metadata: newV1ReleaseMetadataFromRelease(release),
	}
}

// NewCreateReleaseResponseFromRelease converts new Releases from storage to storage representation (response).
func NewCreateReleaseResponseFromRelease(release *models.Release) *v1.ReleaseServiceCreateResponse {
	return &v1.ReleaseServiceCreateResponse{
		Meta:     newV1MetaFromRelease(release),
		Tag:      release.Tag,
		Commit:   release.Commit,
		Creator:  release.Creator,
		Kind:     v1.ReleaseKind(release.Kind),
		Metadata: newV1ReleaseMetadataFromRelease(release),
	}
}

// NewReleaseFromCreateReleaseResponse converts new Releases from service to storage representation (response).
func NewReleaseFromCreateReleaseResponse(release *v1.ReleaseServiceCreateResponse) *models.Release {
	return &models.Release{
		Model: gorm.Model{
			ID:        uint(release.GetMeta().Id),
			CreatedAt: release.GetMeta().GetCreatedAt().AsTime(),
			UpdatedAt: release.GetMeta().GetCreatedAt().AsTime(),
		},
		Tag:      release.GetTag(),
		Commit:   release.GetCommit(),
		Creator:  release.GetCreator(),
		Kind:     version.Kind(release.GetKind()),
		Metadata: newReleaseMetadataFromV1ReleaseMetadata(release.GetMetadata()),
	}
}

// NewGetReleaseResponseFromRelease converts a Release from storage to service representation.
func NewGetReleaseResponseFromRelease(release *models.Release) *v1.ReleaseServiceGetResponse {
	return &v1.ReleaseServiceGetResponse{
		Meta:              newV1MetaFromRelease(release),
		Tag:               release.Tag,
		Commit:            release.Commit,
		Creator:           release.Creator,
		Kind:              v1.ReleaseKind(release.Kind),
		Metadata:          newV1ReleaseMetadataFromRelease(release),
		Rejected:          release.Rejected,
		QualityMilestones: newV1QualityMilestonesFromQualityMilestones(release.QualityMilestones),
	}
}

// NewUpdateReleaseResponseFromRelease converts a release from storage to service representation.
func NewUpdateReleaseResponseFromRelease(release *models.Release) *v1.ReleaseServiceUpdateResponse {
	return &v1.ReleaseServiceUpdateResponse{
		Meta:              newV1MetaFromRelease(release),
		Tag:               release.Tag,
		Commit:            release.Commit,
		Creator:           release.Creator,
		Kind:              v1.ReleaseKind(release.Kind),
		Metadata:          newV1ReleaseMetadataFromRelease(release),
		Rejected:          release.Rejected,
		QualityMilestones: newV1QualityMilestonesFromQualityMilestones(release.QualityMilestones),
	}
}

// NewUpdateReleaseRequestFromRelease converts a Release from storage to service representation.
func NewUpdateReleaseRequestFromRelease(release *models.Release, includeRejected bool) *v1.ReleaseServiceUpdateRequest {
	return &v1.ReleaseServiceUpdateRequest{
		Tag:             release.Tag,
		Metadata:        newV1ReleaseMetadataFromRelease(release),
		IncludeRejected: includeRejected,
	}
}

// NewReleaseFromUpdateReleaseResponse converts a release from service to storage representation.
func NewReleaseFromUpdateReleaseResponse(release *v1.ReleaseServiceUpdateResponse) *models.Release {
	return &models.Release{
		Model: gorm.Model{
			ID:        uint(release.GetMeta().Id),
			CreatedAt: release.GetMeta().GetCreatedAt().AsTime(),
			UpdatedAt: release.GetMeta().GetCreatedAt().AsTime(),
		},
		Tag:               release.GetTag(),
		Commit:            release.GetCommit(),
		Creator:           release.GetCreator(),
		Kind:              version.Kind(release.GetKind()),
		Metadata:          newReleaseMetadataFromV1ReleaseMetadata(release.GetMetadata()),
		QualityMilestones: newQualityMilestonesFromV1QualityMilestones(release.GetQualityMilestones()),
	}
}

// NewReleaseFromGetReleaseResponse converts a Release from service to storage representation.
func NewReleaseFromGetReleaseResponse(release *v1.ReleaseServiceGetResponse) *models.Release {
	return &models.Release{
		Model: gorm.Model{
			ID:        uint(release.GetMeta().Id),
			CreatedAt: release.GetMeta().GetCreatedAt().AsTime(),
			UpdatedAt: release.GetMeta().GetCreatedAt().AsTime(),
		},
		Tag:               release.GetTag(),
		Commit:            release.GetCommit(),
		Creator:           release.GetCreator(),
		Kind:              version.Kind(release.GetKind()),
		Metadata:          newReleaseMetadataFromV1ReleaseMetadata(release.GetMetadata()),
		QualityMilestones: newQualityMilestonesFromV1QualityMilestones(release.GetQualityMilestones()),
	}
}

// NewListReleaseResponseFromReleases converts a list of Releases from storage to service representation.
func NewListReleaseResponseFromReleases(releases []models.Release) *v1.ReleaseServiceListResponse {
	releaseResponses := &v1.ReleaseServiceListResponse{}

	for i := range releases {
		release := releases[i]
		releaseResponses.Releases = append(releaseResponses.Releases, NewGetReleaseResponseFromRelease(&release))
	}

	return releaseResponses
}

// NewFindLatestReleaseResponseFromRelease converts a latest Release from storage to service representation.
func NewFindLatestReleaseResponseFromRelease(release *models.Release) *v1.ReleaseServiceFindLatestResponse {
	return &v1.ReleaseServiceFindLatestResponse{
		Meta:              newV1MetaFromRelease(release),
		Tag:               release.Tag,
		Commit:            release.Commit,
		Creator:           release.Creator,
		Kind:              v1.ReleaseKind(release.Kind),
		Metadata:          newV1ReleaseMetadataFromRelease(release),
		Rejected:          release.Rejected,
		QualityMilestones: newV1QualityMilestonesFromQualityMilestones(release.QualityMilestones),
	}
}

// NewReleaseFromFindLatestReponse converts a latest Release from service to storage representation.
func NewReleaseFromFindLatestReponse(release *v1.ReleaseServiceFindLatestResponse) *models.Release {
	return &models.Release{
		Model: gorm.Model{
			ID:        uint(release.GetMeta().Id),
			CreatedAt: release.GetMeta().GetCreatedAt().AsTime(),
			UpdatedAt: release.GetMeta().GetCreatedAt().AsTime(),
		},
		Tag:               release.GetTag(),
		Commit:            release.GetCommit(),
		Creator:           release.GetCreator(),
		Kind:              version.Kind(release.GetKind()),
		Metadata:          newReleaseMetadataFromV1ReleaseMetadata(release.GetMetadata()),
		QualityMilestones: newQualityMilestonesFromV1QualityMilestones(release.GetQualityMilestones()),
	}
}

// NewApproveQualityMilestoneResponseFromQualityMilestone converts a QualityMilestone (storage representation)
// to an approve release (service representation).
func NewApproveQualityMilestoneResponseFromQualityMilestone(
	qm *models.QualityMilestone,
) *v1.ReleaseServiceApproveResponse {
	return &v1.ReleaseServiceApproveResponse{
		Meta:                           newV1MetaFromQualityMilestone(qm),
		Tag:                            qm.Release.Tag,
		QualityMilestoneDefinitionName: qm.QualityMilestoneDefinition.Name,
		Approver:                       qm.Approver,
		Kind:                           v1.ReleaseKind(qm.Release.Kind),
		Metadata:                       newV1QualityMilestoneMetadataFromQualityMilestoneMetadata(qm.Metadata),
	}
}

// NewRejectReleaseResponseFromRelease converts a rejected release from storage to service representation.
func NewRejectReleaseResponseFromRelease(release *models.Release) *v1.ReleaseServiceRejectResponse {
	return &v1.ReleaseServiceRejectResponse{
		Meta:     newV1MetaFromRelease(release),
		Tag:      release.Tag,
		Commit:   release.Commit,
		Creator:  release.Creator,
		Kind:     v1.ReleaseKind(release.Kind),
		Metadata: newV1ReleaseMetadataFromRelease(release),
		Rejected: release.Rejected,
	}
}

func newV1ReleaseMetadataFromRelease(release *models.Release) []*v1.ReleaseMetadata {
	releaseMetadata := []*v1.ReleaseMetadata{}

	for _, metadata := range release.Metadata {
		releaseMetadata = append(releaseMetadata, &v1.ReleaseMetadata{
			Key:   metadata.Key,
			Value: metadata.Value,
		})
	}

	return releaseMetadata
}

func newReleaseMetadataFromV1ReleaseMetadata(v1ReleaseMetadata []*v1.ReleaseMetadata) []models.ReleaseMetadata {
	releaseMetadata := []models.ReleaseMetadata{}

	for _, metadata := range v1ReleaseMetadata {
		releaseMetadata = append(releaseMetadata, models.ReleaseMetadata{
			Key: metadata.GetKey(), Value: metadata.GetValue(),
		})
	}

	return releaseMetadata
}
