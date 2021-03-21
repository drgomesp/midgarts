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
}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (*myScene) Setup(u engo.Updater) {
	common.SetBackground(color.White)

	w, _ := u.(*ecs.World)
	w.AddSystem(&common.RenderSystem{})

	charA := Character{BasicEntity: ecs.NewBasic()}
	charA.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: 0},
		Width:    303,
		Height:   641,
	}
	charA.RenderComponent = common.RenderComponent{
		Drawable: monsterSprite.GetTextureAtIndex(0),
		Scale:    engo.Point{X: 1, Y: 1},
	}

	charB := Character{BasicEntity: ecs.NewBasic()}
	charB.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 50, Y: 100},
		Width:    303,
		Height:   641,
	}
	charB.RenderComponent = common.RenderComponent{
		Drawable: charSprite1.GetActionLayerTexture(0, 0),
		Scale:    engo.Point{X: 1, Y: 1},
	}

	charC := Character{BasicEntity: ecs.NewBasic()}
	charC.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 100, Y: 100},
		Width:    303,
		Height:   641,
	}
	charC.RenderComponent = common.RenderComponent{
		Drawable: charSprite2.GetActionLayerTexture(1, 0),
		Scale:    engo.Point{X: 1, Y: 1},
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&charA.BasicEntity, &charA.RenderComponent, &charA.SpaceComponent)
			sys.Add(&charB.BasicEntity, &charB.RenderComponent, &charB.SpaceComponent)
			sys.Add(&charC.BasicEntity, &charC.RenderComponent, &charC.SpaceComponent)
		}
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
