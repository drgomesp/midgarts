package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/EngoEngine/engo/math"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	windowWidth  = 920
	windowHeight = 760

	vertexShaderSource = `
#version 330 core
layout(location = 0) in vec3 vp;
uniform mat4 mvp;

void main() {
	gl_Position = mvp * vec4(vp, 1.0);
}
	` + "\x00"

	fragmentShaderSource = `
#version 330 core
out vec4 frag_colour;
void main() {
	frag_colour = vec4(1.0, 1.0, 1.0, 1);
}
	` + "\x00"
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

	program := initOpenGL()

	cam := NewPerspectiveCamera(
		mgl32.Vec3{0, 0, -3},
		70.0,
		float32(windowWidth/windowHeight),
		0.1,
		1000.0,
	)

	gl.Enable(gl.DEPTH_TEST)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))
	gl.ClearColor(0, 0.5, 1.0, 1.0)

	triangleMesh := NewMesh(
		[]Vertex{
			{mgl32.Vec3{0, 0.5, 0}},
			{mgl32.Vec3{-0.5, -0.5, 0}},
			{mgl32.Vec3{0.5, -0.5, 0}},
		},
		[]uint32{0, 1, 2},
	)

	modelMatrix := mgl32.Ident4()
	counter := float32(0.0)

	shouldStop := false
	for !shouldStop {
		sin := math.Sin(counter)
		cos := math.Cos(counter)

		mvp := cam.ViewProjectionMatrix().Mul4(modelMatrix)
		mvpUniform := gl.GetUniformLocation(program, gl.Str("mvp\x00"))
		gl.UniformMatrix4fv(mvpUniform, 1, false, &mvp[0])

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				shouldStop = true
				break
			}
		}

		pos := cam.Position()
		_ = cos
		cam.SetPosition(pos.Add(mgl32.Vec3{sin, 0, -cos}))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)
		gl.UniformMatrix4fv(mvpUniform, 1, false, &mvp[0])

		triangleMesh.Draw()

		window.GLSwap()

		counter += 0.001
	}
}

// initOpenGL initializes OpenGL and returns an initialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
