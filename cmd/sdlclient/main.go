package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	graphic2 "github.com/project-midgard/midgarts/pkg/graphic"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/pkg/errors"
	rographic "github.com/project-midgard/midgarts/cmd/sdlclient/graphic"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
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
	cam := graphic2.NewPerspectiveCamera(70.0, AspectRatio, 0.1, 1000.0)

	pos := cam.Position()
	cam.Transform.SetPosition(pos.X(), pos.Y(), pos.Z()-17)

	gl.Viewport(0, 0, int32(WindowWidth), int32(WindowHeight))
	gl.ClearColor(0, 0.5, 0.8, 1.0)

	classes := []jobspriteid.Type{
		jobspriteid.Novice,
		jobspriteid.Swordsman,
		jobspriteid.Magician,
		jobspriteid.Archer,
		jobspriteid.Alcolyte,
		jobspriteid.Merchant,
		jobspriteid.Thief,
		jobspriteid.Knight,
		jobspriteid.Priest,
		jobspriteid.Wizard,
		jobspriteid.Blacksmith,
		jobspriteid.Hunter,
		jobspriteid.Assassin,
		jobspriteid.Knight2,
		jobspriteid.Crusader,
		jobspriteid.Monk,
		jobspriteid.Sage,
		jobspriteid.Rogue,
		jobspriteid.Alchemist,
		jobspriteid.Crusader2,
		jobspriteid.MonkH,
	}
	chars := make([]*rographic.CharacterSprite, 21)

	counter := 0.0
	shouldStop := false
	for !shouldStop {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				shouldStop = true
				break
			}
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(gls.Program().ID())

		for i, jid := range classes {
			chars[i] = loadCharOrPanic(grfFile, character.Female, jid, rand.Intn(20-1)+1)
		}

		chars[0].SetPosition(-16, 0, 0)
		chars[1].SetPosition(-14, 0, 0)
		chars[2].SetPosition(-12, 0, 0)
		chars[3].SetPosition(-10, 0, 0)
		chars[4].SetPosition(-8, 0, 0)
		chars[5].SetPosition(-6, 0, 0)
		chars[6].SetPosition(-4, 0, 0)
		chars[7].SetPosition(-2, 0, 0)
		chars[8].SetPosition(0, 0, 0)
		chars[9].SetPosition(0, 4, 0)
		chars[10].SetPosition(2, 4, 0)
		chars[11].SetPosition(4, 4, 0)
		chars[12].SetPosition(6, 4, 0)
		chars[13].SetPosition(8, 4, 0)
		chars[14].SetPosition(10, 4, 0)
		chars[15].SetPosition(12, 4, 0)
		chars[16].SetPosition(14, 4, 0)
		chars[17].SetPosition(16, 4, 0)
		chars[18].SetPosition(-14, -4, 0)
		chars[19].SetPosition(-12, -4, 0)
		chars[20].SetPosition(-10, -4, 0)

		charState := &rographic.CharState{
			Direction: directiontype.Type(rand.Intn(8-1) + 1),
			State:     statetype.Idle,
		}

		for _, c := range chars {
			c.Render(gls, cam, charState)
		}

		window.GLSwap()

		counter += 0.001

		time.Sleep(time.Millisecond * 250)
	}
}

func loadCharOrPanic(grfFile *grf.File, gender character.GenderType, jobspriteid jobspriteid.Type, headIndex int) *rographic.CharacterSprite {
	cm14, err := rographic.LoadCharacterSprite(grfFile, gender, jobspriteid, int32(headIndex))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not load character (%v, %v)\n", gender, jobspriteid))
	}
	return cm14
}
