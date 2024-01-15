package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/atotto/clipboard"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sqweek/dialog"

	"github.com/project-midgard/midgarts/assets"
	"github.com/project-midgard/midgarts/internal/fileformat/act"
	"github.com/project-midgard/midgarts/internal/fileformat/grf"
	"github.com/project-midgard/midgarts/internal/fileformat/spr"
)

var (
	GrfFilePath = os.Getenv("GRF_FILE_PATH")
)

type App struct {
	grfFile          *grf.File
	grfPath          string
	imageWidget      *g.ImageWidget
	fileInfoWidget   g.Widget
	loadedImageName  string
	currentEntry     *grf.Entry
	currentEntryName string
	splitSize        float32
	font             []byte
	openFileSelector bool
	filter           string
	filterByKorean   bool
}

func main() {
	app := &App{
		splitSize:   500,
		imageWidget: &g.ImageWidget{},
	}

	configureLogger()
	if err := app.loadInitialData(); err != nil {
		log.Fatal().Err(err).Send()
	}

	wnd := g.NewMasterWindow("GRF Explorer", 800, 600, 0)
	g.Context.FontAtlas.SetDefaultFontFromBytes(app.font, 16)
	wnd.Run(app.run)
}

func configureLogger() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func (app *App) loadInitialData() error {
	// loads an embedded font
	app.font = assets.FreeSans

	// respects grf env var (if present)
	if GrfFilePath != "" {
		var err error
		app.grfFile, err = grf.Load(GrfFilePath)
		if err != nil {
			log.Warn().Err(err).Msg("Error loading GRF file")

			// nil fail, let the user manually select their grf
			return nil
		}
	}

	return nil
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
		),
		g.Row(
			g.Label("Filter:"),
			g.InputText(&app.filter).Size(200),
			g.Checkbox("Korean", &app.filterByKorean),
		),
		g.SplitLayout(g.DirectionVertical, &app.splitSize, app.buildEntryTreeNodes(), app.fileInfoLayout()),
		app.fileSelectorModal(),
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
			filterValue := e.Name.String()
			if app.filterByKorean {
				filterValue = e.Name.Korean()
			}
			if strings.Contains(e.Name.String(), ".spr") && strings.Contains(filterValue, app.filter) {
				nodeEntries = append(nodeEntries, e.Name)
			}
		}

		if len(nodeEntries) > 0 {
			node := g.TreeNode(fmt.Sprintf("%s (%d)", n.Value, len(nodeEntries)))
			selectableNodes = g.RangeBuilder("selectableNodes", nodeEntries, func(i int, v interface{}) g.Widget {
				entryPath := v.(*grf.Path)
				return g.Style().
					SetColor(g.StyleColorText, color.RGBA{203, 213, 255, 255}).
					To(
						g.Selectable(entryPath.String()).OnClick(func() {
							app.onClickEntry(entryPath.String())
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
	app.currentEntryName = entryName

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

	_, koreanFileName := filepath.Split(app.currentEntry.Name.Korean())

	app.fileInfoWidget = g.Layout{
		g.Table().
			Columns(
				g.TableColumn("File Info"),
				g.TableColumn(""),
			).
			Rows(
				g.TableRow(g.Label("Width").Wrapped(true), g.Label(fmt.Sprintf("%d", sprFile.Frames[0].Width))),
				g.TableRow(g.Label("Height").Wrapped(true), g.Label(fmt.Sprintf("%d", sprFile.Frames[0].Height))),
				g.TableRow(g.Label("Korean").Wrapped(true), g.Row(g.Label(koreanFileName), g.Button("[C]").OnClick(func() {
					clipboard.WriteAll(koreanFileName)
				}))),
				g.TableRow(
					g.Button("Copy Path").OnClick(func() {
						clipboard.WriteAll(app.currentEntryName)
					}),
					g.Button("Copy Korean Path").OnClick(func() {
						clipboard.WriteAll(app.currentEntry.Name.Korean())
					}),
				),
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

func errorDialog(title, content string) {
	dialog.Message("%v", content).Title(title).Error()
}
