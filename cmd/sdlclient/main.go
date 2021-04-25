package main

import (
	"log"
	"runtime"
	"time"

	"github.com/project-midgard/midgarts/pkg/common/character/statetype"

	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"

	"github.com/EngoEngine/ecs"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/project-midgard/midgarts/internal/entity"
	"github.com/project-midgard/midgarts/internal/opengl"
	"github.com/project-midgard/midgarts/internal/system"
	"github.com/project-midgard/midgarts/internal/window"
	"github.com/project-midgard/midgarts/pkg/camera"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindowWidth  = 1280
	WindowHeight = 720
	AspectRatio  = float32(WindowWidth) / float32(WindowHeight)
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
	if grfFile, err = grf.Load("/home/drgomesp/grf/data.grf"); err != nil {
		log.Fatal(err)
	}

	gl.Viewport(0, 0, WindowWidth, WindowHeight)

	cam := camera.NewPerspectiveCamera(0.638, AspectRatio, 0.1, 1000.0)
	cam.ResetAngleAndY(WindowWidth, WindowHeight)

	ks := window.NewKeyState(win)

	w := ecs.World{}
	renderSys := system.NewCharacterRenderSystem(grfFile)
	actionSystem := system.NewCharacterActionSystem(grfFile)

	c1 := entity.NewCharacter(character.Female, jobspriteid.Alchemist, 25)
	c2 := entity.NewCharacter(character.Female, jobspriteid.Knight, 29)
	c3 := entity.NewCharacter(character.Male, jobspriteid.Swordsman, 19)
	c4 := entity.NewCharacter(character.Female, jobspriteid.Crusader, 31)
	c5 := entity.NewCharacter(character.Male, jobspriteid.Alcolyte, 11)

	c1.SetPosition(mgl32.Vec3{0, 38, 0})
	c2.SetPosition(mgl32.Vec3{4, 38, 0})
	c3.SetPosition(mgl32.Vec3{8, 38, 0})
	c4.SetPosition(mgl32.Vec3{0, 34, 0})
	c5.SetPosition(mgl32.Vec3{4, 34, 0})

	var actionable *system.CharacterActionable
	var renderable *system.CharacterRenderable
	w.AddSystemInterface(actionSystem, actionable, nil)
	w.AddSystemInterface(renderSys, renderable, nil)
	w.AddSystem(system.NewOpenGLRenderSystem(gls, cam, renderSys.RenderCommands))

	w.AddEntity(c1)
	w.AddEntity(c2)
	w.AddEntity(c3)
	w.AddEntity(c4)
	w.AddEntity(c5)

	shouldStop := false

	frameStart := time.Now()

	for !shouldStop {
		gl.ClearColor(0, 0.5, 0.8, 1.0)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				shouldStop = true
				break
			}
		}

		// char controls
		if ks.Pressed(sdl.K_UP) && ks.Pressed(sdl.K_RIGHT) {
			c1.Direction = directiontype.NorthEast
			c1.SetState(statetype.Walking)
		} else if ks.Pressed(sdl.K_UP) && ks.Pressed(sdl.K_LEFT) {
			c1.Direction = directiontype.NorthWest
			c1.SetState(statetype.Walking)
		} else if ks.Pressed(sdl.K_DOWN) && ks.Pressed(sdl.K_RIGHT) {
			c1.Direction = directiontype.SouthEast
			c1.SetState(statetype.Walking)
		} else if ks.Pressed(sdl.K_DOWN) && ks.Pressed(sdl.K_LEFT) {
			c1.Direction = directiontype.SouthWest
			c1.SetState(statetype.Walking)
		} else if ks.Pressed(sdl.K_UP) {
			c1.Direction = directiontype.North
			c1.SetState(statetype.Walking)
		} else if ks.Pressed(sdl.K_DOWN) {
			c1.Direction = directiontype.South
			c1.SetState(statetype.Walking)
		} else if ks.Pressed(sdl.K_RIGHT) {
			c1.Direction = directiontype.East
			c1.SetState(statetype.Walking)
		} else if ks.Pressed(sdl.K_LEFT) {
			c1.Direction = directiontype.West
			c1.SetState(statetype.Walking)
		} else {
			c1.SetState(statetype.Idle)
		}

		// camera controls
		if ks.Pressed(sdl.K_w) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() + 0.2})
		} else if ks.Pressed(sdl.K_s) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() - 0.2})
		}

		now := time.Now()
		frameDelta := now.Sub(frameStart)
		frameStart = now

		w.Update(float32(frameDelta.Seconds()))

		win.GLSwap()
	}
}

//
//func loadCharOrPanic(grfFile *grf.File, gender character.GenderType, jobspriteid jobspriteid.Type, headIndex int) *entity.CharacterSprite {
//	c, err := entity.LoadCharacterSprite(grfFile, gender, jobspriteid, int32(headIndex))
//	if err != nil {
//		log.Fatal(errors.Wrapf(err, "could not load character (%v, %v)\n", gender, jobspriteid))
//	}
//	return c
//}
