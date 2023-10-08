package main

import (
	_ "embed"
	"flag"
	"io"
	"os"
	"path"
	"strings"

	"github.com/alacrity-engine/core/math/geometry"
	"github.com/golang-collections/collections/queue"
	lua "github.com/yuin/gopher-lua"
	bolt "go.etcd.io/bbolt"
	luar "layeh.com/gopher-luar"
)

const (
	sceneBucketName   = "scene"
	sceneBucketPrefix = "__scene"
)

var (
	projectPath      string
	resourceFilePath string
	//go:embed funcs.lua
	luaFuncs string
)

func parseFlags() {
	flag.StringVar(&projectPath, "project", ".",
		"Path to the project to pack animations for.")
	flag.StringVar(&resourceFilePath, "out", "./stage.res",
		"Resource file to store animations and spritesheets.")

	flag.Parse()
}

func main() {
	parseFlags()

	state := lua.NewState()
	defer state.Close()

	// Import Lua functions.
	err := state.DoString(luaFuncs)
	handleError(err)

	// Open the resource file.
	resourceFile, err := bolt.Open(resourceFilePath, 0666, nil)
	handleError(err)
	defer resourceFile.Close()

	err = resourceFile.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(
			[]byte(sceneBucketName))

		if err != nil {
			return err
		}

		return nil
	})
	handleError(err)

	// Create constants.
	state.SetGlobal("DrawModeStatic", luar.New(state, 0x88E4))
	state.SetGlobal("DrawModeDynamic", luar.New(state, 0x88E8))
	state.SetGlobal("DrawModeStream", luar.New(state, 0x88E0))

	// Import Go functions.
	state.SetGlobal("vec", luar.New(state, geometry.V))
	state.SetGlobal("rect", luar.New(state, geometry.R))
	state.SetGlobal("list", luar.New(state, list))
	state.SetGlobal("colorMask", luar.New(state, colorMask))
	state.SetGlobal("colorRepeat4", luar.New(state, colorRepeat4))

	// Resource preload functions.
	state.SetGlobal("preloadAnimations", luar.New(state,
		preloadAnimations(resourceFile, handleError)))
	state.SetGlobal("preloadAudios", luar.New(state,
		preloadAudios(resourceFile, handleError)))
	state.SetGlobal("preloadTextures", luar.New(state,
		preloadTextures(resourceFile, handleError)))
	state.SetGlobal("preloadFonts", luar.New(state,
		preloadFonts(resourceFile, handleError)))
	state.SetGlobal("preloadPictures", luar.New(state,
		preloadPictures(resourceFile, handleError)))
	state.SetGlobal("preloadShaders", luar.New(state,
		preloadShaders(resourceFile, handleError)))
	state.SetGlobal("preloadShaderPrograms", luar.New(state,
		preloadShaderPrograms(resourceFile, handleError)))

	// Graphics creation functions.
	state.SetGlobal("createCanvases", luar.New(state,
		createCanvases(resourceFile, handleError)))
	state.SetGlobal("createBatches", luar.New(state,
		createBatches(resourceFile, handleError)))

	// Game object creation functions.
	state.SetGlobal("createGameObjects", luar.New(state,
		createGameObjects(resourceFile, handleError)))
	state.SetGlobal("storePrefabs", luar.New(state,
		storePrefabs(resourceFile, handleError)))

	entries, err := os.ReadDir(projectPath)
	handleError(err)

	traverseQueue := queue.New()

	if len(entries) <= 0 {
		return
	}

	for _, entry := range entries {
		traverseQueue.Enqueue(FileTracker{
			EntryPath: ".",
			Entry:     entry,
		})
	}

	for traverseQueue.Len() > 0 {
		fsEntry := traverseQueue.Dequeue().(FileTracker)

		if fsEntry.Entry.IsDir() {
			entries, err = os.ReadDir(path.Join(fsEntry.EntryPath, fsEntry.Entry.Name()))
			handleError(err)

			for _, entry := range entries {
				traverseQueue.Enqueue(FileTracker{
					EntryPath: path.Join(fsEntry.EntryPath, fsEntry.Entry.Name()),
					Entry:     entry,
				})
			}

			continue
		}

		if !strings.HasSuffix(fsEntry.Entry.Name(), ".main.lua") {
			continue
		}

		// Execute the main script.
		file, err := os.Open(projectPath)
		handleError(err)
		data, err := io.ReadAll(file)
		handleError(err)
		err = file.Close()
		handleError(err)
		err = state.DoString(string(data))
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
