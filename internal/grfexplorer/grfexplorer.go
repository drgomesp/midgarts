package grfexplorer

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"
	"time"

	"github.com/project-midgard/midgarts/internal/fileformat/grf"
	"github.com/project-midgard/midgarts/internal/fileformat/spr"

	"github.com/inkyblackness/imgui-go/v3"
)

type clipboard struct {
	platform Platform
}

func (board clipboard) Text() (string, error) {
	return board.platform.ClipboardText()
}

func (board clipboard) SetText(text string) {
	board.platform.SetClipboardText(text)
}

// Renderer covers rendering imgui draw data.
type Renderer interface {
	// PreRender causes the display buffer to be prepared for new output.
	PreRender(clearColor [3]float32)
	// Render draws the provided imgui draw data.
	Render(displaySize [2]float32, framebufferSize [2]float32, drawData imgui.DrawData)

	LoadImage(img *image.RGBA) (uint32, error)
}

const (
	sleepDuration = time.Millisecond * 25
)

// Run implements the main program loop of the demo. It returns when the platform signals to stop.
// This demo application shows some basic features of ImGui, as well as exposing the standard demo window.
func Run(p Platform, r Renderer) {
	grfFile, err := grf.Load(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	dir := "data\\sprite\\npc"
	entry := "data\\sprite\\npc\\4_f_kafra3.spr"
	e, err := grfFile.GetEntry(dir, entry)

	if err != nil {
		log.Fatal(err)
	}

	sprFile, err := spr.Load(e.Data)
	if err != nil {
		log.Fatal(err)
	}

	img := sprFile.ImageAt(0)

	path := fmt.Sprintf("./out/%s.png", entry)
	outputFile, err := os.Create(strings.Trim(path, `'`))
	if err != nil {
		log.Fatal(err)
	}

	if err = png.Encode(outputFile, img); err != nil {
		log.Fatal(err)
	}

	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, rect.Min, draw.Src)

	imgui.CurrentIO().SetClipboard(clipboard{platform: p})

	clearColor := [3]float32{0.0, 0.0, 0.0}

	showTreeViewWindow := true
	showSpriteViewWindow := true

	for !p.ShouldStop() {
		p.ProcessEvents()

		// Signal start of a new frame
		p.NewFrame()
		imgui.NewFrame()

		// MenuBar
		if imgui.BeginMainMenuBar() {
			if imgui.BeginMenu("File") {
				if imgui.MenuItem("Open GRF file...") {
					//showSpriteViewWindow = true
				}

				imgui.EndMenu()
			}

			if imgui.BeginMenu("Help") {
				imgui.EndMenu()
			}

			imgui.EndMainMenuBar()
		}

		if showTreeViewWindow {
			//imgui.SetNextWindowPos(imgui.Vec2{
			//	X: 0,
			//	Y: 18,
			//})
			//imgui.SetNextWindowSize(imgui.Vec2{
			//	X: 250,
			//	Y: 750,
			//})

			imgui.BeginV("Tree View", &showTreeViewWindow, windowFlags{
				noTitlebar:     false,
				noScrollbar:    false,
				noMenu:         true,
				noMove:         false,
				noResize:       true,
				noCollapse:     true,
				noNav:          false,
				noBackground:   false,
				noBringToFront: false,
				noDecoration:   false,
			}.combined())

			//for d := range grfFile.GetEntryDirectories() {
			if imgui.TreeNode(dir) {
				for _, e = range grfFile.GetEntries(dir) {
					if imgui.TreeNodeV(e.Name, imgui.TreeNodeFlagsLeaf) {
						imgui.TreePop()
					}
				}

				imgui.TreePop()

			}
			//}

			imgui.PushItemWidth(imgui.FontSize() * -12)
			imgui.End()
		}

		if showSpriteViewWindow {
			imgui.BeginV("Sprite View", &showSpriteViewWindow, 0)

			var imageID uint32
			imageID, err = r.LoadImage(rgba)
			imgui.Image(imgui.TextureID(imageID), imgui.Vec2{X: float32(img.Bounds().Max.X), Y: float32(img.Bounds().Max.Y)})

			imgui.End()
		}

		// Rendering
		imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

		r.PreRender(clearColor)
		// A this point, the application could perform its own rendering...
		// app.RenderScene()

		r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
		p.PostRender()

		// sleep to avoid 100% CPU usage for this demo
		<-time.After(sleepDuration)
	}
}
