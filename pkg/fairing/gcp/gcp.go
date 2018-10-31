package gcp

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/kubeflow/fairing/pkg/fairing/constants"
)

type ServiceAccount struct {
	ProjectID string `json:"project_id"`
}

func ProjectID() (string, error) {
	secretPath := os.Getenv(constants.GoogleCredentialsEnv)
	if secretPath == "" {
		return "", fmt.Errorf("could not get credentials file path from env")
	}

	f, err := os.Open(secretPath)
	if err != nil {
		return "", errors.Wrap(err, "getting secret file")
	}
	defer f.Close()
	var sa *ServiceAccount
	d := json.NewDecoder(f)
	if err := d.Decode(&sa); err != nil {
		return "", errors.Wrap(err, "decoding credentials file")
	}

	return sa.ProjectID, nil
}
