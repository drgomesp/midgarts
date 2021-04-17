package graphic

import (
	"fmt"
	"log"
	"math"

	graphic2 "github.com/project-midgard/midgarts/pkg/graphic"

	"github.com/go-gl/gl/v3.3-core/gl"
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
}

type fileSet struct {
	ACT *act.ActionFile
	SPR *spr.SpriteFile
}

type CharacterSprite struct {
	*graphic2.Transform

	Gender character.GenderType

	files   [NumCharacterSpriteElements]fileSet
	sprites [NumCharacterSpriteElements]*graphic2.Sprite
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
		Transform: graphic2.NewTransform(graphic2.Origin),
		Gender:    gender,
		files: [NumCharacterSpriteElements]fileSet{
			CharacterSpriteElementShadow: {shadowActFile, shadowSprFile},
			CharacterSpriteElementHead:   {headActFile, headSprFile},
			CharacterSpriteElementBody:   {bodyActFile, bodySprFile},
		},
		sprites: [NumCharacterSpriteElements]*graphic2.Sprite{},
	}

	return characterSprite, nil
}

func getDecodedFolder(buf []byte) (folder []byte, err error) {
	if folder, err = charmap.Windows1252.NewDecoder().Bytes(buf); err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *CharacterSprite) Render(gls *opengl.State, cam *graphic2.Camera, char *CharState) {
	offset := [2]float32{}

	// TODO: this must come from character state
	action := actionindex.GetActionIndex(char.State)

	if action != actionindex.Dead && action != actionindex.Sitting {
		s.renderElement(gls, cam, char, CharacterSpriteElementShadow, &offset)
	}

	s.renderElement(gls, cam, char, CharacterSpriteElementBody, &offset)

	s.renderElement(gls, cam, char, CharacterSpriteElementHead, &offset)
}

func (s *CharacterSprite) renderElement(
	gls *opengl.State,
	cam *graphic2.Camera,
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

	// TODO: make this come from char entity
	playMode := actionplaymode.Repeat
	var frameIndex int64
	switch playMode {
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

		s.renderLayer(gls, cam, layer, fileSet.SPR, pos, elem)
	}

	// Save offset reference
	if elem == CharacterSpriteElementBody && len(frame.Positions) > 0 {
		*offset = [2]float32{float32(frame.Positions[frameIndex][0]), float32(frame.Positions[frameIndex][1])}
	}
}

func (s *CharacterSprite) renderLayer(
	gls *opengl.State,
	cam *graphic2.Camera,
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
	texture, err := graphic2.NewTextureFromImage(img)

	width *= layer.Scale[0] * SpriteScaleFactor * graphic2.OnePixelSize
	height *= layer.Scale[1] * SpriteScaleFactor * graphic2.OnePixelSize
	if err != nil {
		log.Fatal(err)
	}

	offsetX := (float32(layer.Position[0]) + offset[0]) * graphic2.OnePixelSize
	offsetY := (float32(layer.Position[1]) + offset[1]) * graphic2.OnePixelSize

	if layer.Mirrored {
		width = -width
	}

	sprite := graphic2.NewSprite(width, height, texture)
	sprite.SetPosition(s.Position().X()-offsetX, s.Position().Y()-offsetY, 0)

	log.Printf(
		"elem=(%s) size=(%v, %v), scale=(%+v) position=(%+v), rotation=(%+v)\n",
		elem,
		width/graphic2.OnePixelSize,
		height/graphic2.OnePixelSize,
		sprite.Scale(),
		sprite.Position(),
		sprite.Rotation(),
	)

	{
		mvp := cam.ViewProjectionMatrix().Mul4(sprite.Model())
		mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])

		sprite.Texture.Bind(0)
		sprite.Render(gls, cam)
	}

	s.sprites[elem] = sprite
}
