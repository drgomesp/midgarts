package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"

	"github.com/project-midgard/midgarts/pkg/common/character"

	"github.com/project-midgard/midgarts/pkg/common/character/jobid"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/pkg/client/graphics"
)

var monsterSprite *graphics.CharacterSprite
var charSprite1 *graphics.CharacterSprite
var charSprite2 *graphics.CharacterSprite
var charSprite3 *graphics.CharacterSprite
var charSprite4 *graphics.CharacterSprite

var characterSpriteSheets = []map[string]*common.Spritesheet{
	jobid.Archer: {
		"f": nil,
		"m": nil,
	},
	jobid.Thief: {
		"f": nil,
		"m": nil,
	},
	jobid.Monk: {
		"f": nil,
		"m": nil,
	},
}

type myScene struct {
}

// Type uniquely defines your game type
func (*myScene) Type() string { return "myGame" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (*myScene) Preload() {
	var err error
	//err = engo.Files.Load("build/15-1.png")
	//err = engo.Files.Load("build/4016-1.png")

	if err = engo.Files.Load("build/m/3-1.xml"); err != nil {
		log.Fatal(err)
	}

	if err = engo.Files.Load("build/m/6-1.xml"); err != nil {
		log.Fatal(err)
	}

	if err = engo.Files.Load("build/f/15-1.xml"); err != nil {
		log.Fatal(err)
	}

	if err = engo.Files.Load("build/m/15-1.xml"); err != nil {
		log.Fatal(err)
	}

	characterSpriteSheets[jobid.Archer]["m"] = common.NewAsymmetricSpritesheetFromFile(
		"build/m/3-1.png",
		BuildSpriteRegionsFromTextureAtlas(character.Male, jobid.Archer),
	)
	characterSpriteSheets[jobid.Thief]["f"] = common.NewAsymmetricSpritesheetFromFile(
		"build/m/6-1.png",
		BuildSpriteRegionsFromTextureAtlas(character.Male, jobid.Thief),
	)
	characterSpriteSheets[jobid.Monk]["f"] = common.NewAsymmetricSpritesheetFromFile(
		"build/f/15-1.png",
		BuildSpriteRegionsFromTextureAtlas(character.Female, jobid.Monk),
	)
	characterSpriteSheets[jobid.Monk]["m"] = common.NewAsymmetricSpritesheetFromFile(
		"build/m/15-1.png",
		BuildSpriteRegionsFromTextureAtlas(character.Male, jobid.Monk),
	)

	//f, err := grf.Load("/home/drgomesp/grf/data.grf")
	//tmp, err := graphics.LoadSprite(f, `data/sprite/ork_warrior`)
	//if err != nil {
	//	panic(err)
	//}
	//monsterSprite, err = graphics.LoadMonsterSprite(tmp)
	//if err != nil {
	//	panic(err)
	//}
	//
	//charSprite1, err = graphics.LoadCharacterSprite(f, character.Female, jobspriteid.Merchant)
	//if err != nil {
	//	panic(err)
	//}
	//
	//charSprite2, err = graphics.LoadCharacterSprite(f, character.Male, jobspriteid.Merchant)
	//if err != nil {
	//	panic(err)
	//}
	//
	//charSprite3, err = graphics.LoadCharacterSprite(f, character.Male, jobspriteid.Thief)
	//if err != nil {
	//	panic(err)
	//}
	//
	//charSprite4, err = graphics.LoadCharacterSprite(f, character.Male, jobspriteid.Monk)
	//if err != nil {
	//	panic(err)
	//}
}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (s *myScene) Setup(u engo.Updater) {
	common.SetBackground(color.White)

	w, _ := u.(*ecs.World)
	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.AnimationSystem{})

	heroA := s.CreateEntity(engo.Point{X: 100, Y: 200}, jobid.Archer, character.Male)
	heroB := s.CreateEntity(engo.Point{X: 150, Y: 200}, jobid.Thief, character.Female)
	heroC := s.CreateEntity(engo.Point{X: 200, Y: 200}, jobid.Monk, character.Female)
	heroD := s.CreateEntity(engo.Point{X: 250, Y: 200}, jobid.Monk, character.Male)

	//monster := NewCharacterEntity(monsterSprite, engo.Point{X: 300, Y: 300}, 27)
	//charA := NewCharacterEntity(charSprite1, engo.Point{X: 0, Y: 0}, 96)
	//charB := NewCharacterEntity(charSprite2, engo.Point{X: 50, Y: 100}, 0)
	//charD := NewCharacterEntity(charSprite3, engo.Point{X: 200, Y: 100}, 0)
	//charE := NewCharacterEntity(charSprite4, engo.Point{X: 250, Y: 100}, 95)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			//sys.Add(&monster.BasicEntity, &monster.RenderComponent, &monster.SpaceComponent)
			//sys.Add(&charA.BasicEntity, &charA.RenderComponent, &charA.SpaceComponent)
			//sys.Add(&charB.BasicEntity, &charB.RenderComponent, &charB.SpaceComponent)
			//sys.Add(&charD.BasicEntity, &charD.RenderComponent, &charD.SpaceComponent)
			//sys.Add(&charE.BasicEntity, &charE.RenderComponent, &charE.SpaceComponent)
			sys.Add(&heroA.BasicEntity, &heroA.RenderComponent, &heroA.SpaceComponent)
			sys.Add(&heroB.BasicEntity, &heroB.RenderComponent, &heroB.SpaceComponent)
			sys.Add(&heroC.BasicEntity, &heroC.RenderComponent, &heroC.SpaceComponent)
			sys.Add(&heroD.BasicEntity, &heroD.RenderComponent, &heroD.SpaceComponent)
			break
		case *common.AnimationSystem:
			sys.Add(&heroA.BasicEntity, &heroA.AnimationComponent, &heroA.RenderComponent)
			sys.Add(&heroB.BasicEntity, &heroB.AnimationComponent, &heroB.RenderComponent)
			sys.Add(&heroC.BasicEntity, &heroC.AnimationComponent, &heroC.RenderComponent)
			sys.Add(&heroD.BasicEntity, &heroD.AnimationComponent, &heroD.RenderComponent)
		}
	}
}

