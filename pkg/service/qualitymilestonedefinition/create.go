package qualitymilestonedefinition

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	v1Meta "github.com/stackrox/release-registry/gen/go/proto/shared/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (*server) Create(
	ctx context.Context, newQmdReq *v1.QualityMilestoneDefinitionServiceCreateRequest,
) (*v1.QualityMilestoneDefinitionServiceCreateResponse, error) {
	name := newQmdReq.GetName()
	qmd, err := models.CreateQualityMilestoneDefinition(name, newQmdReq.GetExpectedMetadataKeys())

	if err != nil {
		message := "could not create QualityMilestoneDefinition"
		log.Infow(message, "name", name, "error", err.Error())

		return nil, errors.Wrap(err, fmt.Sprintf("%s '%s'", message, name))
	}

	qmdResponse := newCreateResponseFromQualityMilestoneDefinition(qmd)

	return qmdResponse, nil
}

func newCreateResponseFromQualityMilestoneDefinition(
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
