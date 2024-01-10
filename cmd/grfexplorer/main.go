package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sqweek/dialog"
	"github.com/suapapa/go_hangul/encoding/cp949"
	"golang.org/x/text/encoding/charmap"

	"github.com/project-midgard/midgarts/assets"
	"github.com/project-midgard/midgarts/internal/fileformat/act"
	"github.com/project-midgard/midgarts/internal/fileformat/grf"
	"github.com/project-midgard/midgarts/internal/fileformat/spr"
)

type SupportedEncodings string

const (
	EncodingWindows1252 SupportedEncodings = "Windows1252"
	EncodingCP949       SupportedEncodings = "CP949"
)

type App struct {
	grfFile          *grf.File
	grfPath          string
	imageWidget      *g.ImageWidget
	fileInfoWidget   g.Widget
	loadedImageName  string
	currentEntry     *grf.Entry
	splitSize        float32
	font             []byte
	currentEncoding  SupportedEncodings
	openFileSelector bool
}

func main() {
	app := &App{
		splitSize:       500,
		currentEncoding: EncodingWindows1252,
		imageWidget:     &g.ImageWidget{},
	}

	configureLogger()
	if err := app.loadInitialData(); err != nil {
		log.Fatal().Err(err).Send()
	}

	wnd := g.NewMasterWindow("GRF Explorer", 800, 600, g.MasterWindowFlagsNotResizable)
	g.Context.FontAtlas.SetDefaultFontFromBytes(app.font, 16)
	wnd.Run(app.run)
}

func configureLogger() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func (app *App) loadInitialData() (err error) {
	app.font = assets.FreeSans
	return err
}

func (app *App) run() {
	g.SingleWindowWithMenuBar().Layout(
		g.MenuBar().Layout(
			g.Menu("File").Layout(
				g.MenuItem("Open GRF").OnClick(app.openGRFSelector),
				g.MenuItem("Quit").OnClick(func() {
					os.Exit(0)
				}),
			),
			app.fileEncodingMenu(),
		),
		g.SplitLayout(g.DirectionVertical, &app.splitSize, app.buildEntryTreeNodes(), app.fileInfoLayout()),
		app.fileSelectorModal(),
	)
}

func (app *App) fileEncodingMenu() *g.MenuWidget {
	return g.Menu("Settings").Layout(
		g.Menu("File Path Encoding").Layout(
			g.MenuItem(string(EncodingWindows1252)).OnClick(func() {
				app.currentEncoding = EncodingWindows1252
			}).Selected(app.currentEncoding == EncodingWindows1252),
			g.MenuItem(string(EncodingCP949)).OnClick(func() {
				app.currentEncoding = EncodingCP949
			}).Selected(app.currentEncoding == EncodingCP949),
		),
	)
}

func (app *App) openGRFSelector() {
	app.openFileSelector = true
}

func (app *App) fileSelectorModal() g.Widget {
	return g.Custom(func() {
		if app.openFileSelector {
			var err error
			defer func() {
				app.openFileSelector = false
			}()
			app.grfPath, err = dialog.File().Filter("", "grf").Load()
			if err != nil {
				log.Error().Err(err).Msg("Error opening GRF file")
				return
			}
			app.grfFile, err = grf.Load(app.grfPath)
			if err != nil {
				log.Error().Err(err).Msg("Error loading GRF file")
				errorDialog("⚠️  Loading GRF File", "Error loading GRF file: \n"+err.Error())

				return
			}
		}
	})
}

