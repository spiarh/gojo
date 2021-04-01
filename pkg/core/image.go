package core

import (
	"fmt"

	"github.com/spiarh/gojo/pkg/util"
)

func NewImageFromFQIN(fqin string) (Image, error) {
	var image = Image{}
	registry, name, tag, err := util.ParseFQIN(fqin)
	if err != nil {
		return image, err
	}
	image = NewImage(registry, name, tag)

	return image, nil
}

func NewImage(registry, name, tag string) Image {
	return Image{
		Registry: registry,
		Name:     name,
		Tag:      tag,
	}
}

func (b *Image) String() string {
	return fmt.Sprintf("%s/%s:%s", b.Registry, b.Name, b.Tag)
}

func (b *Image) StringWithTag(tag string) string {
	return fmt.Sprintf("%s/%s:%s", b.Registry, b.Name, tag)
}

func (b *Image) StringWithTagLatest() string {
	return fmt.Sprintf("%s/%s:%s", b.Registry, b.Name, "latest")
}
