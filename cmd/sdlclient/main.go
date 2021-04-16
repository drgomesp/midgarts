package main

import (
	"log"
	"runtime"

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

	cm1, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Novice, 1)
	cm2, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Swordsman, 2)
	cm3, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Magician, 3)
	cm4, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Archer, 4)
	cm5, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Alcolyte, 5)
	cm6, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Merchant, 6)
	cm7, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Thief, 7)
	cm8, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Monk, 8)
	cm9, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.MonkH, 9)

	cf1, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Novice, 10)
	cf2, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Swordsman, 11)
	cf3, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Magician, 12)
	cf4, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Archer, 13)
	cf5, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Alcolyte, 14)
	cf6, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Merchant, 15)
	cf7, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Thief, 16)
	cf8, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Monk, 17)
	cf9, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.MonkH, 18)

	if err != nil {
		log.Fatal(err)
	}

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

		cm1.SetPosition(-10, 0, 0)
		cm2.SetPosition(-8, 0, 0)
		cm3.SetPosition(-6, 0, 0)
		cm4.SetPosition(-4, 0, 0)
		cm5.SetPosition(-2, 0, 0)
		cm6.SetPosition(0, 0, 0)
		cm7.SetPosition(2, 0, 0)
		cm8.SetPosition(4, 0, 0)
		cm9.SetPosition(6, 0, 0)

		cf1.SetPosition(-10, 4, 0)
		cf2.SetPosition(-8, 4, 0)
		cf3.SetPosition(-6, 4, 0)
		cf4.SetPosition(-4, 4, 0)
		cf5.SetPosition(-2, 4, 0)
		cf6.SetPosition(0, 4, 0)
		cf7.SetPosition(2, 4, 0)
		cf8.SetPosition(4, 4, 0)
		cf9.SetPosition(6, 4, 0)

		//sin := math.Sin(counter)
		//cos := math.Cos(counter)
		//cs1.SetRotation(mgl32.Vec3{float32(sin) * 25, float32(cos) * 25, 0})
		cm1.Render(gls, cam)
		cm2.Render(gls, cam)
		cm3.Render(gls, cam)
		cm4.Render(gls, cam)
		cm5.Render(gls, cam)
		cm6.Render(gls, cam)
		cm7.Render(gls, cam)
		cm8.Render(gls, cam)
		cm9.Render(gls, cam)

		cf1.Render(gls, cam)
		cf2.Render(gls, cam)
		cf3.Render(gls, cam)
		cf4.Render(gls, cam)
		cf5.Render(gls, cam)
		cf6.Render(gls, cam)
		cf7.Render(gls, cam)
		cf8.Render(gls, cam)
		cf9.Render(gls, cam)

		window.GLSwap()

		counter += 0.001
	}
}
