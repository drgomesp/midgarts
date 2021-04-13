package main

import (
	"log"
	"math"
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

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	gls := opengl.InitOpenGL()
	cam := graphic.NewPerspectiveCamera(
		70.0,
		float32(windowWidth/windowHeight),
		0.1,
		1000.0,
	)

	gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))
	gl.ClearColor(0, 0.5, 0.8, 1.0)

	//t1 := NewMesh(
	//	[]Vertex{
	//		{
	//			mgl32.Vec3{0, 0.75, 0},
	//			Red,
	//			mgl32.Vec2{0.1, 0.1},
	//		},
	//		{
	//			mgl32.Vec3{-0.75, -0.75, 0},
	//			Red,
	//			mgl32.Vec2{0.1, 0.1},
	//		},
	//		{
	//			mgl32.Vec3{0.75, -0.75, 0},
	//			Red,
	//			mgl32.Vec2{0.1, 0.1},
	//		},
	//	},
	//	[]uint32{0, 1, 2},
	//)
	//t1.SetPosition(mgl32.Vec3{-1, 1, 1})
	//
	//t2 := NewMesh(
	//	[]Vertex{
	//		{
	//			mgl32.Vec3{0, 0.75, 0},
	//			Red,
	//			mgl32.Vec2{0.5, 0.5},
	//		},
	//		{
	//			mgl32.Vec3{-0.75, -0.75, 0},
	//			Green,
	//			mgl32.Vec2{0.1, 0.1},
	//		},
	//		{
	//			mgl32.Vec3{0.75, -0.75, 0},
	//			Blue,
	//			mgl32.Vec2{0.1, 0.1},
	//		},
	//	},
	//	[]uint32{0, 1, 2},
	//)
	//t2.SetPosition(mgl32.Vec3{3, 0, 5})

	tex, err := NewTextureFromImage("assets/out/4016/f/0.png")
	if err != nil {
		log.Fatal(err)
	}

	w := float32(35) * OnePixelSize
	h := float32(75) * OnePixelSize
	//rect := NewMesh(
	//	[]Vertex{
	//		{
	//			mgl32.Vec3{w, h, 0},
	//			White,
	//			mgl32.Vec2{0, 0},
	//		},
	//		{
	//			mgl32.Vec3{-w, -h, 0},
	//			White,
	//			mgl32.Vec2{0.05, 0.05},
	//		},
	//		{
	//			mgl32.Vec3{w, -h, 0},
	//			White,
	//			mgl32.Vec2{0, 0.05},
	//		},
	//		{
	//			mgl32.Vec3{-w, h, 0},
	//			White,
	//			mgl32.Vec2{0.05, 0},
	//		},
	//	},
	//	[]uint32{0, 1, 2, 3, 1, 0},
	//)
	//rect.SetPosition(mgl32.Vec3{-1, -1, 0})

	_ = w
	_ = h

	//sprite := graphic.NewSprite(w, h)
	sprite := graphic.NewSprite(2.0, 2.0)
	sprite.SetPosition(mgl32.Vec3{-1, -1, 1})

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

		_ = tex
		tex.Bind(0)

		sin := math.Sin(counter)
		cos := math.Cos(counter)

		//t1.SetRotation(mgl32.Vec3{0, 0, counter * 50})

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(gls.Program().ID())

		//mvp := cam.ViewProjectionMatrix().Mul4(t1.Model())
		//mvpUniform := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		//gl.UniformMatrix4fv(mvpUniform, 1, false, &mvp[0])
		//t1.Render()

		//t2.SetRotation(mgl32.Vec3{sin * 25, cos * 25, 0})
		//mvp = cam.ViewProjectionMatrix().Mul4(t2.Model())
		//mvpUniform = gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		//gl.UniformMatrix4fv(mvpUniform, 1, false, &mvp[0])
		//t2.Render()

		//mvp := cam.ViewProjectionMatrix().Mul4(rect.Model())
		//mvpUniform := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		//gl.UniformMatrix4fv(mvpUniform, 1, false, &mvp[0])
		//
		//rect.Render()

		mvp := cam.ViewProjectionMatrix().Mul4(sprite.Model())
		mvpu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
		gl.UniformMatrix4fv(mvpu, 1, false, &mvp[0])

		_ = sin
		_ = cos
		//sprite.SetRotation(mgl32.Vec3{float32(sin) * 25, float32(cos) * 25, 0})
		sprite.Render(gls, cam)

		window.GLSwap()

		counter += 0.001
	}
}
