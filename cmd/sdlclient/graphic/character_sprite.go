package graphic

import (
	"fmt"
	"log"
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
	graphic "github.com/project-midgard/midgarts/pkg/graphic"
	"golang.org/x/text/encoding/charmap"
)

type CharacterSpriteElement int

const (
	SpriteScaleFactor    = float32(1.0)
	FPSMultiplier        = float64(1.0)
	FixedCameraDirection = 6
)

const (
	CharacterSpriteElementShadow = CharacterSpriteElement(iota)
	CharacterSpriteElementBody
	CharacterSpriteElementHead

	NumCharacterSpriteElements
)

func (e CharacterSpriteElement) String() string {
	switch e {
	case CharacterSpriteElementShadow:
		return "CharacterSpriteElementShadow"
	case CharacterSpriteElementBody:
		return "CharacterSpriteElementBody"
	case CharacterSpriteElementHead:
		return "CharacterSpriteElementHead"
	default:
		return "Unknown"
	}
}

var DirectionTable = [8]int{6, 5, 4, 3, 2, 1, 0, 7}

type CharState struct {
	Direction directiontype.Type
	State     statetype.Type
	PlayMode  actionplaymode.Type
}

type fileSet struct {
	ACT *act.ActionFile
	SPR *spr.SpriteFile
}

type CharacterSprite struct {
	*graphic.Transform

	Gender character.GenderType

	files   [NumCharacterSpriteElements]fileSet
	sprites [NumCharacterSpriteElements]*graphic.Sprite
}

func LoadCharacterSprite(f *grf.File, gender character.GenderType, jobSpriteID jobspriteid.Type, headIndex int32) (
	sprite *CharacterSprite,
	err error,
) {
	jobFileName := character.JobSpriteNameTable[jobSpriteID]
	if "" == jobFileName {
		return nil, fmt.Errorf("unsupported jobSpriteID %v", jobSpriteID)
	}

	decodedFolderA, err := getDecodedFolder([]byte{0xC0, 0xCE, 0xB0, 0xA3, 0xC1, 0xB7})
	if err != nil {
		return nil, err
	}

	decodedFolderB, err := getDecodedFolder([]byte{0xB8, 0xF6, 0xC5, 0xEB})
	if err != nil {
		return nil, err
	}

	var (
		bodyFilePath   string
		shadowFilePath = "data/sprite/shadow"
		headFilePathf  = "data/sprite/ÀÎ°£Á·/¸Ó¸®Åë/%s/%d_%s"
	)

	if character.Male == gender {
		bodyFilePath = fmt.Sprintf(character.MaleFilePathf, decodedFolderA, decodedFolderB, jobFileName)
		headFilePathf = fmt.Sprintf(headFilePathf, "³²", headIndex, "³²")
	} else {
		bodyFilePath = fmt.Sprintf(character.FemaleFilePathf, decodedFolderA, decodedFolderB, jobFileName)
		headFilePathf = fmt.Sprintf(headFilePathf, "¿©", headIndex, "¿©")
	}

	shadowActFile, shadowSprFile, err := f.GetActionAndSpriteFiles(shadowFilePath)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not load shadow act and spr files (%v, %s)", gender, jobSpriteID))
	}

	bodyActFile, bodySprFile, err := f.GetActionAndSpriteFiles(bodyFilePath)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not load body act and spr files (%v, %s)", gender, jobSpriteID))
	}

	headActFile, headSprFile, err := f.GetActionAndSpriteFiles(headFilePathf)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not load head act and spr files (%v, %s)", gender, jobSpriteID))
	}

	characterSprite := &CharacterSprite{
		Transform: graphic.NewTransform(graphic.Origin),
		Gender:    gender,
		files: [NumCharacterSpriteElements]fileSet{
			CharacterSpriteElementShadow: {shadowActFile, shadowSprFile},
			CharacterSpriteElementHead:   {headActFile, headSprFile},
			CharacterSpriteElementBody:   {bodyActFile, bodySprFile},
		},
		sprites: [NumCharacterSpriteElements]*graphic.Sprite{},
	}

	return characterSprite, nil
}

func getDecodedFolder(buf []byte) (folder []byte, err error) {
	if folder, err = charmap.Windows1252.NewDecoder().Bytes(buf); err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *CharacterSprite) Render(gls *opengl.State, renderInfo graphic.RenderInfo, char *CharState) {
	offset := [2]float32{}
	action := actionindex.GetActionIndex(char.State)

	if action != actionindex.Dead && action != actionindex.Sitting {
		s.renderElement(gls, renderInfo, char, CharacterSpriteElementShadow, &offset)
	}

	s.renderElement(gls, renderInfo, char, CharacterSpriteElementBody, &offset)

	s.renderElement(gls, renderInfo, char, CharacterSpriteElementHead, &offset)
}

