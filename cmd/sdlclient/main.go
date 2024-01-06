package main

import (
	"fmt"
	"os"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/project-midgard/midgarts/internal/camera"
	"github.com/project-midgard/midgarts/internal/character"
	"github.com/project-midgard/midgarts/internal/character/directiontype"
	"github.com/project-midgard/midgarts/internal/character/jobspriteid"
	"github.com/project-midgard/midgarts/internal/character/statetype"
	"github.com/project-midgard/midgarts/internal/entity"
	"github.com/project-midgard/midgarts/internal/fileformat/gat"
	"github.com/project-midgard/midgarts/internal/fileformat/grf"
	"github.com/project-midgard/midgarts/internal/graphic/caching"
	"github.com/project-midgard/midgarts/internal/system"
	"github.com/project-midgard/midgarts/internal/system/opengl"
	"github.com/project-midgard/midgarts/internal/window"
	"github.com/project-midgard/midgarts/pkg/version"
)

const (
	WindowWidth  = 960
	WindowHeight = 720
	AspectRatio  = float32(WindowWidth) / float32(WindowHeight)
	FPS          = 60
)

var (
	GrfFilePath = os.Getenv("GRF_FILE_PATH")
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatal().Err(err).Msg("failed to load sdl")
	}
	defer sdl.Quit()

	desktop, err := sdl.GetDesktopDisplayMode(0)
	if err != nil {
		log.Fatal().Err(err).Msg("getting desktop display mode")
	}

	var win *sdl.Window
	if win, err = sdl.CreateWindow(
		fmt.Sprintf("Midgarts Client - %s", version.Get()),
		desktop.W-WindowWidth,
		0,
		WindowWidth,
		WindowHeight,
		sdl.WINDOW_OPENGL,
	); err != nil {
		panic(err)
	}
	defer func() {
		_ = win.Destroy()
	}()

	err = sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	if err != nil {
		panic(err)
	}

	err = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 2)
	if err != nil {
		panic(err)
	}

	err = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	if err != nil {
		panic(err)
	}

	context, err := win.GLCreateContext()
	if err != nil {
		panic(err)
	}
	defer sdl.GLDeleteContext(context)

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Info().Msgf("OpenGL version: %s", version)

	var grfFile *grf.File
	if grfFile, err = grf.Load(GrfFilePath); err != nil {
		log.Fatal().Err(err).Msg("failed to load grf file")
	}

	e, err := grfFile.GetEntry("data/izlude.gat")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get gat entry")
	}
	groundAltitude, err := gat.Load(e.Data)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load gat")
	}
	_ = groundAltitude

	gl.Viewport(0, 0, WindowWidth, WindowHeight)

	cam := camera.NewPerspectiveCamera(0.638, AspectRatio, 0.1, 1000.0)
	cam.ResetAngleAndY(WindowWidth, WindowHeight)

	ks := window.NewKeyState(win)

	w := ecs.World{}
	renderSys := system.NewCharacterRenderSystem(grfFile, caching.NewCachedTextureProvider())
	actionSystem := system.NewCharacterActionSystem(grfFile)

	c1 := entity.NewCharacter(character.Male, jobspriteid.Knight, 23)
	c1.HasShield = true
	c2 := entity.NewCharacter(character.Male, jobspriteid.Knight, 22)
	c2.HasShield = true
	c3 := entity.NewCharacter(character.Male, jobspriteid.Swordsman, 14)
	c3.HasShield = true
	c4 := entity.NewCharacter(character.Female, jobspriteid.Alchemist, 16)
	c5 := entity.NewCharacter(character.Male, jobspriteid.Bard, 19)
	c6 := entity.NewCharacter(character.Female, jobspriteid.MonkH, 27)
	c7 := entity.NewCharacter(character.Male, jobspriteid.Crusader2, 30)
	c8 := entity.NewCharacter(character.Male, jobspriteid.Assassin, 17)
	c9 := entity.NewCharacter(character.Male, jobspriteid.Monk, 15)
	c10 := entity.NewCharacter(character.Female, jobspriteid.Wizard, 19)
	c11 := entity.NewCharacter(character.Female, jobspriteid.Sage, 4)
	c12 := entity.NewCharacter(character.Female, jobspriteid.Dancer, 16)

	c1.SetPosition(mgl32.Vec3{0, 44, -1})
	c2.SetPosition(mgl32.Vec3{4, 44, 0})
	c3.SetPosition(mgl32.Vec3{8, 44, 0})
	c4.SetPosition(mgl32.Vec3{0, 40, 0})
	c5.SetPosition(mgl32.Vec3{4, 40, 0})
	c6.SetPosition(mgl32.Vec3{8, 40, 0})
	c7.SetPosition(mgl32.Vec3{0, 36, 0})
	c8.SetPosition(mgl32.Vec3{4, 36, 0})
	c9.SetPosition(mgl32.Vec3{8, 36, 0})
	c10.SetPosition(mgl32.Vec3{0, 32, 0})
	c11.SetPosition(mgl32.Vec3{4, 32, 0})
	c12.SetPosition(mgl32.Vec3{8, 32, 0})

	var actionable *system.CharacterActionable
	var renderable *system.CharacterRenderable
	w.AddSystemInterface(actionSystem, actionable, nil)
	w.AddSystemInterface(renderSys, renderable, nil)
	w.AddSystem(opengl.NewOpenGLRenderSystem(cam, renderSys.RenderCommands))

	w.AddEntity(c1)
	w.AddEntity(c2)
	w.AddEntity(c3)
	w.AddEntity(c4)
	w.AddEntity(c5)
	w.AddEntity(c6)
	w.AddEntity(c7)
	w.AddEntity(c8)
	w.AddEntity(c9)
	w.AddEntity(c10)
	w.AddEntity(c11)
	w.AddEntity(c12)

	c1.SetState(statetype.StandBy)

	shouldStop := false

	var refreshPeriod = time.Second / FPS

	for !shouldStop {
		frameStart := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch eventType := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				shouldStop = true
				break
			case *sdl.MouseButtonEvent:
				halfWidth := WindowWidth / 2
				halfHeight := WindowHeight / 2

				if eventType.Type == sdl.MOUSEBUTTONDOWN {
					if int(eventType.X) < halfWidth {
						if int(eventType.Y) < halfHeight {
							c1.Direction = directiontype.NorthWest
						} else if int(eventType.Y) > halfHeight {
							c1.Direction = directiontype.SouthWest
						} else {
							c1.Direction = directiontype.West
						}
					}

					if int(eventType.X) > halfWidth {
						if int(eventType.Y) < halfHeight {
							c1.Direction = directiontype.NorthEast
						} else if int(eventType.Y) > halfHeight {
							c1.Direction = directiontype.SouthEast
						} else {
							c1.Direction = directiontype.East
						}
					}

					spew.Dump(eventType)
					//}
					break
				}
			}

			ks.Update(event)
		}

		const MovementRate = float32(0.065)

		p1 := c1.Position()

		// char controls
		if ks.Pressed(sdl.K_w) && ks.Pressed(sdl.K_d) {
			c1.Direction = directiontype.NorthEast
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X() - MovementRate, p1.Y() + MovementRate, p1.Z()})
		} else if ks.Pressed(sdl.K_w) && ks.Pressed(sdl.K_a) {
			c1.Direction = directiontype.NorthWest
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X() + MovementRate, p1.Y() + MovementRate, p1.Z()})
		} else if ks.Pressed(sdl.K_s) && ks.Pressed(sdl.K_d) {
			c1.Direction = directiontype.SouthEast
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X() - MovementRate, p1.Y() - MovementRate, p1.Z()})
		} else if ks.Pressed(sdl.K_s) && ks.Pressed(sdl.K_a) {
			c1.Direction = directiontype.SouthWest
			c1.SetPosition(mgl32.Vec3{p1.X() + MovementRate, p1.Y() - MovementRate, p1.Z()})
			c1.SetState(statetype.Walking)
		} else if ks.Pressed(sdl.K_w) {
			c1.Direction = directiontype.North
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X(), p1.Y() + MovementRate, p1.Z()})
		} else if ks.Pressed(sdl.K_s) {
			c1.Direction = directiontype.South
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X(), p1.Y() - MovementRate, p1.Z()})
		} else if ks.Pressed(sdl.K_d) {
			c1.Direction = directiontype.East
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X() - MovementRate, p1.Y(), p1.Z()})
		} else if ks.Pressed(sdl.K_a) {
			c1.Direction = directiontype.West
			c1.SetState(statetype.Walking)
			c1.SetPosition(mgl32.Vec3{p1.X() + MovementRate, p1.Y(), p1.Z()})
			//
			//p2 := c2.Position()
			//c2.Direction = directiontype.West
			//c2.SetState(statetype.Walking)
			//c2.SetPosition(mgl32.Vec3{p2.X() + MovementRate, p2.Y(), p2.Z()})
		} else {
			//c1.SetState(statetype.StandBy)
			c1.SetState(statetype.StandBy)
		}

		// camera controls
		if ks.Pressed(sdl.K_z) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() + 0.2})
		} else if ks.Pressed(sdl.K_x) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X(), cam.Position().Y(), cam.Position().Z() - 0.2})
		} else if ks.Pressed(sdl.K_c) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X() - 0.2, cam.Position().Y(), cam.Position().Z()})
		} else if ks.Pressed(sdl.K_v) {
			cam.SetPosition(mgl32.Vec3{cam.Position().X() + 0.2, cam.Position().Y(), cam.Position().Z()})
		}

		//c2.SetState(statetype.StandBy)

		frameStart = time.Now()
		frameDelta := frameStart.Sub(frameStart)

		w.Update(float32(frameDelta.Seconds()))

		win.GLSwap()

		time.Sleep(refreshPeriod)
	}
}
