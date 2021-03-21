package main

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/pkg/client/graphics"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
)

var monsterSprite *graphics.Sprite
var charSprite1 *graphics.CharacterSprite
var charSprite2 *graphics.CharacterSprite
var charSprite3 *graphics.CharacterSprite

type myScene struct{}

// Type uniquely defines your game type
func (*myScene) Type() string { return "myGame" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (*myScene) Preload() {
	var err error
	err = engo.Files.Load("textures/test.png")

	f, err := grf.Load("/home/drgomesp/grf/data.grf")
	monsterSprite, err = graphics.LoadSprite(f, `data/sprite/ork_warrior`)
	if err != nil {
		panic(err)
	}

	charSprite1, err = graphics.LoadCharacterSprite(f, jobspriteid.Thief)
	if err != nil {
		panic(err)
	}

	charSprite2, err = graphics.LoadCharacterSprite(f, jobspriteid.Merchant)
	if err != nil {
		panic(err)
	}

	charSprite3, err = graphics.LoadCharacterSprite(f, jobspriteid.MonkH)
	if err != nil {
		panic(err)
	}
}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (*myScene) Setup(u engo.Updater) {
	common.SetBackground(color.White)

	w, _ := u.(*ecs.World)
	w.AddSystem(&common.RenderSystem{})

	charA := NewCharacterEntity(charSprite1, engo.Point{X: 0, Y: 0})
	charB := NewCharacterEntity(charSprite2, engo.Point{X: 50, Y: 100})
	charC := NewCharacterEntity(charSprite2, engo.Point{X: 0, Y: 150})
	charD := NewCharacterEntity(charSprite3, engo.Point{X: 200, Y: 100})

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&charA.BasicEntity, &charA.RenderComponent, &charA.SpaceComponent)
			sys.Add(&charB.BasicEntity, &charB.RenderComponent, &charB.SpaceComponent)
			sys.Add(&charC.BasicEntity, &charC.RenderComponent, &charC.SpaceComponent)
			sys.Add(&charD.BasicEntity, &charD.RenderComponent, &charD.SpaceComponent)
		}
	}
}

func NewCharacterEntity(sprite *graphics.CharacterSprite, initialPos engo.Point) *Character {
	texture := sprite.GetActionLayerTexture(0, 0)

	return &Character{
		BasicEntity: ecs.NewBasic(),
		RenderComponent: common.RenderComponent{
			Drawable: texture,
			Scale:    engo.Point{X: 1, Y: 1},
		},
		SpaceComponent: common.SpaceComponent{
			Position: initialPos,
			Width:    texture.Width(),
			Height:   texture.Height(),
		},
	}
}

type Character struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func main() {
	opts := engo.RunOptions{
		Title:  "Hello World",
		Width:  400,
		Height: 400,
	}
	engo.Run(opts, &myScene{})
}
