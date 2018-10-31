package train

import (
	"fmt"

	"github.com/kubeflow/fairing/pkg/fairing/build"
	"github.com/kubeflow/fairing/pkg/fairing/gcp"
	"github.com/pkg/errors"
)

func Train(trainingImg, dstTag, srcTar string) error {
	if dstTag == "" {
		projectID, err := gcp.ProjectID()
		if err != nil {
			return errors.Wrap(err, "getting project id")
		}

		dstTag = fmt.Sprintf("gcr.io/%s/kubeflow-training-model:%s", projectID, srcTar)
	}

	if err := build.Build(trainingImg, dstTag, srcTar); err != nil {
		return errors.Wrap(err, "building image")
	}

	return nil
}
