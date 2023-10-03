package main

import (
	"flag"
	"io"
	"os"

	"github.com/alacrity-engine/core/geometry"
	"github.com/alacrity-engine/core/render"
	lua "github.com/yuin/gopher-lua"
	bolt "go.etcd.io/bbolt"
	luar "layeh.com/gopher-luar"
)

const (
	sceneBucketName = "scene"
)

var (
	mainScriptFilePath string
	resourceFilePath   string
)

func parseFlags() {
	flag.StringVar(&mainScriptFilePath, "script", "./main.lua",
		"Main Lua script to construct a scene")
	flag.StringVar(&resourceFilePath, "out", "./stage.res",
		"Resource file to store animations and spritesheets.")

	flag.Parse()
}

func main() {
	parseFlags()

	state := lua.NewState()
	defer state.Close()

	// Import Lua functions.
	file, err := os.Open("funcs.lua")
	handleError(err)
	data, err := io.ReadAll(file)
	handleError(err)
	err = file.Close()
	handleError(err)
	err = state.DoString(string(data))
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

	// Import Go functions.
	state.SetGlobal("vec", luar.New(state, geometry.V))
	state.SetGlobal("ortho2DStandard",
		luar.New(state, render.Ortho2DStandard))
	state.SetGlobal("list", luar.New(state, list))

	// Resource preload functions.
	state.SetGlobal("preloadAnimations", luar.New(state,
		preloadAnimations(resourceFile, handleError)))
	state.SetGlobal("preloadAudios", luar.New(state,
		preloadAudios(resourceFile, handleError)))
	state.SetGlobal("preloadTextures", luar.New(state,
		preloadTextures(resourceFile, handleError)))
	state.SetGlobal("preloadFonts", luar.New(state,
		preloadFonts(resourceFile, handleError)))
	state.SetGlobal("preloadSPictures", luar.New(state,
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

	// Execute the main script.
	file, err = os.Open(mainScriptFilePath)
	handleError(err)
	data, err = io.ReadAll(file)
	handleError(err)
	err = file.Close()
	handleError(err)
	err = state.DoString(string(data))
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
