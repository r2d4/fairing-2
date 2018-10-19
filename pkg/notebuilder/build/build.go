package build

import (
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"

	"github.com/pkg/errors"
)

// Append appends a tarball to a base image and uploads the new image
// base is the base image to append to
// tag is the new tag that will be uploaded
// tarPath is the path to a tarball that will be appended as a layer
func Append(base, tag, tarPath string) error {
	baseRef, err := name.ParseReference(base, name.WeakValidation)
	if err != nil {
		return errors.Wrap(err, "parsing source tag")
	}

	baseImage, err := remote.Image(baseRef, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return errors.Wrap(err, "getting source image ref")
	}

	l, err := tarball.LayerFromFile(tarPath)
	if err != nil {
		return errors.Wrapf(err, "generating layer from tarball %s", tarPath)
	}

	image, err := mutate.AppendLayers(baseImage, l)
	if err != nil {
		return errors.Wrap(err, "appending layer")
	}

	if err := UploadImage(image, tag); err != nil {
		return errors.Wrapf(err, "uploading image %s", tag)
	}

	return nil
}

// UploadImage uploads an image to a remote registry, using the default keychain
// for that registry.
func UploadImage(image v1.Image, tag string) error {
	dstTag, err := name.NewTag(tag, name.WeakValidation)
	if err != nil {
		return errors.Wrap(err, "parsing dst tag")
	}

	dstAuth, err := authn.DefaultKeychain.Resolve(dstTag.Context().Registry)
	if err != nil {
		return errors.Wrap(err, "getting credentials")
	}

	if err := remote.Write(dstTag, image, dstAuth, http.DefaultTransport); err != nil {
		return errors.Wrap(err, "uploading image")
	}
	return nil
}
