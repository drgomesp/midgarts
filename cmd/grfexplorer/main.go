package main

import (
	"fmt"
	"os"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/suapapa/go_hangul/encoding/cp949"

	"github.com/project-midgard/midgarts/internal/fileformat/act"
	"github.com/project-midgard/midgarts/internal/fileformat/grf"
	"github.com/project-midgard/midgarts/internal/fileformat/spr"
)

var grfFile *grf.File
var imageWidget = &g.ImageWidget{}
var fileInfoWidget g.Widget
var loadedImageName string
var currentEntry *grf.Entry
var splitSize float32 = 500
var font []byte

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var err error
	grfFile, err = grf.Load("./assets/grf/data.grf")
	noErr(err)

	font, err = os.ReadFile("./assets/Fonts/FreeSans.ttf")
	noErr(err)
}

func main() {
	wnd := g.NewMasterWindow("GRF Explorer", 800, 600, g.FocusedFlagsNone)
	g.Context.FontAtlas.SetDefaultFontFromBytes(font, 16)
	wnd.Run(Run)
}

func Run() {
	g.SingleWindow().Layout(
		g.SplitLayout(g.DirectionVertical, &splitSize,
			buildEntryTreeNodes(),
			g.Layout{
				fileInfoWidget,
				imageWidget,
				g.Custom(func() {
					if g.IsItemActive() {
						loadImage(loadedImageName)
					}
				}),
			},
		),
	)
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
		g.Table().
			Columns(
				g.TableColumn("File Info"),
				g.TableColumn(""),
			).
			Rows(
				g.TableRow(g.Label("Width").Wrapped(true), g.Label(fmt.Sprintf("%d", sprFile.Frames[0].Width))),
				g.TableRow(g.Label("Height").Wrapped(true), g.Label(fmt.Sprintf("%d", sprFile.Frames[0].Height))),
			).Flags(g.TableFlagsBordersH),
	}
}

func getDecodedFolder(buf []byte) (string, error) {
	folderNameBytes, err := cp949.From(buf)
	return string(folderNameBytes), err
}

func buildEntryTreeNodes() g.Layout {
	entries := grfFile.GetEntryTree()

	var nodes []interface{}

	entries.Traverse(entries.Root, func(n *grf.EntryTreeNode) {
		selectableNodes := make([]g.Widget, 0)
		nodeEntries := make([]interface{}, 0)
		//body_file_path=data/sprite/ÀÎ°£Á·/¸öÅë/³²/Á¦Ã¶°ø_³²

		decodedFolderA, err := getDecodedFolder([]byte{0xC0, 0xCE, 0xB0, 0xA3, 0xC1, 0xB7})
		if err != nil {
			panic(err)
		}

		_ = decodedFolderA
		if strings.Contains(n.Value, "data/sprite") {
			for _, e := range grfFile.GetEntries(n.Value) {
				if strings.Contains(e.Name, ".spr") {
					nodeEntries = append(nodeEntries, e.Name)
				}
			}

			var decodedDirName []byte
			var err error
			if decodedDirName, err = cp949.From([]byte(n.Value)); err != nil {
				panic(err)
			}

			if len(nodeEntries) > 0 {
				node := g.TreeNode(fmt.Sprintf("%s (%d)", decodedDirName, len(nodeEntries)))
				selectableNodes = g.RangeBuilder("selectableNodes", nodeEntries, func(i int, v interface{}) g.Widget {
					var decodedBytes []byte
					var err error
					if decodedBytes, err = cp949.From([]byte(v.(string))); err != nil {
						panic(err)
					}
					return g.Selectable(string(decodedBytes)).OnClick(func() {
						onClickEntry(v.(string))
					})
				})

				node.Layout(selectableNodes...)
				nodes = append(nodes, node)
			}
		}
	})

	tree := g.RangeBuilder("entries", nodes, func(i int, v interface{}) g.Widget {
		return v.(g.Widget)
	})

	return g.Layout{tree}
}

func loadImage(name string) *g.Texture {
	if grfFile == nil {
		return nil
	}

	sprFile, _ := spr.Load(currentEntry.Data)
	img := sprFile.ImageAt(0)

	g.NewTextureFromRgba(img.RGBA, func(spriteTexture *g.Texture) {
		imageWidget = g.Image(spriteTexture).Size(float32(img.Bounds().Max.X), float32(img.Bounds().Max.Y))
		loadedImageName = name
	})

	return nil
}

func noErr(err error) bool {
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return true
}
