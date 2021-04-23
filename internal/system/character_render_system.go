package system

import (
	"log"
	"math"
	"strconv"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/project-midgard/midgarts/internal/component"
	"github.com/project-midgard/midgarts/internal/entity"
	"github.com/project-midgard/midgarts/internal/system/rendercmd"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
	"github.com/project-midgard/midgarts/pkg/graphic"
)

const (
	SpriteScaleFactor    = float32(1.0)
	FPSMultiplier        = float64(1.0)
	FixedCameraDirection = 6
)

type CharacterRenderable interface {
	common.BasicFace
	component.CharacterStateComponentFace
	component.CharacterAttachmentComponentFace
}

type CharacterRenderSystem struct {
	grfFile        *grf.File
	characters     map[string]*entity.Character
	RenderCommands *RenderCommands
}

func NewCharacterRenderSystem(grfFile *grf.File) *CharacterRenderSystem {
	return &CharacterRenderSystem{
		grfFile:    grfFile,
		characters: map[string]*entity.Character{},
		RenderCommands: &RenderCommands{
			sprite: []rendercmd.SpriteRenderCommand{},
		},
	}
}

func (s *CharacterRenderSystem) Update(dt float32) {
	s.RenderCommands.sprite = []rendercmd.SpriteRenderCommand{}

	for _, char := range s.characters {
		s.renderCharacter(dt, char)
	}
}

func (s *CharacterRenderSystem) AddByInterface(o ecs.Identifier) {
	char := o.(*entity.Character)
	s.Add(char)
}

func (s *CharacterRenderSystem) Add(char *entity.Character) {
	cmp, e := component.NewCharacterAttachmentComponent(s.grfFile, char.Gender, char.JobSpriteID, char.HeadIndex)
	if e != nil {
		log.Fatal(e)
	}

	char.SetCharacterAttachmentComponent(cmp)
	s.characters[strconv.Itoa(int(char.ID()))] = char
}

func (s *CharacterRenderSystem) Remove(e ecs.BasicEntity) {
	delete(s.characters, strconv.Itoa(int(e.ID())))
}

func (s *CharacterRenderSystem) renderCharacter(dt float32, char *entity.Character) {
	offset := [2]float32{}

	if char.ActionIndex != actionindex.Dead && char.ActionIndex != actionindex.Sitting {
		s.renderAttachment(dt, char, character.AttachmentShadow, &offset)
	}

	s.renderAttachment(dt, char, character.AttachmentBody, &offset)
	s.renderAttachment(dt, char, character.AttachmentHead, &offset)
}

func (s *CharacterRenderSystem) renderAttachment(
	dt float32,
	char *entity.Character,
	elem character.AttachmentType,
	offset *[2]float32,
) {
	if len(char.Files[elem].ACT.Actions) == 0 {
		return
	}

	var idx int
	if elem == character.AttachmentShadow {
		idx = 0
	} else {
		idx = int(char.ActionIndex) + (int(char.Direction)+directiontype.DirectionTable[FixedCameraDirection])%8
	}

	action := char.Files[elem].ACT.Actions[idx]
	fileSet := char.Files[elem]

	frameCount := len(action.Frames)
	timeNeededForOneFrame := int64(float64(action.Delay) * (1.0 / char.FPSMultiplier))
	if char.ForcedDuration != 0 {
		timeNeededForOneFrame = int64(char.ForcedDuration) / int64(frameCount)
	}
	timeNeededForOneFrame = int64(math.Max(float64(timeNeededForOneFrame), 100))
	elapsedTime := time.Since(char.AnimationStartedAt).Milliseconds() - int64(dt)
	realIndex := elapsedTime / timeNeededForOneFrame

	var frameIndex int64
	switch char.PlayMode {
	case actionplaymode.Repeat:
		frameIndex = realIndex % int64(frameCount)
		break
	}

	frame := action.Frames[frameIndex]

	if len(frame.Layers) == 0 {
		return
	}

	position := [2]float32{0, 0}

	if len(frame.Positions) > 0 && elem != character.AttachmentBody {
		position[0] = offset[0] - float32(frame.Positions[frameIndex][0])
		position[1] = offset[1] - float32(frame.Positions[frameIndex][1])
	}

	// Render all frames
	for _, layer := range frame.Layers {
		if layer.SpriteFrameIndex < 0 {
			continue
		}

		s.renderLayer(char, frameIndex, layer, fileSet.SPR, position)
	}

	// Save offset reference
	if len(frame.Positions) > 0 {
		*offset = [2]float32{
			float32(frame.Positions[frameIndex][0]),
			float32(frame.Positions[frameIndex][1]),
		}
	}
}

func (s *CharacterRenderSystem) renderLayer(
	char *entity.Character,
	frameIndex int64,
	layer *act.ActionFrameLayer,
	spr *spr.SpriteFile,
	prevOffset [2]float32,
) {
	frameIndex = int64(layer.SpriteFrameIndex)
	if frameIndex < 0 {
		return
	}

	texture, err := graphic.NewTextureFromImage(spr.ImageAt(frameIndex))
	if err != nil {
		log.Fatal(err)
	}

	frame := spr.Frames[frameIndex]
	width, height := float32(frame.Width), float32(frame.Height)
	width *= layer.Scale[0] * SpriteScaleFactor * graphic.OnePixelSize
	height *= layer.Scale[1] * SpriteScaleFactor * graphic.OnePixelSize

	offset := [2]float32{
		(float32(layer.Position[0]) + prevOffset[0]) * graphic.OnePixelSize,
		(float32(layer.Position[1]) + prevOffset[1]) * graphic.OnePixelSize,
	}

	// This is the current API to render a sprite. Commands will
	// be collected by the lower-level rendering system (OpenGL).
	s.renderSpriteCommand(rendercmd.SpriteRenderCommand{
		Scale:           layer.Scale,
		Size:            mgl32.Vec2{width, height},
		Position:        mgl32.Vec3{char.Position().X(), char.Position().Y(), char.Position().Z()},
		Offset:          mgl32.Vec2{offset[0], offset[1]},
		RotationRadians: 0,
		Texture:         texture,
		FlipVertically:  layer.Mirrored,
	})
}

func (s *CharacterRenderSystem) renderSpriteCommand(cmd ...rendercmd.SpriteRenderCommand) {
	s.RenderCommands.sprite = append(s.RenderCommands.sprite, cmd...)
}