func BuildSpriteRegionsFromTextureAtlas(gender character.GenderType, jid jobid.Type) []common.SpriteRegion {
	textureAtlas, err := engo.Files.Resource(fmt.Sprintf("build/%s/%d-1.xml", gender, jid))
	if err != nil {
		log.Fatal(err)
	}
	atlas := textureAtlas.(*common.TextureAtlasResource).Atlas
	regions := make([]common.SpriteRegion, len(atlas.SubTextures))

	for _, subTexture := range atlas.SubTextures {
		idx, err := strconv.ParseInt(strings.TrimSuffix(subTexture.Name, ".png"), 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		regions[idx] = common.SpriteRegion{
			Position: engo.Point{
				X: subTexture.X,
				Y: subTexture.Y,
			},
			Width:  int(subTexture.Width),
			Height: int(subTexture.Height),
		}
	}

	return regions
}

func (*myScene) CreateEntity(
	point engo.Point,
	jid jobid.Type,
	gender character.GenderType,
) *AnimatedEntity {
	spriteSheet := characterSpriteSheets[jid][gender.String()]
	if spriteSheet == nil {
		log.Fatalf("character spritesheet not found for jobid '%d' and gender '%d'", jid, gender)
	}

	entity := &AnimatedEntity{BasicEntity: ecs.NewBasic()}

	entity.SpaceComponent = common.SpaceComponent{
		Position: point,
		Width:    100,
		Height:   100,
	}
	entity.RenderComponent = common.RenderComponent{
		Drawable: spriteSheet.Cell(0),
		Scale:    engo.Point{X: 1, Y: 1},
	}
	entity.AnimationComponent = common.NewAnimationComponent(spriteSheet.Drawables(), .09)

	animationAction0 := &common.Animation{Name: "run", Frames: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}

	entity.AnimationComponent.AddAnimations([]*common.Animation{animationAction0})
	entity.AnimationComponent.AddDefaultAnimation(animationAction0)

	return entity
}

func NewCharacterEntity(sprite *graphics.CharacterSprite, initialPos engo.Point, initialActIndex int) *Character {
	texture := sprite.GetActionLayerTexture(initialActIndex, 0)

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

type AnimatedEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.AnimationComponent
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
