package main

import (
	stdlog "log"
	"math"
	"os"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang-ui/nuklear/nk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/xlab/closer"

	"github.com/project-midgard/midgarts/internal/camera"
	"github.com/project-midgard/midgarts/internal/graphic"
	"github.com/project-midgard/midgarts/internal/graphic/geometry"
	"github.com/project-midgard/midgarts/internal/opengl"
	"github.com/project-midgard/midgarts/internal/window"
)

const (
	WindowWidth  = int32(1920)
	WindowHeight = int32(1080)
	AspectRatio  = float32(WindowWidth) / float32(WindowHeight)
	FPS          = 60

	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	var err error
	sdl.Init(sdl.INIT_EVERYTHING)

	win, err := sdl.CreateWindow("Nuklear Demo", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WindowWidth, WindowHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		closer.Fatalln(err)
	}
	defer win.Destroy()

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 2)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	glContext, err := win.GLCreateContext()
	if err != nil {
		closer.Fatalln(err)
	}
	_ = glContext

	width, height := win.GetSize()
	log.Printf("SDL2: created window %dx%d", width, height)

	if err := gl.Init(); err != nil {
		closer.Fatalln("opengl: init failed:", err)
	}

	gl.Viewport(0, 0, int32(width), int32(height))

	exitC := make(chan struct{}, 1)
	doneC := make(chan struct{}, 1)
	closer.Bind(func() {
		close(exitC)
		<-doneC
	})

	state := &State{
		bgColor: nk.NkRgba(28, 48, 62, 255),
	}
	_ = state
	fpsTicker := time.NewTicker(time.Second / 30)

	cam := camera.NewPerspectiveCamera(0.638, AspectRatio, 0.1, 1000.0)
	cam.ResetAngleAndY(WindowWidth, WindowHeight)

	ks := window.NewKeyState(win)

	gl.Viewport(0, 0, WindowWidth, WindowHeight)
	context, err := win.GLCreateContext()
	if err != nil {
		closer.Fatalln(err)
	}

	ctx := nk.NkPlatformInit(win, context, nk.PlatformInstallCallbacks)
	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	sansFont := nk.NkFontAtlasAddFromBytes(atlas, MustAsset("assets/FreeSans.ttf"), 16, nil)
	nk.NkFontStashEnd()
	if sansFont != nil {
		nk.NkStyleSetFont(ctx, sansFont.Handle())
	}

	gls := opengl.InitOpenGL()
	for {
		select {
		case <-exitC:
			nk.NkPlatformShutdown()
			sdl.Quit()
			fpsTicker.Stop()
			close(doneC)
			return
		case <-fpsTicker.C:
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					close(exitC)
				}

				ks.Update(event)
			}

			// camera controls
			if ks.Pressed(sdl.K_z) {
				cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() + 1.0})
			} else if ks.Pressed(sdl.K_x) {
				cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() - 1.0})
			} else if ks.Pressed(sdl.K_c) {
				cam.SetPosition(mgl32.Vec3{cam.Position().X() - 1.0, cam.Position().Y(), cam.Position().Z()})
			} else if ks.Pressed(sdl.K_v) {
				cam.SetPosition(mgl32.Vec3{cam.Position().X() + 1.0, cam.Position().Y(), cam.Position().Z()})
			}

			{
				gfxMain(win, state, ctx, cam, gls)
			}

		}
	}
}

func flag(v bool) int32 {
	if v {
		return 1
	}
	return 0
}

func gfxMain(win *sdl.Window, state *State, ctx *nk.Context, cam *camera.Camera, gls *opengl.State) {
	nk.NkPlatformNewFrame()

	// Layout
	bounds := nk.NkRect(50, 50, 230, 250)
	update := nk.NkBegin(ctx, "Demo", bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)

	if update > 0 {
		nk.NkLayoutRowStatic(ctx, 30, 80, 1)
		{
			if nk.NkButtonLabel(ctx, "button") > 0 {
				stdlog.Println("[INFO] button pressed!")
			}
		}
		nk.NkLayoutRowDynamic(ctx, 30, 2)
		{
			if nk.NkOptionLabel(ctx, "easy", flag(state.opt == Easy)) > 0 {
				state.opt = Easy
			}
			if nk.NkOptionLabel(ctx, "hard", flag(state.opt == Hard)) > 0 {
				state.opt = Hard
			}
		}
		nk.NkLayoutRowDynamic(ctx, 25, 1)
		{
			nk.NkPropertyInt(ctx, "Compression:", 0, &state.prop, 100, 10, 1)
		}
		nk.NkLayoutRowDynamic(ctx, 20, 1)
		{
			nk.NkLabel(ctx, "background:", nk.TextLeft)
		}
		nk.NkLayoutRowDynamic(ctx, 25, 1)
		{
			size := nk.NkVec2(nk.NkWidgetWidth(ctx), 400)
			if nk.NkComboBeginColor(ctx, state.bgColor, size) > 0 {
				nk.NkLayoutRowDynamic(ctx, 120, 1)
				cf := nk.NkColorCf(state.bgColor)
				cf = nk.NkColorPicker(ctx, cf, nk.ColorFormatRGBA)
				state.bgColor = nk.NkRgbCf(cf)
				nk.NkLayoutRowDynamic(ctx, 25, 1)
				r, g, b, a := state.bgColor.RGBAi()
				r = nk.NkPropertyi(ctx, "#R:", 0, r, 255, 1, 1)
				g = nk.NkPropertyi(ctx, "#G:", 0, g, 255, 1, 1)
				b = nk.NkPropertyi(ctx, "#B:", 0, b, 255, 1, 1)
				a = nk.NkPropertyi(ctx, "#A:", 0, a, 255, 1, 1)
				state.bgColor.SetRGBAi(r, g, b, a)
				nk.NkComboEnd(ctx)
			}
		}
	}
	nk.NkEnd(ctx)

	// Render
	bg := make([]float32, 4)
	nk.NkColorFv(bg, state.bgColor)
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	{
		gl.UseProgram(gls.Program().ID())
		plane := geometry.NewPlane(50, 60, nil)

		view := cam.ViewMatrix()
		viewu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("view\x00"))
		gl.UniformMatrix4fv(viewu, 1, false, &view[0])

		model := plane.Model()
		modelu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("model\x00"))
		gl.UniformMatrix4fv(modelu, 1, false, &model[0])

		projection := cam.ProjectionMatrix()
		projectionu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("projection\x00"))
		gl.UniformMatrix4fv(projectionu, 1, false, &projection[0])

		sizeu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("size\x00"))
		s := float32(1.0)
		gl.Uniform2fv(sizeu, 1, &s)

		offsetu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("offset\x00"))
		o := float32(0.0)
		gl.Uniform2fv(offsetu, 1, &o)

		iden := mgl32.Ident4()
		rotation := iden.Mul4(mgl32.HomogRotate3D(math.Pi/180, graphic.Backwards))
		rotationu := gl.GetUniformLocation(gls.Program().ID(), gl.Str("rotation\x00"))
		gl.UniformMatrix4fv(rotationu, 1, false, &rotation[0])

		plane.Render(gls)
	}

	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)

	win.GLSwap()
}

type Option uint8

const (
	Easy Option = 0
	Hard Option = 1
)

type State struct {
	bgColor nk.Color
	prop    int32
	opt     Option
}

func onError(code int32, msg string) {
	log.Printf("[ERR]: error %d: %s", code, msg)
}
