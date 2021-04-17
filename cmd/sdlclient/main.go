package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/pkg/errors"
	rographic "github.com/project-midgard/midgarts/cmd/sdlclient/graphic"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
	"github.com/project-midgard/midgarts/internal/graphic"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
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
	cam.Transform.SetPosition(pos.X(), pos.Y(), pos.Z()-17)

	gl.Viewport(0, 0, int32(WindowWidth), int32(WindowHeight))
	gl.ClearColor(0, 0.5, 0.8, 1.0)

	chars := make([]*rographic.CharacterSprite, 10)

	//chars = append(chars, loadCharOrPanic(grfFile, character.Male, jobspriteid.Wizard, 11))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Male, jobspriteid.Blacksmith, 12))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Male, jobspriteid.Hunter, 13))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Male, jobspriteid.Assassin, 14))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Male, jobspriteid.Crusader, 15))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Male, jobspriteid.Sage, 16))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Male, jobspriteid.Rogue, 17))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Male, jobspriteid.Alchemist, 18))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Male, jobspriteid.Bard, 19))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Male, jobspriteid.Dancer, 20))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Novice, 1))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Swordsman, 2))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Magician, 3))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Archer, 4))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Alcolyte, 5))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Merchant, 6))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Thief, 7))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Monk, 8))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Knight, 9))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Priest, 10))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Wizard, 11))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Blacksmith, 12))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Hunter, 13))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Assassin, 14))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Crusader, 15))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Sage, 16))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Rogue, 17))
	//chars = append(chars, loadCharOrPanic(grfFile, character.Female, jobspriteid.Alchemist, 18))
	//cf19 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Bard, 19)
	//cf20 := loadCharOrPanic(grfFile, character.Female, jobspriteid.Dancer, 20)

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

		headIndex := rand.Intn(20-1) + 1
		chars[0] = loadCharOrPanic(grfFile, character.Female, jobspriteid.Knight2, headIndex)
		chars[1] = loadCharOrPanic(grfFile, character.Male, jobspriteid.Novice, headIndex)
		chars[2] = loadCharOrPanic(grfFile, character.Male, jobspriteid.Magician, headIndex)
		chars[3] = loadCharOrPanic(grfFile, character.Male, jobspriteid.Archer, headIndex)
		chars[4] = loadCharOrPanic(grfFile, character.Male, jobspriteid.Alcolyte, headIndex)
		chars[5] = loadCharOrPanic(grfFile, character.Male, jobspriteid.Merchant, headIndex)
		chars[6] = loadCharOrPanic(grfFile, character.Male, jobspriteid.Thief, headIndex)
		chars[7] = loadCharOrPanic(grfFile, character.Male, jobspriteid.Monk, headIndex)
		chars[8] = loadCharOrPanic(grfFile, character.Male, jobspriteid.Knight, headIndex)
		chars[9] = loadCharOrPanic(grfFile, character.Male, jobspriteid.Priest, headIndex)

		chars[0].SetPosition(0, 0, 0)
		chars[1].SetPosition(-12, 0, 0)
		chars[2].SetPosition(-10, 0, 0)
		chars[3].SetPosition(-8, 0, 0)
		chars[4].SetPosition(-6, 0, 0)
		chars[5].SetPosition(-4, 0, 0)
		chars[6].SetPosition(-2, 0, 0)
		chars[7].SetPosition(2, 0, 0)
		chars[8].SetPosition(4, 0, 0)
		chars[9].SetPosition(6, 0, 0)
		//chars[10].SetPosition(6, 0, 0)
		//chars[11].SetPosition(8, 0, 0)
		//chars[12].SetPosition(10, 0, 0)
		//chars[13].SetPosition(12, 0, 0)
		//chars[14].SetPosition(-12, 4, 0)
		//chars[15].SetPosition(-10, 4, 0)
		//chars[16].SetPosition(-8, 4, 0)
		//chars[17].SetPosition(-6, 4, 0)
		//chars[18].SetPosition(-4, 4, 0)
		//chars[19].SetPosition(-2, 4, 0)
		//chars[20].SetPosition(-12, 8, 0)
		//chars[21].SetPosition(-10, 8, 0)
		//chars[22].SetPosition(-8, 8, 0)
		//chars[23].SetPosition(-6, 8, 0)
		//chars[24].SetPosition(-4, 8, 0)
		//chars[25].SetPosition(-2, 8, 0)
		//chars[26].SetPosition(0, 8, 0)
		//chars[27].SetPosition(2, 8, 0)
		//chars[28].SetPosition(4, 8, 0)
		//chars[29].SetPosition(6, 8, 0)
		//chars[30].SetPosition(8, 8, 0)
		//chars[31].SetPosition(10, 8, 0)
		//chars[32].SetPosition(12, 8, 0)
		//chars[33].SetPosition(-12, -4, 0)
		//chars[34].SetPosition(-8, -4, 0)
		//chars[35].SetPosition(-6, -4, 0)
		//chars[36].SetPosition(-4, -4, 0)
		//chars[37].SetPosition(-2, -4, 0)
		//chars[38].SetPosition(0, -4, 0)
		//cf19.SetPosition(2, -4, 0)
		//cf20.SetPosition(4, -4, 0)

		//sin := math.Sin(counter)
		//cos := math.Cos(counter)

		charState := &rographic.CharState{
			Direction: directiontype.South,
			State:     statetype.Idle,
		}

		for _, c := range chars {
			c.Render(gls, cam, charState)
		}

		window.GLSwap()

		counter += 0.001

		time.Sleep(time.Millisecond * 250)
	}
}

func loadCharOrPanic(grfFile *grf.File, gender character.GenderType, jobspriteid jobspriteid.Type, headIndex int) *rographic.CharacterSprite {
	cm14, err := rographic.LoadCharacterSprite(grfFile, gender, jobspriteid, int32(headIndex))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not load character (%v, %v)\n", gender, jobspriteid))
	}
	return cm14
}
