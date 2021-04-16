package graphic

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
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

type fileSet struct {
	ACT *act.ActionFile
	SPR *spr.SpriteFile
}

type CharacterSprite struct {
	*Transform

	Gender character.GenderType

	files          [NumCharacterSpriteElements]fileSet
	elementSprites [NumCharacterSpriteElements]*Sprite
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
	bodyActFile, bodySprFile, err := f.GetActionAndSpriteFiles(bodyFilePath)
	headActFile, headSprFile, err := f.GetActionAndSpriteFiles(headFilePathf)

	if err != nil {
		log.Fatal(err)
	}

	characterSprite := &CharacterSprite{
		Transform: NewTransform(Origin),
		Gender:    gender,
		files: [NumCharacterSpriteElements]fileSet{
			CharacterSpriteElementShadow: {shadowActFile, shadowSprFile},
			CharacterSpriteElementHead:   {headActFile, headSprFile},
			CharacterSpriteElementBody:   {bodyActFile, bodySprFile},
		},
		elementSprites: [NumCharacterSpriteElements]*Sprite{},
	}

	return characterSprite, nil

}

func getDecodedFolder(buf []byte) (folder []byte, err error) {
	if folder, err = charmap.Windows1252.NewDecoder().Bytes(buf); err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *CharacterSprite) Render(gls *opengl.State, cam *Camera) {
	position := [2]float32{}

	s.renderElement(gls, cam, CharacterSpriteElementShadow, &position)

	s.renderElement(gls, cam, CharacterSpriteElementBody, &position)

	s.renderElement(gls, cam, CharacterSpriteElementHead, &position)
}

func (s *CharacterSprite) renderElement(
	gls *opengl.State,
	cam *Camera,
	elem CharacterSpriteElement,
	position *[2]float32,
) {
	fileSet := s.files[elem]

	actionIndex := actionindex.GetActionIndex(statetype.Idle)
	idx := int(actionIndex) +
		(int(directiontype.South)+DirectionTable[FixedCameraDirection])%8

	currentFrame := 0
	action := fileSet.ACT.Actions[idx]
	frame := action.Frames[currentFrame]

	pos := [2]float32{0, 0}

	if len(frame.Positions) > 0 && elem != CharacterSpriteElementBody {
		pos[0] = position[0] - float32(frame.Positions[0][0])
		pos[1] = position[1] - float32(frame.Positions[0][1])
	}

	// Render all frames
	for _, layer := range frame.Layers {
		s.renderLayer(gls, cam, layer, fileSet.SPR, pos, elem)
	}

	// Save position reference
	if elem == CharacterSpriteElementBody && len(frame.Positions) > 0 {
		*position = [2]float32{float32(frame.Positions[0][0]), float32(frame.Positions[0][1])}
	}
}

func (s *CharacterSprite) renderLayer(
	gls *opengl.State,
	cam *Camera,
	layer *act.ActionFrameLayer,
	spr *spr.SpriteFile,
	position [2]float32,
	elem CharacterSpriteElement,
) {
	currentFrame := 0
	frame := spr.Frames[currentFrame]
	width, height := float32(frame.Width), float32(frame.Height)

	width *= layer.Scale[0] * SpriteScaleFactor * OnePixelSize
	height *= layer.Scale[1] * SpriteScaleFactor * OnePixelSize

	texture, err := NewTextureFromImage(spr.ImageAt(currentFrame))
	if err != nil {
		log.Fatal(err)
	}

	offsetX := (float32(layer.Position[0]) + position[0]) * OnePixelSize
	offsetY := (float32(layer.Position[1]) + position[1]) * OnePixelSize

	s.elementSprites[elem] = NewSprite(width, height, texture)
	s.elementSprites[elem].SetPosition(s.position.X()+offsetX, s.position.Y()-offsetY, 0)

	log.Printf("elem=(%s) w=(%v) h=(%v) pos=(%v) offset=(%v, %v)\n", elem, width, height, position, offsetX, offsetY)

	{
		mvp := cam.ViewProjectionMatrix().Mul4(s.elementSprites[elem].Model())
		mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])

		s.elementSprites[elem].Texture.Bind(0)
		s.elementSprites[elem].Render(gls, cam)
	}
}
