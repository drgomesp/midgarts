package main

import (
	"log"
	"math"
	"runtime"

	"github.com/pkg/errors"

	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
	"github.com/project-midgard/midgarts/internal/graphic"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
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
	cam := graphic.NewPerspectiveCamera(70.0, AspectRatio, 0.1, 1000.0)

	pos := cam.Position()
	cam.Transform.SetPosition(pos.X(), pos.Y(), pos.Z()-18)

	gl.Viewport(0, 0, int32(WindowWidth), int32(WindowHeight))
	gl.ClearColor(0, 0.5, 0.8, 1.0)

	cm1 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Novice, 1)
	cm2 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Swordsman, 2)
	cm3 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Magician, 3)
	cm4 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Archer, 4)
	cm5 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Alcolyte, 5)
	cm6 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Merchant, 6)
	cm7 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Thief, 7)
	cm8 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Monk, 8)
	cm9 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Knight, 9)
	cm10 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Priest, 10)
	cm11 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Wizard, 11)
	cm12 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Blacksmith, 12)
	cm13 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Hunter, 13)
	cm14 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Assassin, 14)
	cm15 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Crusader, 15)
	cm16 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Sage, 16)
	cm17 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Rogue, 17)
	cm18 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Alchemist, 18)
	cm19 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Bard, 19)
	//cm20 := loadCharOrPanic(grfFile, character.Male, jobspriteid.Dancer, 20)

	cf1 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Novice, 1)
	cf2 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Swordsman, 2)
	cf3 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Magician, 3)
	cf4 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Archer, 4)
	cf5 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Alcolyte, 5)
	cf6 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Merchant, 6)
	cf7 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Thief, 7)
	cf8 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Monk, 8)
	cf9 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Knight, 9)
	cf10 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Priest, 10)
	cf11 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Wizard, 11)
	cf12 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Blacksmith, 12)
	cf13 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Hunter, 13)
	cf14 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Assassin, 14)
	cf15 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Crusader, 15)
	cf16 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Sage, 16)
	cf17 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Rogue, 17)
	cf18 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Alchemist, 18)
	//cf19 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Bard, 19)
	//cf20 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Dancer, 20)

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

		cm1.SetPosition(-12, 0, 0)
		cm2.SetPosition(-10, 0, 0)
		cm3.SetPosition(-8, 0, 0)
		cm4.SetPosition(-6, 0, 0)
		cm5.SetPosition(-4, 0, 0)
		cm6.SetPosition(-2, 0, 0)
		cm7.SetPosition(0, 0, 0)
		cm8.SetPosition(2, 0, 0)
		cm9.SetPosition(4, 0, 0)
		cm10.SetPosition(6, 0, 0)
		cm11.SetPosition(8, 0, 0)
		cm12.SetPosition(10, 0, 0)
		cm13.SetPosition(12, 0, 0)
		cm14.SetPosition(-12, 4, 0)
		cm15.SetPosition(-10, 4, 0)
		cm16.SetPosition(-8, 4, 0)
		cm17.SetPosition(-6, 4, 0)
		cm18.SetPosition(-4, 4, 0)
		cm19.SetPosition(-2, 4, 0)
		//cm20.SetPosition(0, 4, 0)

		cf1.SetPosition(-12, 8, 0)
		cf2.SetPosition(-10, 8, 0)
		cf3.SetPosition(-8, 8, 0)
		cf4.SetPosition(-6, 8, 0)
		cf5.SetPosition(-4, 8, 0)
		cf6.SetPosition(-2, 8, 0)
		cf7.SetPosition(0, 8, 0)
		cf8.SetPosition(2, 8, 0)
		cf9.SetPosition(4, 8, 0)
		cf9.SetPosition(6, 8, 0)
		cf10.SetPosition(8, 8, 0)
		cf11.SetPosition(10, 8, 0)
		cf12.SetPosition(12, 8, 0)
		cf13.SetPosition(-12, -4, 0)
		cf14.SetPosition(-8, -4, 0)
		cf15.SetPosition(-6, -4, 0)
		cf16.SetPosition(-4, -4, 0)
		cf17.SetPosition(-2, -4, 0)
		cf18.SetPosition(0, -4, 0)
		//cf19.SetPosition(2, -4, 0)
		//cf20.SetPosition(4, -4, 0)

		sin := math.Sin(counter)
		cos := math.Cos(counter)
		cm1.SetRotation(mgl32.Vec3{float32(sin) * 25, float32(cos) * 25, 0})

		char := &graphic.CharState{Action: actionindex.Idle}

		cm1.Render(gls, cam, char)
		cm2.Render(gls, cam, char)
		cm3.Render(gls, cam, char)
		cm4.Render(gls, cam, char)
		cm5.Render(gls, cam, char)
		cm6.Render(gls, cam, char)
		cm7.Render(gls, cam, char)
		cm8.Render(gls, cam, char)
		cm9.Render(gls, cam, char)
		cm10.Render(gls, cam, char)
		cm11.Render(gls, cam, char)
		cm12.Render(gls, cam, char)
		cm13.Render(gls, cam, char)
		cm14.Render(gls, cam, char)
		cm15.Render(gls, cam, char)
		cm16.Render(gls, cam, char)
		cm17.Render(gls, cam, char)
		cm18.Render(gls, cam, char)
		cm19.Render(gls, cam, char)
		//cm20.Render(gls, cam, char)

		cf1.Render(gls, cam, char)
		cf2.Render(gls, cam, char)
		cf3.Render(gls, cam, char)
		cf4.Render(gls, cam, char)
		cf5.Render(gls, cam, char)
		cf6.Render(gls, cam, char)
		cf7.Render(gls, cam, char)
		cf8.Render(gls, cam, char)
		cf9.Render(gls, cam, char)
		cf10.Render(gls, cam, char)
		cf11.Render(gls, cam, char)
		cf12.Render(gls, cam, char)
		cf13.Render(gls, cam, char)
		cf14.Render(gls, cam, char)
		cf15.Render(gls, cam, char)
		cf16.Render(gls, cam, char)
		cf17.Render(gls, cam, char)
		cf18.Render(gls, cam, char)
		//cf19.Render(gls, cam, char)

		window.GLSwap()

		counter += 0.001
	}
}

func loadCharOrPanic(grfFile *grf.File, gender character.GenderType, jobspriteid jobspriteid.Type, headIndex int32) *graphic.CharacterSprite {
	cm14, err := graphic.LoadCharacterSprite(grfFile, gender, jobspriteid, headIndex)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not load character (%v, %v)\n", gender, jobspriteid))
	}
	return cm14
}
