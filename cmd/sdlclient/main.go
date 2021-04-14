package main

import (
	"log"
	"runtime"

	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
	"github.com/project-midgard/midgarts/internal/graphic"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 1024
	windowHeight = 768
)

func main() {
	runtime.LockOSThread()

	var err error
	var grfFile *grf.File
	if grfFile, err = grf.Load("/home/drgomesp/grf/data.grf"); err != nil {
		log.Fatal(err)
	}

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	var window *sdl.Window
	if window, err = sdl.CreateWindow(
		"Midgarts Client",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		windowWidth,
		windowHeight,
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

	gl.Viewport(0, 0, windowWidth, windowHeight)

	//cam := graphic.NewPerspectiveCamera(
	//	70.0,
	//	float32(windowWidth/windowHeight),
	//	0.1,
	//	1000.0,
	//)
	//aspect:=float32(windowWidth/windowHe)
	cam := graphic.NewOrthographicCamera(
		-windowWidth/2,
		windowWidth/2,
		-windowHeight/2,
		windowHeight/2,
	)

	pos := cam.Position()
	cam.Transform.SetPosition(pos.X(), pos.Y(), pos.Z())

	gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))
	gl.ClearColor(0, 0.5, 0.8, 1.0)

	cs1, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Archer)
	if err != nil {
		log.Fatal(err)
	}

	cs2, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.MonkH)
	if err != nil {
		log.Fatal(err)
	}

	cs3, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Novice)
	if err != nil {
		log.Fatal(err)
	}

	cs4, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Alcolyte)
	if err != nil {
		log.Fatal(err)
	}

	cs5, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Merchant)
	if err != nil {
		log.Fatal(err)
	}

	cs6, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Swordsman)
	if err != nil {
		log.Fatal(err)
	}

	p1 := graphic.NewPlane(100, 100)

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

		p1.SetPosition(-200, -400, 0)
		cs1.SetPosition(0, 0, 0)
		cs2.Sprite.SetPosition(2, 1, 0)
		cs3.Sprite.SetPosition(4, 1, 0)
		cs4.Sprite.SetPosition(-2, 1, 0)
		cs5.Sprite.SetPosition(-4, 1, 0)
		cs6.Sprite.SetPosition(-1, 1, 0)

		//sin := math.Sin(counter)
		//cos := math.Cos(counter)
		//cs1.SetRotation(mgl32.Vec3{float32(sin) * 25, float32(cos) * 25, 0})
		mvp := cam.ViewProjectionMatrix().Mul4(cs1.Model())
		mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
		cs1.Render(gls, cam)

		mvp = cam.ViewProjectionMatrix().Mul4(p1.Model())
		mvpu = gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
		p1.Render(gls, cam)

		//mvp = cam.ViewProjectionMatrix().Mul4(cs2.BodySprite.Model())
		//mvpu = gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		//gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
		//cs2.Render(gls, cam)
		//
		//mvp = cam.ViewProjectionMatrix().Mul4(cs3.BodySprite.Model())
		//mvpu = gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		//gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
		//cs3.Render(gls, cam)
		//
		//mvp = cam.ViewProjectionMatrix().Mul4(cs4.BodySprite.Model())
		//mvpu = gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		//gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
		//cs4.Render(gls, cam)
		//
		//mvp = cam.ViewProjectionMatrix().Mul4(cs5.BodySprite.Model())
		//mvpu = gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		//gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
		//cs5.Render(gls, cam)
		//
		//mvp = cam.ViewProjectionMatrix().Mul4(cs6.BodySprite.Model())
		//mvpu = gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		//gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
		//cs6.Render(gls, cam)

		window.GLSwap()

		counter += 0.001
	}
}
