package main

import (
	"log"
	"runtime"

	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
	"github.com/project-midgard/midgarts/internal/graphic"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 1024
	windowHeight = 768
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

	var grfFile *grf.File
	if grfFile, err = grf.Load("/home/drgomesp/grf/data.grf"); err != nil {
		log.Fatal(err)
	}

	e, err := grfFile.GetEntry("data/sprite/ÀÎ°£Á·/¸Ó¸®Åë/³²/1_³².spr")
	if err != nil {
		log.Fatal(err)
	}
	sprFile, err := spr.Load(e.Data)
	if err != nil {
		log.Fatal(err)
	}
	img := sprFile.ImageAt(0)
	t, _ := graphic.NewTextureFromImage(img)
	w, h := float32(img.Bounds().Size().X)*graphic.OnePixelSize, float32(img.Bounds().Size().Y)*graphic.OnePixelSize
	h += h / 3

	s := graphic.NewSprite(w, h, t)
	s.SetPosition(0, 2, 0)

	gl.Viewport(0, 0, windowWidth, windowHeight)

	cam := graphic.NewPerspectiveCamera(
		70.0,
		float32(windowWidth/windowHeight),
		0.1,
		1000.0,
	)

	pos := cam.Position()
	cam.Transform.SetPosition(pos.X(), pos.Y(), pos.Z()-25)

	gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))
	gl.ClearColor(0, 0.5, 0.8, 1.0)

	cs3, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Novice)
	if err != nil {
		log.Fatal(err)
	}
	cs6, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Swordsman)
	if err != nil {
		log.Fatal(err)
	}
	cs4, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Alcolyte)
	if err != nil {
		log.Fatal(err)
	}
	cs8, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Thief)
	if err != nil {
		log.Fatal(err)
	}
	cs9, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Magician)
	if err != nil {
		log.Fatal(err)
	}
	cs1, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Archer)
	if err != nil {
		log.Fatal(err)
	}
	cs5, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Merchant)
	if err != nil {
		log.Fatal(err)
	}
	cs7, err := graphic.LoadCharacterSprite(grfFile, character.Male, jobspriteid.Monk)
	if err != nil {
		log.Fatal(err)
	}
	cs2, err := graphic.LoadCharacterSprite(grfFile, character.Female, jobspriteid.MonkH)
	if err != nil {
		log.Fatal(err)
	}

	p1 := graphic.NewPlane(2.5, 2.5)

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

		p1.SetPosition(5, 5, 0)

		cs3.SetPosition(6, 0, 0)
		cs6.SetPosition(4, 0, 0)
		cs4.SetPosition(2, 0, 0)
		cs8.SetPosition(0, 0, 0)
		cs1.SetPosition(-2, 0, 0)
		cs5.SetPosition(-4, 0, 0)
		cs7.SetPosition(-6, 0, 0)
		cs2.SetPosition(-8, 0, 0)
		cs9.SetPosition(-10, 0, 0)

		//mvp = cam.ViewProjectionMatrix().Mul4(p1.Model())
		//mvpu = gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		//gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
		//p1.Render(gls, cam)
		{
			mvp := cam.ViewProjectionMatrix().Mul4(cs3.Model())
			mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
			gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
			cs3.Render(gls, cam)
		}

		{
			mvp := cam.ViewProjectionMatrix().Mul4(cs6.Model())
			mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
			gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
			cs6.Render(gls, cam)
		}

		{

			mvp := cam.ViewProjectionMatrix().Mul4(cs4.Model())
			mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
			gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
			cs4.Render(gls, cam)
		}

		{
			mvp := cam.ViewProjectionMatrix().Mul4(cs8.Model())
			mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
			gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
			cs8.Render(gls, cam)
		}

		{
			//sin := math.Sin(counter)
			//cos := math.Cos(counter)
			//cs1.SetRotation(mgl32.Vec3{float32(sin) * 25, float32(cos) * 25, 0})
			mvp := cam.ViewProjectionMatrix().Mul4(cs1.Model())
			mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
			gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
			cs1.Render(gls, cam)
		}

		{
			mvp := cam.ViewProjectionMatrix().Mul4(cs5.Model())
			mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
			gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
			cs5.Render(gls, cam)
		}

		{
			mvp := cam.ViewProjectionMatrix().Mul4(cs7.Model())
			mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
			gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
			cs7.Render(gls, cam)
		}

		{
			mvp := cam.ViewProjectionMatrix().Mul4(cs2.Model())
			mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
			gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
			cs2.Render(gls, cam)
		}

		{
			mvp := cam.ViewProjectionMatrix().Mul4(cs9.Model())
			mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
			gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
			cs9.Render(gls, cam)
		}

		{
			t.Bind(0)
			mvp := cam.ViewProjectionMatrix().Mul4(s.Model())
			mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
			gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])
			s.Render(gls, cam)
		}

		window.GLSwap()

		counter += 0.001
	}
}
