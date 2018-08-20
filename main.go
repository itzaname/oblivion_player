package main

import (
	"github/itzaname/oblivion_player/manager"

	"github.com/veandco/go-sdl2/sdl"
	"os"
	"runtime"
)

var renderer *sdl.Renderer
var window *sdl.Window

func init() {
	runtime.LockOSThread()
}

func main() {
	os.Mkdir(os.TempDir()+string(os.PathSeparator)+"oblivion", os.ModePerm)

	instance, err := manager.New()
	if err != nil {
		panic(err)
	}

	instance.Run()
}
