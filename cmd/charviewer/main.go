package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/drgomesp/midgarts/internal/entity"
	"github.com/drgomesp/midgarts/internal/opengl"
	"github.com/drgomesp/midgarts/internal/system"
	"github.com/drgomesp/midgarts/internal/window"
	"github.com/drgomesp/midgarts/pkg/camera"
	"github.com/drgomesp/midgarts/pkg/character"
	"github.com/drgomesp/midgarts/pkg/character/jobspriteid"
	"github.com/drgomesp/midgarts/pkg/character/statetype"
	"github.com/drgomesp/midgarts/pkg/fileformat/grf"
	"github.com/drgomesp/midgarts/pkg/graphic/caching"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"runtime"
	"time"
)

const (
	H           = 480
	W           = 600
	AspectRatio = float32(W) / float32(H)
	FPS         = 120
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	runtime.LockOSThread()

	checkErr(sdl.Init(sdl.INIT_EVERYTHING), "failed to load sdl")

	defer sdl.Quit()

	win, err := sdl.CreateWindow(
		"charviewer",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		W,
		H,
		sdl.WINDOW_OPENGL,
	)
	checkErr(err)
	defer func() { _ = win.Destroy() }()

	checkErr(sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3))
	checkErr(sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3))
	checkErr(sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE))

	context, err := win.GLCreateContext()
	checkErr(err)
	defer sdl.GLDeleteContext(context)

	gls := opengl.InitOpenGL()

	var grfFile *grf.File
	if grfFile, err = grf.Load("./assets/grf/data.grf"); err != nil {
		log.Fatal().Err(err).Msg("failed to load grf file")
	}

	gl.Viewport(0, 0, W, H)

	cam := camera.NewPerspectiveCamera(0.638, AspectRatio, 0.1, 1000.0)
	cam.ResetAngleAndY(W, H)

	ks := window.NewKeyState(win)

	w := ecs.World{}
	renderSys := system.NewCharacterRenderSystem(grfFile, caching.NewCachedTextureProvider())

	c1 := entity.NewCharacter(character.Male, jobspriteid.Blacksmith, 23)
	c1.SetPosition(mgl32.Vec3{0, 44, -1})

	var renderable *system.CharacterRenderable
	w.AddSystemInterface(renderSys, renderable, nil)
	w.AddSystem(system.NewOpenGLRenderSystem(gls, cam, renderSys.RenderCommands))

	w.AddEntity(c1)
	c1.SetState(statetype.StandBy)

	shouldStop := false

	var refreshPeriod = time.Second / FPS

	for !shouldStop {
		frameStart := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				shouldStop = true
				break
			}

			ks.Update(event)
		}

		if ks.Pressed(sdl.K_z) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() + 0.2})
		} else if ks.Pressed(sdl.K_x) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() - 0.2})
		}

		frameStart = time.Now()
		frameDelta := frameStart.Sub(frameStart)

		w.Update(float32(frameDelta.Seconds()))

		win.GLSwap()

		time.Sleep(refreshPeriod)
	}
}

func checkErr(err error, msg ...interface{}) {
	if err != nil {
		log.Fatal().Err(err).Msgf("", msg...)
	}
}
