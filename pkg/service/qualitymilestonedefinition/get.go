package qualitymilestonedefinition

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
	"github.com/stackrox/release-registry/pkg/utils/conversions"
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

	qmdResponse := conversions.NewGetQualityMilestoneDefinitionResponseFromQualityMilestoneDefinition(qmd)

	return qmdResponse, nil
}
