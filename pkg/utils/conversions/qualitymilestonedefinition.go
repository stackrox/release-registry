package conversions

import (
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	v1Meta "github.com/stackrox/release-registry/gen/go/proto/shared/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewCreateQualityMilestoneDefinitionResponseFromQualityMilestoneDefinition converts QualityMilestoneDefinition
// from storage to service representation.
func NewCreateQualityMilestoneDefinitionResponseFromQualityMilestoneDefinition(
	qmd *models.QualityMilestoneDefinition,
) *v1.QualityMilestoneDefinitionServiceCreateResponse {
	return &v1.QualityMilestoneDefinitionServiceCreateResponse{
		Meta: &v1Meta.Meta{
			Id:        int64(qmd.ID),
			CreatedAt: timestamppb.New(qmd.CreatedAt),
			UpdatedAt: timestamppb.New(qmd.UpdatedAt),
		},
		Name:                 qmd.Name,
		ExpectedMetadataKeys: qmd.ExpectedMetadataKeys,
	}
}

// NewGetQualityMilestoneDefinitionResponseFromQualityMilestoneDefinition converts QualityMilestoneDefinition
// from storage to service representation.
func NewGetQualityMilestoneDefinitionResponseFromQualityMilestoneDefinition(
	qmd *models.QualityMilestoneDefinition,
) *v1.QualityMilestoneDefinitionServiceGetResponse {
	return &v1.QualityMilestoneDefinitionServiceGetResponse{
		Meta:                 newV1MetaFromQualityMilestoneDefinition(qmd),
		Name:                 qmd.Name,
		ExpectedMetadataKeys: qmd.ExpectedMetadataKeys,
	}
}
