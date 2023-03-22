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

func (*server) Get(
	ctx context.Context, qmdReq *v1.QualityMilestoneDefinitionServiceGetRequest,
) (*v1.QualityMilestoneDefinitionServiceGetResponse, error) {
	name := qmdReq.GetName()
	qmd, err := models.GetQualityMilestoneDefinition(name)

	if err != nil {
		message := "could not find requested QualityMilestoneDefinition"
		log.Infow(message, "name", name, "error", err.Error())

		return nil, errors.Wrap(err, fmt.Sprintf("%s: %s", message, name))
	}

	qmdResponse := newGetResponseFromQualityMilestoneDefinition(qmd)

	return qmdResponse, nil
}

func newGetResponseFromQualityMilestoneDefinition(
	qmd *models.QualityMilestoneDefinition,
) *v1.QualityMilestoneDefinitionServiceGetResponse {
	return &v1.QualityMilestoneDefinitionServiceGetResponse{
		Meta: &v1Meta.Meta{
			Id:        int64(qmd.ID),
			CreatedAt: timestamppb.New(qmd.CreatedAt),
			UpdatedAt: timestamppb.New(qmd.UpdatedAt),
		},
		Name:                 qmd.Name,
		ExpectedMetadataKeys: qmd.ExpectedMetadataKeys,
	}
}
