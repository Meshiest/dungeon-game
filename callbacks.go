package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"os"
)

var keys map[glfw.Key]bool

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyEscape:
			glfw.Terminate()
			os.Exit(0)
		}
	}
	if action != glfw.Repeat {
		keys[key] = action == glfw.Press
	}
}

func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {

}

func SizeCallback(w *glfw.Window, width, height int) {
	screenWidth = width
	screenHeight = height
	projection = mgl32.Perspective(mgl32.DegToRad(float32(fov)), float32(screenWidth)/float32(screenHeight), 0.01, 20.0)
	gl.Viewport(0, 0, int32(width), int32(height))
}
