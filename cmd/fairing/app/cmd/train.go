/*
Copyright 2018 COMPANY

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

package cmd

import (
	"fmt"
	"io"

	"github.com/kubeflow/fairing/pkg/fairing/train"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	trainingImage, dstImage, srcTar string
)

func NewCmdTrain(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "train",
		Short: "trains an image from a notebook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if trainingImage == "" || srcTar == "" {
				return fmt.Errorf("training-image and source are all required flags")
			}
			return RunTrain(out, cmd)
		},
	}

	cmd.Flags().StringVar(&trainingImage, "training-image", "", "the base training image to append to")
	cmd.Flags().StringVar(&dstImage, "tag", "", "the tag to push to, if not provided one will be generated")
	cmd.Flags().StringVar(&srcTar, "source", "", "a tar.gz file to append as a layer")
	return cmd
}

func RunTrain(out io.Writer, cmd *cobra.Command) error {
	if err := train.Train(trainingImage, dstImage, srcTar); err != nil {
		return errors.Wrap(err, "executing template")
	}
	return nil
}
