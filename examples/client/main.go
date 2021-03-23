package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"

	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"

	"github.com/project-midgard/midgarts/pkg/client"

	"github.com/project-midgard/midgarts/pkg/common/character"

	"github.com/project-midgard/midgarts/pkg/common/character/jobid"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/pkg/client/graphics"
)

var F = character.Female.String()
var M = character.Male.String()

var monsterSprite *graphics.CharacterSprite
var charSprite1 *graphics.CharacterSprite
var charSprite2 *graphics.CharacterSprite
var charSprite3 *graphics.CharacterSprite
var charSprite4 *graphics.CharacterSprite

var spritesheets = map[string][]*common.Spritesheet{
	F: {
		jobid.Archer: {},
		jobid.Thief:  {},
		jobid.Monk:   {},
	},
	M: {
		jobid.Monk: {},
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

	spritesheets[M][jobid.Archer] = common.NewAsymmetricSpritesheetFromFile(
		"build/m/3-1.png",
		BuildSpriteRegionsFromTextureAtlas(character.Male, jobid.Archer),
	)
	spritesheets[F][jobid.Thief] = common.NewAsymmetricSpritesheetFromFile(
		"build/m/6-1.png",
		BuildSpriteRegionsFromTextureAtlas(character.Male, jobid.Thief),
	)
	spritesheets[F][jobid.Monk] = common.NewAsymmetricSpritesheetFromFile(
		"build/f/15-1.png",
		BuildSpriteRegionsFromTextureAtlas(character.Female, jobid.Monk),
	)
	spritesheets[M][jobid.Monk] = common.NewAsymmetricSpritesheetFromFile(
		"build/m/15-1.png",
		BuildSpriteRegionsFromTextureAtlas(character.Male, jobid.Monk),
	)

	//F, err := grf.Load("/home/drgomesp/grf/data.grf")
	//tmp, err := graphics.LoadSprite(F, `data/sprite/ork_warrior`)
	//if err != nil {
	//	panic(err)
	//}
	//monsterSprite, err = graphics.LoadMonsterSprite(tmp)
	//if err != nil {
	//	panic(err)
	//}
	//
	//charSprite1, err = graphics.LoadCharacterSprite(F, character.Female, jobspriteid.Merchant)
	//if err != nil {
	//	panic(err)
	//}
	//
	//charSprite2, err = graphics.LoadCharacterSprite(F, character.Male, jobspriteid.Merchant)
	//if err != nil {
	//	panic(err)
	//}
	//
	//charSprite3, err = graphics.LoadCharacterSprite(F, character.Male, jobspriteid.Thief)
	//if err != nil {
	//	panic(err)
	//}
	//
	//charSprite4, err = graphics.LoadCharacterSprite(F, character.Male, jobspriteid.Monk)
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

	heroA := s.CreateCharacter(engo.Point{X: 100, Y: 200}, character.Male, jobid.Archer)
	heroB := s.CreateCharacter(engo.Point{X: 150, Y: 200}, character.Female, jobid.Thief)
	heroC := s.CreateCharacter(engo.Point{X: 200, Y: 200}, character.Female, jobid.Monk)
	heroD := s.CreateCharacter(engo.Point{X: 250, Y: 200}, character.Male, jobid.Monk)

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

func (*myScene) CreateCharacter(
	point engo.Point,
	gender character.GenderType,
	jid jobid.Type,
) *client.Character {
	spritesheet := spritesheets[gender.String()][jid]
	if spritesheet == nil {
		log.Fatalf("character spritesheet not found for jobid '%d' and gender '%d'", jid, gender)
	}

	char := client.NewCharacter(spritesheet, gender, jid)

	char.SpaceComponent = common.SpaceComponent{
		Position: point,
		Width:    100,
		Height:   100,
	}
	char.RenderComponent = common.RenderComponent{
		Drawable: spritesheet.Cell(0),
		Scale:    engo.Point{X: 1, Y: 1},
	}

	char.SetCurrentAction(client.CharacterActions[actionindex.Idle])

	return char
}

func main() {
	opts := engo.RunOptions{
		Title:  "Hello World",
		Width:  400,
		Height: 400,
	}
	engo.Run(opts, &myScene{})
}
