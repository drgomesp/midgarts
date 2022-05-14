package graphic

import (
	"github.com/google/uuid"
	"image"
)

type UniqueRGBA struct {
	ID uuid.UUID
	*image.RGBA
}

func NewUniqueRGBA(r image.Rectangle) *UniqueRGBA {
	return &UniqueRGBA{
		ID:   uuid.New(),
		RGBA: image.NewRGBA(r),
	}
}
