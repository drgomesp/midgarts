package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	_ "github.com/joho/godotenv/autoload"
	"github.com/project-midgard/midgarts/internal/entity"
	"github.com/project-midgard/midgarts/internal/opengl"
	"github.com/project-midgard/midgarts/internal/system"
	"github.com/project-midgard/midgarts/internal/window"
	"github.com/project-midgard/midgarts/pkg/camera"
	"github.com/project-midgard/midgarts/pkg/character"
	"github.com/project-midgard/midgarts/pkg/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/character/statetype"
	"github.com/project-midgard/midgarts/pkg/fileformat/gnd"
	"github.com/project-midgard/midgarts/pkg/fileformat/grf"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindowWidth  = 1280
	WindowHeight = 720
	AspectRatio  = float32(WindowWidth) / float32(WindowHeight)
	FPS          = 120
)

var (
	GrfFilePath = os.Getenv("GRF_FILE_PATH")
)

func main() {
	runtime.LockOSThread()

	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	var win *sdl.Window
	if win, err = sdl.CreateWindow(
		"Midgarts Client",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		WindowWidth,
		WindowHeight,
		sdl.WINDOW_OPENGL,
	); err != nil {
		panic(err)
	}
	defer func() {
		_ = win.Destroy()
	}()

	context, err := win.GLCreateContext()
	if err != nil {
		panic(err)
	}
	defer sdl.GLDeleteContext(context)

	gls := opengl.InitOpenGL()

	var grfFile *grf.File
	if grfFile, err = grf.Load(GrfFilePath); err != nil {
		log.Fatal(err)
	}

	e, err := grfFile.GetEntry("data/prontera.gnd")
	if err != nil {
		log.Fatal(err)
	}
	groundFile, err := gnd.Load(e.Data)
	if err != nil {
		log.Fatal(err)
	}
	_ = groundFile

	gl.Viewport(0, 0, WindowWidth, WindowHeight)

	cam := camera.NewPerspectiveCamera(0.638, AspectRatio, 0.1, 1000.0)
	cam.ResetAngleAndY(WindowWidth, WindowHeight)

	ks := window.NewKeyState(win)

	w := ecs.World{}
	renderSys := system.NewCharacterRenderSystem(grfFile)
	actionSystem := system.NewCharacterActionSystem(grfFile)

	c1 := entity.NewCharacter(character.Male, jobspriteid.Blacksmith, 23)
	c1.HasShield = true
	c2 := entity.NewCharacter(character.Male, jobspriteid.Knight, 22)
	c2.HasShield = true
	c3 := entity.NewCharacter(character.Male, jobspriteid.Swordsman, 14)
	c3.HasShield = true
	c4 := entity.NewCharacter(character.Female, jobspriteid.Alchemist, 16)
	c5 := entity.NewCharacter(character.Male, jobspriteid.Bard, 19)
	c6 := entity.NewCharacter(character.Female, jobspriteid.MonkH, 27)
	c7 := entity.NewCharacter(character.Male, jobspriteid.Crusader2, 30)
	c8 := entity.NewCharacter(character.Male, jobspriteid.Assassin, 17)
	c9 := entity.NewCharacter(character.Male, jobspriteid.Monk, 15)
	c10 := entity.NewCharacter(character.Female, jobspriteid.Wizard, 19)
	c11 := entity.NewCharacter(character.Female, jobspriteid.Sage, 4)
	c12 := entity.NewCharacter(character.Female, jobspriteid.Dancer, 16)

	c1.SetPosition(mgl32.Vec3{0, 44, -1})
	c2.SetPosition(mgl32.Vec3{4, 44, 0})
	c3.SetPosition(mgl32.Vec3{8, 44, 0})
	c4.SetPosition(mgl32.Vec3{0, 40, 0})
	c5.SetPosition(mgl32.Vec3{4, 40, 0})
	c6.SetPosition(mgl32.Vec3{8, 40, 0})
	c7.SetPosition(mgl32.Vec3{0, 36, 0})
	c8.SetPosition(mgl32.Vec3{4, 36, 0})
	c9.SetPosition(mgl32.Vec3{8, 36, 0})
	c10.SetPosition(mgl32.Vec3{0, 32, 0})
	c11.SetPosition(mgl32.Vec3{4, 32, 0})
	c12.SetPosition(mgl32.Vec3{8, 32, 0})

	var actionable *system.CharacterActionable
	var renderable *system.CharacterRenderable
	w.AddSystemInterface(actionSystem, actionable, nil)
	w.AddSystemInterface(renderSys, renderable, nil)
	w.AddSystem(system.NewOpenGLRenderSystem(gls, cam, renderSys.RenderCommands))

	w.AddEntity(c2)
	w.AddEntity(c3)
	w.AddEntity(c4)
	w.AddEntity(c5)
	w.AddEntity(c6)
	w.AddEntity(c7)
	w.AddEntity(c8)
	w.AddEntity(c9)
	w.AddEntity(c10)
	w.AddEntity(c11)
	w.AddEntity(c12)
	w.AddEntity(c1)

	//c1.SetState(statetype.StandBy)

	shouldStop := false

	var refreshPeriod = time.Second / FPS

	for !shouldStop {
		frameStart := time.Now()
		gl.ClearColor(0, 0.5, 0.8, 1.0)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				shouldStop = true
				break
			}
		}

		const MovementRate = float32(0.065)

		p1 := c1.Position()

		// char controls
		if ks.Pressed(sdl.K_w) && ks.Pressed(sdl.K_d) {
			c1.Direction = directiontype.NorthEast
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X() - MovementRate, p1.Y() + MovementRate, p1.Z()})
		} else if ks.Pressed(sdl.K_w) && ks.Pressed(sdl.K_a) {
			c1.Direction = directiontype.NorthWest
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X() + MovementRate, p1.Y() + MovementRate, p1.Z()})
		} else if ks.Pressed(sdl.K_s) && ks.Pressed(sdl.K_d) {
			c1.Direction = directiontype.SouthEast
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X() - MovementRate, p1.Y() - MovementRate, p1.Z()})
		} else if ks.Pressed(sdl.K_s) && ks.Pressed(sdl.K_a) {
			c1.Direction = directiontype.SouthWest
			c1.SetPosition(mgl32.Vec3{p1.X() + MovementRate, p1.Y() - MovementRate, p1.Z()})
			c1.SetState(statetype.Walking)
		} else if ks.Pressed(sdl.K_w) {
			c1.Direction = directiontype.North
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X(), p1.Y() + MovementRate, p1.Z()})
		} else if ks.Pressed(sdl.K_s) {
			c1.Direction = directiontype.South
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X(), p1.Y() - MovementRate, p1.Z()})
		} else if ks.Pressed(sdl.K_d) {
			c1.Direction = directiontype.East
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X() - MovementRate, p1.Y(), p1.Z()})
		} else if ks.Pressed(sdl.K_a) {
			c1.Direction = directiontype.West
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X() + MovementRate, p1.Y(), p1.Z()})
			//
			//p2 := c2.Position()
			//c2.Direction = directiontype.West
			//c2.SetState(statetype.Walking)
			//c2.SetPosition(mgl32.Vec3{p2.X() + MovementRate, p2.Y(), p2.Z()})
		} else {
			c1.SetState(statetype.StandBy)
		}

		// camera controls
		if ks.Pressed(sdl.K_z) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() + 0.2})
		} else if ks.Pressed(sdl.K_x) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() - 0.2})
		}

		//c2.SetState(statetype.StandBy)

		frameStart = time.Now()
		frameDelta := frameStart.Sub(frameStart)

		w.Update(float32(frameDelta.Seconds()))

		win.GLSwap()

		time.Sleep(refreshPeriod)
	}
}
