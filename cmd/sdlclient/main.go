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

	cs3, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Novice)
	if err != nil {
		log.Fatal(err)
	}
	cs6, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Swordsman)
	if err != nil {
		log.Fatal(err)
	}
	cs4, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Alcolyte)
	if err != nil {
		log.Fatal(err)
	}
	cs8, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Thief)
	if err != nil {
		log.Fatal(err)
	}
	cs9, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Magician)
	if err != nil {
		log.Fatal(err)
	}
	cs1, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Archer)
	if err != nil {
		log.Fatal(err)
	}
	cs5, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Merchant)
	if err != nil {
		log.Fatal(err)
	}
	cs7, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.Monk)
	if err != nil {
		log.Fatal(err)
	}
	cs2, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.MonkH)
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

		cs3.SetPosition(6, 0, 0)
		cs6.SetPosition(4, 0, 0)
		cs4.SetPosition(2, 0, 0)
		cs8.SetPosition(0, 0, 0)
		cs1.SetPosition(-2, 0, 0)
		cs5.SetPosition(-4, 0, 0)
		cs7.SetPosition(-6, 0, 0)
		cs2.SetPosition(-8, 0, 0)
		cs9.SetPosition(-10, 0, 0)

		//sin := math.Sin(counter)
		//cos := math.Cos(counter)
		//cs1.SetRotation(mgl32.Vec3{float32(sin) * 25, float32(cos) * 25, 0})
		cs3.Render(gls, cam)
		cs6.Render(gls, cam)
		cs4.Render(gls, cam)
		cs8.Render(gls, cam)
		cs1.Render(gls, cam)
		cs5.Render(gls, cam)
		cs7.Render(gls, cam)
		cs2.Render(gls, cam)
		cs8.Render(gls, cam)

		window.GLSwap()

		counter += 0.001
	}
}
