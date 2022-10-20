package caching

import (
	"github.com/google/uuid"

	graphic2 "github.com/project-midgard/midgarts/internal/graphic"
)

type CachedTextureProvider map[uuid.UUID]*graphic2.Texture

func NewCachedTextureProvider() CachedTextureProvider {
	return make(map[uuid.UUID]*graphic2.Texture)
}

func (t CachedTextureProvider) NewTextureFromRGBA(rgba *graphic2.UniqueRGBA) (*graphic2.Texture, error) {
	if txt, ok := t[rgba.ID]; ok {
		return txt, nil
	}

	tex, err := graphic2.NewTextureFromRGBA(rgba)
	if err != nil {
		return nil, err
	}

	t[rgba.ID] = tex
	return tex, nil
}
