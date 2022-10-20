package system

import (
	"math"
	"strconv"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/rs/zerolog/log"

	"github.com/project-midgard/midgarts/internal/character"
	"github.com/project-midgard/midgarts/internal/character/actionindex"
	"github.com/project-midgard/midgarts/internal/character/actionplaymode"
	"github.com/project-midgard/midgarts/internal/character/directiontype"
	"github.com/project-midgard/midgarts/internal/component"
	"github.com/project-midgard/midgarts/internal/entity"
	"github.com/project-midgard/midgarts/internal/fileformat/act"
	"github.com/project-midgard/midgarts/internal/fileformat/grf"
	"github.com/project-midgard/midgarts/internal/fileformat/spr"
	"github.com/project-midgard/midgarts/internal/graphic"
	"github.com/project-midgard/midgarts/internal/graphic/geometry"
	"github.com/project-midgard/midgarts/internal/system/rendercmd"
)

const (
	SpriteScaleFactor    = float32(1.0)
	FixedCameraDirection = 6
)

type CharacterRenderable interface {
	ecs.BasicFace
	component.CharacterStateComponentFace
	component.CharacterAttachmentComponentFace
}

type CharacterRenderSystem struct {
	grfFile         *grf.File
	characters      map[string]*entity.Character
	RenderCommands  *RenderCommands
	textureProvider graphic.TextureProvider
}

func NewCharacterRenderSystem(grfFile *grf.File, textureProvider graphic.TextureProvider) *CharacterRenderSystem {
	return &CharacterRenderSystem{
		grfFile:    grfFile,
		characters: map[string]*entity.Character{},
		RenderCommands: &RenderCommands{
			sprite: []rendercmd.SpriteRenderCommand{},
		},
		textureProvider: textureProvider,
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
	cmp, err := component.NewCharacterAttachmentComponent(s.grfFile, component.CharacterAttachmentComponentConfig{
		Gender:           char.Gender,
		JobSpriteID:      char.JobSpriteID,
		HeadIndex:        char.HeadIndex,
		EnableShield:     char.HasShield,
		ShieldSpriteName: char.ShieldSpriteName,
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	char.SetCharacterAttachmentComponent(cmp)
	s.characters[strconv.Itoa(int(char.ID()))] = char
}

func (s *CharacterRenderSystem) Remove(e ecs.BasicEntity) {
	delete(s.characters, strconv.Itoa(int(e.ID())))
}

func (s *CharacterRenderSystem) renderCharacter(dt float32, char *entity.Character) {
	offset := [2]float32{0, 0}

	direction := int(char.Direction) + directiontype.DirectionTable[FixedCameraDirection]%8
	behind := direction > 1 && direction < 6
	renderShield := char.HasShield && char.ActionIndex == actionindex.StandBy

	if char.ActionIndex != actionindex.Dead && char.ActionIndex != actionindex.Sitting {
		s.renderAttachment(dt, char, character.AttachmentShadow, &offset)
	}

	if behind && renderShield {
		s.renderAttachment(dt, char, character.AttachmentShield, &offset)
	}

	s.renderAttachment(dt, char, character.AttachmentBody, &offset)
	s.renderAttachment(dt, char, character.AttachmentHead, &offset)

	if !behind && renderShield {
		s.renderAttachment(dt, char, character.AttachmentShield, &offset)
	}
}

func (s *CharacterRenderSystem) renderAttachment(
	dt float32,
	char *entity.Character,
	elem character.AttachmentType,
	offset *[2]float32,
) {
	var actions []*act.Action
	if actions = char.Files[elem].ACT.Actions; len(actions) == 0 {
		return
	}

	idx := (int(char.ActionIndex) + (int(char.Direction)+directiontype.DirectionTable[FixedCameraDirection])%8) % len(actions)
	action := actions[idx]
	frameCount := int64(len(action.Frames))
	timeNeededForOneFrame := int64(float64(action.Delay) * (1.0 / char.FPSMultiplier))

	if char.ForcedDuration != 0 {
		timeNeededForOneFrame = int64(char.ForcedDuration) / frameCount
	}

	timeNeededForOneFrame = int64(math.Max(float64(timeNeededForOneFrame), 100))
	elapsedTime := time.Since(char.AnimationStartedAt).Milliseconds() - int64(dt)
	realIndex := elapsedTime / timeNeededForOneFrame

	var frameIndex int64
	switch char.PlayMode {
	case actionplaymode.Repeat:
		frameIndex = realIndex % frameCount
		break
	}

	// Ignore "doridori" animation
	if len(action.Frames) == 3 {
		frameIndex = 0
	}

	var frame *act.ActionFrame
	if frame = action.Frames[frameIndex]; len(frame.Layers) == 0 {
		*offset = [2]float32{0, 0}
		return
	}

	position := [2]float32{0, 0}

	if len(frame.Positions) > 0 &&
		elem != character.AttachmentBody &&
		elem != character.AttachmentShield {

		position[0] = offset[0] - float32(frame.Positions[0][0])
		position[1] = offset[1] - float32(frame.Positions[0][1])
	}

	// Render all layers
	for _, layer := range frame.Layers {
		if layer.SpriteFrameIndex < 0 {
			continue
		}

		s.renderLayer(char, layer, char.Files[elem].SPR, position)
	}

	// Save offset reference
	if len(frame.Positions) > 0 {
		*offset = [2]float32{
			float32(frame.Positions[0][0]),
			float32(frame.Positions[0][1]),
		}
	}

	char.AnimationDelay = time.Duration(action.DurationMilliseconds) * time.Millisecond
}

func (s *CharacterRenderSystem) renderLayer(
	char *entity.Character,
	layer *act.ActionFrameLayer,
	spr *spr.SpriteFile,
	offset [2]float32,
) {
	frameIndex := character.SpriteIndex(layer.SpriteFrameIndex)
	if frameIndex < 0 {
		return
	}

	texture, err := s.textureProvider.NewTextureFromRGBA(spr.ImageAt(frameIndex))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	frame := spr.Frames[frameIndex]
	width, height := float32(frame.Width), float32(frame.Height)
	width *= layer.Scale[0] * SpriteScaleFactor * geometry.OnePixelSize
	height *= layer.Scale[1] * SpriteScaleFactor * geometry.OnePixelSize
	rot := float64(layer.Angle) * (math.Pi / 180)

	offset = [2]float32{
		(float32(layer.Position[0]) + offset[0]) * geometry.OnePixelSize,
		(float32(layer.Position[1]) + offset[1]) * geometry.OnePixelSize,
	}

	cmd := rendercmd.SpriteRenderCommand{
		Scale:           layer.Scale,
		Size:            mgl32.Vec2{width, height},
		Position:        char.Position(),
		Offset:          mgl32.Vec2{offset[0], offset[1]},
		RotationRadians: float32(rot),
		Texture:         texture,
		FlipVertically:  layer.Mirrored,
	}

	// This is the current API to render a shader. Commands will
	// be collected by the lower-level rendering system (OpenGL).
	s.renderSpriteCommand(cmd)
}

func (s *CharacterRenderSystem) renderSpriteCommand(cmd ...rendercmd.SpriteRenderCommand) {
	s.RenderCommands.sprite = append(s.RenderCommands.sprite, cmd...)
}
