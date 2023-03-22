package qualitymilestonedefinition

import (
	"context"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/storage/models"
)

func (*server) List(
	context.Context, *v1.QualityMilestoneDefinitionServiceListRequest,
) (*v1.QualityMilestoneDefinitionServiceListResponse, error) {
	knownQmds, err := models.ListQualityMilestoneDefinitions()
	if err != nil {
		message := "could not list all QualityMilestoneDefinitions"
		log.Infow(message, "error", err.Error())

		return nil, errors.Wrap(err, message)
	}

	qmdListResponse := &v1.QualityMilestoneDefinitionServiceListResponse{}
	for i := range knownQmds {
		qmdListResponse.QualityMilestoneDefinition = append(
			qmdListResponse.QualityMilestoneDefinition,
			newGetResponseFromQualityMilestoneDefinition(&knownQmds[i]),
		)
	}

	return qmdListResponse, nil
}
