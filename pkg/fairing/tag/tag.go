/*
Copyright 2018 The Kubeflow Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tag

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/kubeflow/fairing/pkg/fairing/gcp"
	"github.com/pkg/errors"
)

func ChecksumTag(path string) (string, error) {
	projectID, err := gcp.ProjectID()
	if err != nil {
		return "", errors.Wrap(err, "getting project id")
	}
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "opening file for checksum")
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", errors.Wrap(err, "checksumming")
	}

	return fmt.Sprintf("gcr.io/%s/fairing-job:%x", projectID, h.Sum(nil)), nil
}
