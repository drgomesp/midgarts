package main

import (
	"log"
	"runtime"
	"time"

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
	"github.com/project-midgard/midgarts/pkg/common/character/jobid"
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

	log.Printf("Window Aspect Ratio = %f\n", AspectRatio)
	cam := camera.NewPerspectiveCamera(0.638, AspectRatio, 0.1, 1000.0)
	cam.ResetAngleAndY(WindowWidth, WindowHeight)

	w := ecs.World{}
	renderSys := system.NewCharacterRenderSystem(grfFile)
	c1 := entity.NewCharacter(character.Female, jobid.Knight, 29)
	c2 := entity.NewCharacter(character.Female, jobid.Crusader, 31)
	c1.SetPosition(mgl32.Vec3{0, 38, 0})
	c2.SetPosition(mgl32.Vec3{4, 38, 0})

	var renderable *system.CharacterRenderable
	w.AddSystemInterface(renderSys, renderable, nil)
	w.AddSystem(system.NewOpenGLRenderSystem(gls, cam, renderSys.RenderCommands))

	w.AddEntity(c1)
	w.AddEntity(c2)

	counter := 0.0
	shouldStop := false

	ks := window.NewKeyState(win)

	for !shouldStop {
		gl.ClearColor(0, 0.5, 0.8, 1.0)

		start := time.Now()

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
			c2.Direction = directiontype.NorthEast
		} else if ks.Pressed(sdl.K_UP) && ks.Pressed(sdl.K_LEFT) {
			c1.Direction = directiontype.NorthWest
			c2.Direction = directiontype.NorthWest
		} else if ks.Pressed(sdl.K_DOWN) && ks.Pressed(sdl.K_RIGHT) {
			c1.Direction = directiontype.SouthEast
			c2.Direction = directiontype.SouthEast
		} else if ks.Pressed(sdl.K_DOWN) && ks.Pressed(sdl.K_LEFT) {
			c1.Direction = directiontype.SouthWest
			c2.Direction = directiontype.SouthWest
		} else if ks.Pressed(sdl.K_UP) {
			c1.Direction = directiontype.North
			c2.Direction = directiontype.North
		} else if ks.Pressed(sdl.K_DOWN) {
			c1.Direction = directiontype.South
			c2.Direction = directiontype.South
		} else if ks.Pressed(sdl.K_RIGHT) {
			c1.Direction = directiontype.East
			c2.Direction = directiontype.East
		} else if ks.Pressed(sdl.K_LEFT) {
			c1.Direction = directiontype.West
			c2.Direction = directiontype.West
		}

		// camera controls
		if ks.Pressed(sdl.K_w) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() + 0.2})
		} else if ks.Pressed(sdl.K_s) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() - 0.2})
		}

		w.Update(float32(time.Since(start).Seconds()))

		win.GLSwap()

		counter += 0.001
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
