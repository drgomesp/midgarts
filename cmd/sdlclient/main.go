package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
	rographic "github.com/project-midgard/midgarts/cmd/sdlclient/graphic"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
	"github.com/project-midgard/midgarts/pkg/camera"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math/rand"
	"runtime"
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

	var window *sdl.Window
	if window, err = sdl.CreateWindow(
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
		_ = window.Destroy()
	}()

	context, err := window.GLCreateContext()
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

	gl.Viewport(0, 0, int32(WindowWidth), int32(WindowHeight))
	gl.ClearColor(0, 0.5, 0.8, 1.0)

	jobs := jobspriteid.All()
	//jobs := []jobspriteid.Type{jobspriteid.Knight2}
	chars := make([]*rographic.CharacterSprite, 0)
	for _, jid := range jobs {
		chars = append(chars, loadCharOrPanic(grfFile, character.Female, jid, rand.Intn(20-1)+1))
		chars = append(chars, loadCharOrPanic(grfFile, character.Male, jid, rand.Intn(20-1)+1))
	}

	chars[0].SetPosition(mgl32.Vec3{0, 42, 0})
	chars[1].SetPosition(mgl32.Vec3{2, 42, 0})
	chars[2].SetPosition(mgl32.Vec3{4, 42, 0})
	chars[3].SetPosition(mgl32.Vec3{6, 42, 0})
	chars[4].SetPosition(mgl32.Vec3{8, 42, 0})
	chars[5].SetPosition(mgl32.Vec3{10, 42, 0})
	chars[6].SetPosition(mgl32.Vec3{12, 42, 0})
	chars[7].SetPosition(mgl32.Vec3{14, 42, 0})
	chars[8].SetPosition(mgl32.Vec3{16, 42, 0})
	chars[9].SetPosition(mgl32.Vec3{18, 45, 0})
	chars[10].SetPosition(mgl32.Vec3{2, 38, 0})
	chars[11].SetPosition(mgl32.Vec3{4, 38, 0})
	chars[12].SetPosition(mgl32.Vec3{6, 38, 0})
	chars[13].SetPosition(mgl32.Vec3{8, 38, 0})
	chars[14].SetPosition(mgl32.Vec3{10, 38, 0})
	chars[15].SetPosition(mgl32.Vec3{12, 38, 0})
	chars[16].SetPosition(mgl32.Vec3{14, 38, 0})
	chars[17].SetPosition(mgl32.Vec3{16, 38, 0})
	chars[18].SetPosition(mgl32.Vec3{-14, 38, 0})
	chars[19].SetPosition(mgl32.Vec3{-12, 38, 0})
	chars[20].SetPosition(mgl32.Vec3{-10, 38, 0})
	chars[21].SetPosition(mgl32.Vec3{2, 34, 0})
	chars[22].SetPosition(mgl32.Vec3{4, 34, 0})
	chars[23].SetPosition(mgl32.Vec3{6, 34, 0})
	chars[24].SetPosition(mgl32.Vec3{8, 34, 0})
	chars[25].SetPosition(mgl32.Vec3{10, 34, 0})
	chars[26].SetPosition(mgl32.Vec3{12, 34, 0})
	chars[27].SetPosition(mgl32.Vec3{14, 34, 0})
	chars[28].SetPosition(mgl32.Vec3{16, 34, 0})
	chars[29].SetPosition(mgl32.Vec3{-14, 34, 0})
	chars[30].SetPosition(mgl32.Vec3{-12, 34, 0})
	chars[31].SetPosition(mgl32.Vec3{-10, 34, 0})
	chars[32].SetPosition(mgl32.Vec3{-8, 34, 0})
	chars[33].SetPosition(mgl32.Vec3{-6, 34, 0})
	chars[34].SetPosition(mgl32.Vec3{-4, 34, 0})
	chars[35].SetPosition(mgl32.Vec3{-2, 34, 0})
	chars[36].SetPosition(mgl32.Vec3{0, 34, 0})
	chars[37].SetPosition(mgl32.Vec3{2, 30, 0})
	chars[38].SetPosition(mgl32.Vec3{4, 30, 0})
	chars[39].SetPosition(mgl32.Vec3{6, 30, 0})
	chars[40].SetPosition(mgl32.Vec3{8, 30, 0})
	chars[41].SetPosition(mgl32.Vec3{10, 30, 0})

	counter := 0.0
	shouldStop := false
	charState := &rographic.CharState{
		//Direction: directiontype.Type(rand.Intn(8-1) + 1),
		Direction: directiontype.South,
		State:     statetype.Idle,
		PlayMode:  actionplaymode.Repeat,
	}

	for !shouldStop {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				shouldStop = true
				break

			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym
				pos := cam.Position()

				if keyCode == sdl.K_UP {
					charState.Direction = directiontype.North
				}
				if keyCode == sdl.K_DOWN {
					charState.Direction = directiontype.South
				}
				if keyCode == sdl.K_RIGHT {
					charState.Direction = directiontype.East
				}
				if keyCode == sdl.K_LEFT {
					charState.Direction = directiontype.West
				}
				if keyCode == sdl.K_w {
					cam.SetPosition(mgl32.Vec3{pos.X(), pos.Y(), pos.Z() + 0.2})
				}
				if keyCode == sdl.K_s {
					cam.SetPosition(mgl32.Vec3{pos.X(), pos.Y(), pos.Z() - 0.2})
				}
			}
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(gls.Program().ID())

		for _, c := range chars {
			c.Render(gls, cam, charState)
		}

		window.GLSwap()

		counter += 0.001
	}
}

func loadCharOrPanic(grfFile *grf.File, gender character.GenderType, jobspriteid jobspriteid.Type, headIndex int) *rographic.CharacterSprite {
	cm14, err := rographic.LoadCharacterSprite(grfFile, gender, jobspriteid, int32(headIndex))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not load character (%v, %v)\n", gender, jobspriteid))
	}
	return cm14
}
