package healthz

import (
	"net/http"

	"github.com/stackrox/release-registry/pkg/storage/models"
)

func healthLivenessHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("deep") == "true" {
		if _, err := models.ListQualityMilestoneDefinitions(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)

			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
