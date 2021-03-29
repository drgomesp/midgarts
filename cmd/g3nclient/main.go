package main

import (
	"log"
	"runtime"
	"time"

	"github.com/g3n/engine/geometry"

	"github.com/g3n/engine/light"

	"github.com/g3n/engine/util/helper"

	"github.com/g3n/engine/gui"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/renderer"

	"github.com/EngoEngine/ecs"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/util/logger"
	"github.com/g3n/engine/window"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
)

const (
	CameraPitch      = -60.0
	CameraYaw        = 270.0
	DefaultTargetFPS = 60
)

type ClientOption func(*MidgartsClient)

func WithTargetFPS(fps uint) ClientOption {
	return func(c *MidgartsClient) {
		c.targetFPS = fps
	}
}

type MidgartsClient struct {
	*app.Application

	log                  *logger.Logger
	windowManager        window.IWindow
	dataDir              string
	grfFile              *grf.File
	startTime, frameTime time.Time
	frameDelta           time.Duration
	targetFPS            uint
	world                *ecs.World

	scene  *core.Node
	camera camera.ICamera
}

func NewMidgartsClient(options ...ClientOption) (c *MidgartsClient, err error) {
	var grfFile *grf.File
	if grfFile, err = grf.Load("/home/drgomesp/grf/data.grf"); err != nil {
		log.Fatal(err)
	}

	defaultLogger := logger.New("Midgarts Client", nil)
	defaultLogger.AddWriter(logger.NewConsole(false))
	defaultLogger.SetFormat(logger.FTIME | logger.FMICROS)
	defaultLogger.SetLevel(logger.DEBUG)

	runtime.LockOSThread()

	a := app.App()

	scene := core.NewNode()
	gui.Manager().Set(scene)

	cam := camera.New(1.0)
	cam.SetPosition(0, 0, 3.0)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create a blue torus and add it to the scene
	geom := geometry.NewTorus(2, .4, 12, 32, math32.Pi*2)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(1))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	tex, _ := texture.NewTexture2DFromImage("assets/build/m/4016-1.png")
	mat2 := material.NewStandard(math32.NewColor("white"))
	mat2.AddTexture(tex)
	sprite := graphic.NewSprite(35, 75, mat2)
	char := NewCharacterEntity(sprite)

	world := &ecs.World{}
	var rend *Character
	world.AddSystemInterface(NewCharacterRenderSystem(defaultLogger, a.Renderer(), scene, cam), rend, nil)
	world.AddEntity(char)

	c = &MidgartsClient{
		Application: a,
		log:         defaultLogger,
		grfFile:     grfFile,
		targetFPS:   DefaultTargetFPS,
		scene:       scene,
	}

	for _, opt := range options {
		opt(c)
	}

	c.world = world

	return c, err
}

func (c *MidgartsClient) Run() {
	c.Application.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		c.Application.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		c.world.Update(float32(c.frameTime.Second()))
	})
}

func main() {
	var (
		c   *MidgartsClient
		err error
	)

	if c, err = NewMidgartsClient(WithTargetFPS(60)); err != nil {
		log.Fatal(err)
	}

	c.Run()
}
