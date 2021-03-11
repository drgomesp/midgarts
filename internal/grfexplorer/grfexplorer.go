package grfexplorer

import (
	"image"
	"log"
	"os"

	"github.com/project-midgard/midgarts/internal/fileformat/spr"

	"github.com/project-midgard/midgarts/internal/fileformat/grf"

	g "github.com/AllenDang/giu"
)

var grfFile *grf.File
var imageWidget = &g.ImageWidget{}

// Run implements the main program loop of the demo. It returns when the platform signals to stop.
// This demo application shows some basic features of ImGui, as well as exposing the standard demo window.
func Run() {
	g.SingleWindowWithMenuBar("splitter").Layout(
		g.MenuBar().Layout(
			g.Menu("File").Layout(
				g.MenuItem("Open").OnClick(onOpenFile),
				g.MenuItem("Save"),
				// You could add any kind of widget here, not just menu item.
				g.Menu("Save as ...").Layout(
					g.MenuItem("Excel file"),
					g.MenuItem("CSV file"),
					g.Button("Button inside menu"),
				),
			),
		),
		g.SplitLayout("Split", g.DirectionHorizontal, true, 200,
			buildEntryTreeNodes(),
			g.Layout{
				imageWidget,
			},
		),
	)
}

func onOpenFile() {
	log.Println("loading GRF file...")

	var err error
	grfFile, err = grf.Load(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func onClickEntry(entryName string) {
	log.Println(entryName)
	loadImage(entryName)
}

func buildEntryTreeNodes() g.Layout {
	if grfFile == nil {
		return g.Layout{}
	}

	var nodes []interface{}
	grfFile.GetEntryTree().Traverse(grfFile.GetEntryTree().Root, func(n *grf.EntryTreeNode) {
		node := g.TreeNode(n.Value)
		selectableNodes := make([]g.Widget, 0)
		var nodeEntries []interface{}

		for _, e := range grfFile.GetEntries(n.Value) {
			nodeEntries = append(nodeEntries, e.Name)
		}

		selectableNodes = g.RangeBuilder("selectableNodes", nodeEntries, func(i int, v interface{}) g.Widget {
			return g.Selectable(v.(string)).OnClick(func() {
				onClickEntry(v.(string))
			})
		})

		node.Layout(selectableNodes...)
		nodes = append(nodes, node)
	})

	tree := g.RangeBuilder("entries", nodes, func(i int, v interface{}) g.Widget {
		return v.(g.Widget)
	})

	return g.Layout{tree}
}

var spriteTexture *g.Texture

func loadImage(name string) *g.Texture {
	if grfFile == nil {
		return nil
	}

	entry, _ := grfFile.GetEntry(name)
	sprFile, _ := spr.Load(entry.Data)
	img := sprFile.ImageAt(0).(*image.RGBA)
	go func() {
		spriteTexture, _ = g.NewTextureFromRgba(img)
		imageWidget = g.Image(spriteTexture).Size(float32(img.Bounds().Max.X), float32(img.Bounds().Max.Y))
	}()

	return nil
}
