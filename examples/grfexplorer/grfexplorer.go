package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"

	g "github.com/AllenDang/giu"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
)

var grfFile *grf.File
var imageWidget = &g.ImageWidget{}
var fileInfoWidget g.Widget
var imageScaleMultiplier int32 = 1
var loadedImageName string
var currentEntry *grf.Entry

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
		g.SplitLayout("Split", g.DirectionHorizontal, true, 300,
			buildEntryTreeNodes(),
			g.Layout{
				fileInfoWidget,
				imageWidget,
				//g.SliderInt("SliderInt", &imageScaleMultiplier, 1, 4),
				g.Custom(func() {
					if g.IsItemActive() {
						loadImage(loadedImageName)
					}
				}),
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
	if strings.Contains(entryName, ".act") {
		var err error
		if currentEntry, err = grfFile.GetEntry(entryName); err != nil {
			panic("kurwa!")
		}

		actFile, err := act.Load(currentEntry.Data)
		log.Printf("actFile = %+v\n", actFile)
	}

	if strings.Contains(entryName, ".spr") {
		var err error
		if currentEntry, err = grfFile.GetEntry(entryName); err != nil {
			panic("kurwa!")
		}

		loadImage(entryName)
		loadFileInfo()
	}
}

func loadFileInfo() {
	sprFile, _ := spr.Load(currentEntry.Data)

	fileInfoWidget = g.Layout{
		g.Line(
			g.Group().Layout(
				g.Label("File Info"),
				g.Table("Table").
					Columns(
						g.Column(""),
						g.Column(""),
					).
					Rows(
						g.Row(g.Label("Width").Wrapped(true), g.Label(fmt.Sprintf("%d", sprFile.Frames[0].Width))),
						g.Row(g.Label("Height").Wrapped(true), g.Label(fmt.Sprintf("%d", sprFile.Frames[0].Height))),
					),
			),
		),
	}
}

func buildEntryTreeNodes() g.Layout {
	if grfFile == nil {
		return g.Layout{}
	}

	var nodes []interface{}
	grfFile.GetEntryTree().Traverse(grfFile.GetEntryTree().Root, func(n *grf.EntryTreeNode) {
		selectableNodes := make([]g.Widget, 0)
		var nodeEntries []interface{}

		for _, e := range grfFile.GetEntries(n.Value) {
			nodeEntries = append(nodeEntries, e.Name)
		}

		var decodedDirName []byte
		var err error
		if decodedDirName, err = charmap.Windows1252.NewDecoder().Bytes([]byte(n.Value)); err != nil {
			panic(err)
		}

		node := g.TreeNode(fmt.Sprintf("%s (%d)", decodedDirName, len(nodeEntries)))
		selectableNodes = g.RangeBuilder("selectableNodes", nodeEntries, func(i int, v interface{}) g.Widget {
			var decodedStr string
			var err error
			if decodedStr, err = charmap.Windows1252.NewDecoder().String(v.(string)); err != nil {
				panic(err)
			}

			return g.Selectable(decodedStr).OnClick(func() {
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

	sprFile, _ := spr.Load(currentEntry.Data)
	img := sprFile.ImageAt(0).(*image.RGBA)
	//mul := int(imageScaleMultiplier)
	//img = transform.Resize(img, img.Bounds().Max.X*mul, img.Bounds().Max.Y*mul, transform.Linear)

	go func() {
		spriteTexture, _ = g.NewTextureFromRgba(img)
		imageWidget = g.Image(spriteTexture).Size(float32(img.Bounds().Max.X), float32(img.Bounds().Max.Y))
		loadedImageName = name
	}()

	return nil
}
