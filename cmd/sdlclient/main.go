package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/project-midgard/midgarts/cmd/sdlclient/graphic"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 920
	windowHeight = 760
	OnePixelSize = 1.0 / 35.0
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
		"test",
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
	cam := graphic.NewPerspectiveCamera(
		70.0,
		float32(windowWidth/windowHeight),
		0.1,
		1000.0,
	)

	gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))
	gl.ClearColor(0, 0.5, 0.8, 1.0)

	monkTexture, err := NewTextureFromImage("assets/out/15/m/0.png")
	if err != nil {
		log.Fatal(err)
	}
	champTexture, err := NewTextureFromImage("assets/out/4016/f/0.png")
	if err != nil {
		log.Fatal(err)
	}

	w := float32(35) * OnePixelSize
	h := float32(75) * OnePixelSize
	s1 := graphic.NewSprite(w, h)
	s1.SetPosition(mgl32.Vec3{-1, -1, 21})

	s2 := graphic.NewSprite(w, h)
	s2.SetPosition(mgl32.Vec3{-2, -2, 21})

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

		monkTexture.Bind(0)
		mvp := cam.ViewProjectionMatrix().Mul4(s1.Model())
		mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])

		//sin := math.Sin(counter)
		//cos := math.Cos(counter)
		//s1.SetRotation(mgl32.Vec3{float32(sin) * 25, float32(cos) * 25, 0})

		s1.Render(gls, cam)

		champTexture.Bind(0)
		mvp = cam.ViewProjectionMatrix().Mul4(s2.Model())
		mvpu = gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])

		s2.Render(gls, cam)

		window.GLSwap()

		counter += 0.001
	}
}
