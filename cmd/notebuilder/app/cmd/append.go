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
	"io"

	"github.com/r2d4/notebuilder/pkg/notebuilder/build"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	baseImage, dstImage, layerFile string
)

func NewCmdAppend(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "append",
		Short: "Appends a tarball to an image and up",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunAppend(out, cmd)
		},
	}

	cmd.Flags().StringVar(&baseImage, "base-image", "", "the base image to append to")
	cmd.Flags().StringVar(&dstImage, "dst-image", "", "the image tag to push")
	cmd.Flags().StringVar(&dstImage, "layer-file", "", "a tar.gz file to append as a layer")

	cmd.Flags().VarP(versionFlag, "output", "o", versionFlag.Usage())
	return cmd
}

func RunAppend(out io.Writer, cmd *cobra.Command) error {
	if err := build.Append(baseImage, dstImage, layerFile); err != nil {
		return errors.Wrap(err, "executing template")
	}
	return nil
}
