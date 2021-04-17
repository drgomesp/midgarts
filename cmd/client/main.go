package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"

	"github.com/project-midgard/midgarts/pkg/graphic"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/internal/entity"
	system2 "github.com/project-midgard/midgarts/internal/system"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobid"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
)

var F = character.Female.String()
var M = character.Male.String()

var grfFile *grf.File

var charSpritesheets = map[string][]*graphic.SpritesheetResource{
	F: {
		jobid.Archer: {},
		jobid.Thief:  {},
		jobid.Monk:   {},
		jobid.MonkH:  {},
	},
	M: {
		jobid.Monk:  {},
		jobid.MonkH: {},
	},
}

var monsterSpritesheets = map[string]*graphic.SpritesheetResource{
	"ork_warrior": nil,
}

type myScene struct {
}

// Type uniquely defines your game type
func (*myScene) Type() string { return "myGame" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (s *myScene) Preload() {
	//common.AddShader(&common.SpriteShader{})

	var err error

	if grfFile, err = grf.Load("/home/drgomesp/grf/data.grf"); err != nil {
		log.Fatal(err)
	}

	if err = engo.Files.Load("build/m/3-1.xml"); err != nil {
		log.Fatal(err)
	}

	if err = engo.Files.Load("build/m/6-1.xml"); err != nil {
		log.Fatal(err)
	}

	if err = engo.Files.Load("build/m/5-1.xml"); err != nil {
		log.Fatal(err)
	}

	if err = engo.Files.Load("build/m/15-1.xml"); err != nil {
		log.Fatal(err)
	}

	if err = engo.Files.Load("build/m/15-1.xml"); err != nil {
		log.Fatal(err)
	}

	if err = engo.Files.Load("build/m/4016-1.xml"); err != nil {
		log.Fatal(err)
	}
	//
	//if err = engo.Files.Load("build/ork_warrior/2-1.xml"); err != nil {
	//	log.Fatal(err)
	//}

	//charSpritesheets[M][jobid.Archer] = graphics2.NewSpritesheetResource(
	//	common.NewAsymmetricSpritesheetFromFile(
	//		"build/f/3-1.png",
	//		BuildSpriteRegionsFromTextureAtlas(character.Female, jobid.Archer),
	//	))
	charSpritesheets[F][jobid.Thief] = graphic.NewSpritesheetResource(
		common.NewAsymmetricSpritesheetFromFile(
			"build/m/6-1.png",
			BuildSpriteRegionsFromTextureAtlas(character.Male, jobid.Thief),
		))
	charSpritesheets[F][jobid.Merchant] = graphic.NewSpritesheetResource(
		common.NewAsymmetricSpritesheetFromFile(
			"build/m/5-1.png",
			BuildSpriteRegionsFromTextureAtlas(character.Male, jobid.Merchant),
		))
	charSpritesheets[F][jobid.Monk] = graphic.NewSpritesheetResource(
		common.NewAsymmetricSpritesheetFromFile(
			"build/m/15-1.png",
			BuildSpriteRegionsFromTextureAtlas(character.Male, jobid.Monk),
		))
	charSpritesheets[M][jobid.Archer] = graphic.NewSpritesheetResource(
		common.NewAsymmetricSpritesheetFromFile(
			"build/m/3-1.png",
			BuildSpriteRegionsFromTextureAtlas(character.Male, jobid.Archer),
		))
	charSpritesheets[M][jobid.MonkH] = graphic.NewSpritesheetResource(
		common.NewAsymmetricSpritesheetFromFile(
			"build/m/4016-1.png",
			BuildSpriteRegionsFromTextureAtlas(character.Male, jobid.MonkH),
		))

	//monsterSpritesheets["ork_warrior"] = graphic.NewSpritesheetResource(
	//	common.NewAsymmetricSpritesheetFromFile(
	//		"build/ork_warrior/2-1.png",
	//		BuildMonsterSpriteRegionsFromTextureAtlas("ork_warrior"),
	//	))
}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (s *myScene) Setup(u engo.Updater) {
	engo.Input.RegisterButton("Top", engo.KeyArrowUp)
	engo.Input.RegisterButton("Right", engo.KeyArrowRight)
	engo.Input.RegisterButton("Bot", engo.KeyArrowDown)
	engo.Input.RegisterButton("Left", engo.KeyArrowLeft)

	engo.Input.RegisterButton("A", engo.KeyA)
	common.SetBackground(color.White)

	w, _ := u.(*ecs.World)

	var rend *common.Renderable
	var notRend *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, rend, notRend)

	var anim *common.Animationable
	var notAnim *common.NotAnimationable
	w.AddSystemInterface(system2.NewCharacterAnimationSystem(), anim, notAnim)

	heroA := s.CreateCharacter(engo.Point{X: 100, Y: 100}, character.Male, jobid.Archer)
	heroB := s.CreateCharacter(engo.Point{X: 100, Y: 200}, character.Female, jobid.Merchant)
	heroC := s.CreateCharacter(engo.Point{X: 100, Y: 300}, character.Female, jobid.Monk)
	heroD := s.CreateCharacter(engo.Point{X: 100, Y: 400}, character.Male, jobid.MonkH)
	//monsterA := s.CreateMonsterCharacter(engo.Point{X: 250, Y: 300}, "ork_warrior")

	w.AddEntity(heroA)
	w.AddEntity(heroB)
	w.AddEntity(heroC)
	w.AddEntity(heroD)
	//w.AddEntity(monsterA)

	for _, sys := range w.Systems() {
		switch sys := sys.(type) {
		case *system2.CharacterAnimationSystem:
			{
				sys.Add(heroA)
				sys.Add(heroB)
				sys.Add(heroC)
				sys.Add(heroD)
				//sys.Add(monsterA)
			}
			break
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
) *entity.Character {
	spritesheetResource := charSpritesheets[gender.String()][jid]
	if spritesheetResource == nil {
		log.Fatalf("character spritesheetResource not found for jobid '%d' and gender '%d'", jid, gender)
	}

	actFile := character.LoadCharacterActionFile(grfFile, gender, jobspriteid.GetJobSpriteID(jid))
	char := entity.NewCharacterEntity(spritesheetResource, actFile, gender, jid)

	idleActionSprite := spritesheetResource.Spritesheet.Cell(0)
	char.SpaceComponent = common.SpaceComponent{
		Position: point,
		Width:    idleActionSprite.Width(),
		Height:   idleActionSprite.Height(),
	}
	char.RenderComponent = common.RenderComponent{
		Drawable: spritesheetResource.Spritesheet.Cell(0),
		Scale:    engo.Point{X: 1, Y: 1},
	}
	char.TargetPosition = point

	char.SetAction(statetype.Idle)

	return char
}

func (*myScene) CreateMonsterCharacter(point engo.Point, name string) *entity.Character {
	spritesheetResource := monsterSpritesheets[name]
	if spritesheetResource == nil {
		log.Fatalf("character spritesheetResource not found for jobid '%d' and gender '%d'", 0, 0)
	}

	var err error
	path := fmt.Sprintf("data/sprite/%s", name)
	var entry *grf.Entry
	if entry, err = grfFile.GetEntry(fmt.Sprintf("%s.act", path)); err != nil {
		log.Fatal(err)
	}

	actFile, err := act.Load(entry.Data)
	if err != nil {
		log.Fatal(err)
	}

	char := entity.NewCharacterEntity(spritesheetResource, actFile, 0, 0)

	char.SpaceComponent = common.SpaceComponent{
		Position: point,
		Width:    100,
		Height:   100,
	}
	char.RenderComponent = common.RenderComponent{
		Drawable: spritesheetResource.Spritesheet.Cell(0),
		Scale:    engo.Point{X: 1, Y: 1},
	}

	char.SetAction(statetype.Walking)

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

// NotRenderComponent is used to flag an entity as not in the RenderSystem even
// if it has the proper components
type NotRenderComponent struct{}

// GetNotRenderComponent implements the NotRenderable interface
func (n *NotRenderComponent) GetNotRenderComponent() *NotRenderComponent {
	return n
}

// NotRenderable is an interface used to flag an entity as not in the
// Rendersystem even if it has the proper components
type NotRenderable interface {
	GetNotRenderComponent() *NotRenderComponent
}
