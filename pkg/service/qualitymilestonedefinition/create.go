package qualitymilestonedefinition

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/conversions"
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

	qmdResponse := conversions.NewCreateQualityMilestoneDefinitionResponseFromQualityMilestoneDefinition(qmd)

	return qmdResponse, nil
}