func (s *CharacterSprite) renderElement(
	gls *opengl.State,
	renderInfo graphic.RenderInfo,
	char *CharState,
	elem CharacterSpriteElement,
	offset *[2]float32,
) {
	if len(s.files[elem].ACT.Actions) == 0 {
		return
	}

	actionIndex := actionindex.GetActionIndex(char.State)
	idx := int(actionIndex) + (int(char.Direction)+DirectionTable[FixedCameraDirection])%8

	if elem == CharacterSpriteElementShadow {
		idx = 0
	}

	action := s.files[elem].ACT.Actions[idx]
	fileSet := s.files[elem]

	frameCount := len(action.Frames)
	timeNeededForOneFrame := int64(action.Delay.Seconds() * (1.0 / FPSMultiplier))
	timeNeededForOneFrame = int64(math.Max(float64(timeNeededForOneFrame), 100))
	elapsedTime := int64(0)
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

	pos := [2]float32{0, 0}

	if len(frame.Positions) > 0 && elem != CharacterSpriteElementBody {
		pos[0] = offset[0] - float32(frame.Positions[frameIndex][0])
		pos[1] = offset[1] - float32(frame.Positions[frameIndex][1])
	}

	// Render all frames
	for _, layer := range frame.Layers {
		if layer.SpriteFrameIndex < 0 {
			continue
		}

		s.renderLayer(gls, renderInfo, layer, fileSet.SPR, pos, elem)
	}

	// Save offset reference
	if elem == CharacterSpriteElementBody && len(frame.Positions) > 0 {
		*offset = [2]float32{float32(frame.Positions[frameIndex][0]), float32(frame.Positions[frameIndex][1])}
	}
}

func (s *CharacterSprite) renderLayer(
	gls *opengl.State,
	renderInfo graphic.RenderInfo,
	layer *act.ActionFrameLayer,
	spr *spr.SpriteFile,
	offset [2]float32,
	elem CharacterSpriteElement,
) {
	frameIndex := int(layer.SpriteFrameIndex)
	if frameIndex < 0 {
		return
	}

	frame := spr.Frames[layer.Index]
	width, height := float32(frame.Width), float32(frame.Height)

	img := spr.ImageAt(frameIndex)
	texture, err := graphic.NewTextureFromImage(img)

	width *= layer.Scale[0] * SpriteScaleFactor * graphic.OnePixelSize
	height *= layer.Scale[1] * SpriteScaleFactor * graphic.OnePixelSize
	if err != nil {
		log.Fatal(err)
	}

	offsetX := (float32(layer.Position[0]) + offset[0]) * graphic.OnePixelSize
	offsetY := (float32(layer.Position[1]) + offset[1]) * graphic.OnePixelSize

	if layer.Mirrored {
		width = -width
	}

	sprite := graphic.NewSprite(1.0, 1.0, texture)
	sprite.SetPosition(mgl32.Vec3{s.Position().X() + offsetX, s.Position().Y(), 0})

	log.Printf(
		"elem=(%s) size=(%v, %v), scale=(%+v) position=(%+v), rotation=(%+v)\n",
		elem,
		width/graphic.OnePixelSize,
		height/graphic.OnePixelSize,
		sprite.Scale(),
		sprite.Position(),
		sprite.Rotation(),
	)

	{
		sprite.Texture.Bind(0)

		view := renderInfo.ViewMatrix()
		viewu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("view\x00"))
		gl.UniformMatrix4fv(viewu, 1, false, &view[0])

		model := sprite.Model()
		modelu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("model\x00"))
		gl.UniformMatrix4fv(modelu, 1, false, &model[0])

		projection := renderInfo.ProjectionMatrix()
		projectionu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("projection\x00"))
		gl.UniformMatrix4fv(projectionu, 1, false, &projection[0])

		size := mgl32.Vec2{width, height}
		sizeu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("size\x00"))
		gl.Uniform2fv(sizeu, 1, &size[0])

		offset := mgl32.Vec2{offsetX, offsetY}
		offsetu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("offset\x00"))
		gl.Uniform2fv(offsetu, 1, &offset[0])

		rotation := create3DRotationMatrix(0.0, graphic.Forward)
		rotationu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("rotation\x00"))
		gl.UniformMatrix4fv(rotationu, 1, false, &rotation[0])

		sprite.Render(gls)
	}

	s.sprites[elem] = sprite
}

func create3DRotationMatrix(angle float32, axis mgl32.Vec3) mgl32.Mat4 {
	iden := mgl32.Ident4()
	rotation := mgl32.HomogRotate3D(angle, axis)

	return iden.Mul4(rotation)
}