func (app *App) buildEntryTreeNodes() g.Layout {
	if app.grfFile == nil {
		return g.Layout{
			g.Style().SetStyle(g.StyleVarFramePadding, 10, 5).
				To(
					g.Button("Select a GRF file").OnClick(app.openGRFSelector),
				),
		}
	}
	entries := app.grfFile.GetEntryTree()

	var nodes []any

	entries.Traverse(entries.Root, func(n *grf.EntryTreeNode) {
		if !strings.Contains(n.Value, "data/sprite") {
			return
		}

		selectableNodes := make([]g.Widget, 0)
		nodeEntries := make([]any, 0)

		for _, e := range app.grfFile.GetEntries(n.Value) {
			if strings.Contains(e.Name, ".spr") {
				nodeEntries = append(nodeEntries, e.Name)
			}
		}

		var decodedDirName []byte
		var err error
		if decodedDirName, err = decode([]byte(n.Value), app.currentEncoding); err != nil {
			decodedDirName = []byte(fmt.Sprintf("⚠️ %s", n.Value))
		}

		if len(nodeEntries) > 0 {
			node := g.TreeNode(fmt.Sprintf("%s (%d)", decodedDirName, len(nodeEntries)))
			selectableNodes = g.RangeBuilder("selectableNodes", nodeEntries, func(i int, v interface{}) g.Widget {
				var decodedBytes []byte
				var err error
				if decodedBytes, err = decode([]byte(v.(string)), app.currentEncoding); err != nil {
					decodedBytes = []byte(fmt.Sprintf("⚠️  %s", v.(string)))
				}
				return g.Style().
					SetColor(g.StyleColorText, color.RGBA{203, 213, 255, 255}).
					To(
						g.Selectable(string(decodedBytes)).OnClick(func() {
							app.onClickEntry(v.(string))
						}),
					)

			})

			node.Layout(selectableNodes...)
			nodes = append(nodes, node)
		}
	})

	tree := g.RangeBuilder("entries", nodes, func(i int, v interface{}) g.Widget {
		return v.(g.Widget)
	})

	return g.Layout{tree}
}

func (app *App) onClickEntry(entryName string) {
	if strings.Contains(entryName, ".act") {
		var err error
		app.currentEntry, err = app.grfFile.GetEntry(entryName)
		if err != nil {
			log.Error().Err(err).Msg("Error getting .act file entry")
			return
		}

		actFile, err := act.Load(app.currentEntry.Data)
		if err != nil {
			log.Error().Err(err).Msg("Error loading .act file")
			return
		}
		log.Info().Msgf("actFile: %+v", actFile)
	}

	if strings.Contains(entryName, ".spr") {
		var err error
		app.currentEntry, err = app.grfFile.GetEntry(entryName)
		if err != nil {
			log.Error().Err(err).Msg("Error getting .spr file entry")
			return
		}

		app.loadImage(entryName)
		app.loadFileInfo()
	}
}

func (app *App) loadFileInfo() {
	sprFile, err := spr.Load(app.currentEntry.Data)
	if err != nil {
		log.Error().Err(err).Msg("Error loading .spr file for file info")
		return
	}

	app.fileInfoWidget = g.Layout{
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

func (app *App) loadImage(name string) {
	if app.grfFile == nil {
		return
	}

	sprFile, err := spr.Load(app.currentEntry.Data)
	if err != nil {
		log.Error().Err(err).Msg("Error loading image")
		return
	}

	img := sprFile.ImageAt(0)
	g.NewTextureFromRgba(img.RGBA, func(spriteTexture *g.Texture) {
		app.imageWidget = g.Image(spriteTexture).Size(float32(img.Bounds().Max.X), float32(img.Bounds().Max.Y))
		app.loadedImageName = name
	})
}

func (app *App) fileInfoLayout() g.Layout {
	return g.Layout{
		app.fileInfoWidget,
		app.imageWidget,
		g.Custom(func() {
			if g.IsItemActive() {
				app.loadImage(app.loadedImageName)
			}
		}),
	}
}

func decode(buf []byte, encoding SupportedEncodings) ([]byte, error) {
	switch encoding {
	case EncodingCP949:
		return cp949.From(buf)
	case EncodingWindows1252:
		return charmap.Windows1252.NewDecoder().Bytes(buf)
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

func errorDialog(title, content string) {
	dialog.Message("%v", content).Title(title).Error()
}
