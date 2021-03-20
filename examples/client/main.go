package main

import (
	"image/color"

	"github.com/project-midgard/midgarts/pkg/character"
	"github.com/project-midgard/midgarts/pkg/character/jobspriteid"

	"github.com/project-midgard/midgarts/pkg/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/graphics"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var spriteResource *graphics.Sprite

type myScene struct{}

// Type uniquely defines your game type
func (*myScene) Type() string { return "myGame" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (*myScene) Preload() {
	var err error
	err = engo.Files.Load("textures/test.png")

	f, err := grf.Load("/home/drgomesp/grf/data.grf")
	//spriteResource, err = graphics.LoadSprite(f, `data/sprite/ork_warrior`)
	//if err != nil {
	//	panic(err)
	//}

	charSprite, err := character.LoadCharacterSprite(f, jobspriteid.Swordsman)
	if err != nil {
		panic(err)
	}

	spriteResource = charSprite
}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (*myScene) Setup(u engo.Updater) {
	common.SetBackground(color.White)

	w, _ := u.(*ecs.World)
	w.AddSystem(&common.RenderSystem{})

	char := Character{BasicEntity: ecs.NewBasic()}
	char.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 10, Y: 10},
		Width:    303,
		Height:   641,
	}

	char.RenderComponent = common.RenderComponent{
		Drawable: spriteResource.Textures[0],
		Scale:    engo.Point{X: 2, Y: 2},
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&char.BasicEntity, &char.RenderComponent, &char.SpaceComponent)
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
