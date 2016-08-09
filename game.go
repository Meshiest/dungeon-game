package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/meshiest/dungeon-game/platform"
	"github.com/meshiest/go-dungeon/dungeon"
)

func FloorTile(xInt int, zInt int) []float32 {
	x := float32(xInt)
	z := float32(zInt)
	return []float32{
		-0.5 + x, 0, -0.5 + z, 0, 0,
		0.5 + x, 0, -0.5 + z, 1, 0,
		-0.5 + x, 0, 0.5 + z, 0, 1,
		0.5 + x, 0, -0.5 + z, 1, 0,
		0.5 + x, 0, 0.5 + z, 1, 1,
		-0.5 + x, 0, 0.5 + z, 0, 1,
	}
}

func WallTile(xInt int, zInt int, dir bool, offset mgl32.Vec2) []float32 {
	x := float32(xInt)
	z := float32(zInt)
	var xAxis, zAxis float32
	if dir {
		xAxis = 1
		zAxis = 0
	} else {
		xAxis = 0
		zAxis = 1
	}
	return []float32{
		-0.5*xAxis + x + offset.X(), 0, -0.5*zAxis + z + offset.Y(), 1, 0,
		0.5*xAxis + x + offset.X(), 0, 0.5*zAxis + z + offset.Y(), 0, 0,
		-0.5*xAxis + x + offset.X(), 1, -0.5*zAxis + z + offset.Y(), 1, 1,
		0.5*xAxis + x + offset.X(), 0, 0.5*zAxis + z + offset.Y(), 0, 0,
		0.5*xAxis + x + offset.X(), 1, 0.5*zAxis + z + offset.Y(), 0, 1,
		-0.5*xAxis + x + offset.X(), 1, -0.5*zAxis + z + offset.Y(), 1, 1,
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var screenWidth = 800
var screenHeight = 600

var vertexArray []float32
var numFloorTiles, numWallTiles int
var fov float64
var projection mgl32.Mat4

const (
	FOG_DISTANCE float32 = 8.0
)

func init() {
	runtime.LockOSThread()
}

func main() {
	keys = map[glfw.Key]bool{}
	rand.Seed(time.Now().Unix())

	player := &Player{
		Yaw:    0,
		Pitch:  0,
		X:      0,
		Y:      0,
		Speed:  2,
		Size:   0.7,
		Health: 1,
	}
	fov = 90.0
	fmt.Println("Generating Dungeon...")
	dungeon := dungeon.NewDungeon(50, 200)
	fmt.Println("Generated!")
	room := dungeon.Rooms[0]
	player.X = float64(room.Y + room.Height/2)
	player.Y = float64(room.X + room.Width/2)
	dungeon.Print()

	enemies := []*Enemy{}
	bloods := []*Blood{}

	vertexArray = []float32{
		// Enemy Sprite
		-0.3, 0.6, 0, 1, 0,
		0.3, 0.6, 0, 0, 0,
		-0.3, 0.0, 0, 1, 1,
		0.3, 0.6, 0, 0, 0,
		0.3, 0.0, 0, 0, 1,
		-0.3, 0.0, 0, 1, 1,

		// Blood Sprite
		-0.5, 0, -0.5, 0, 0,
		0.5, 0, -0.5, 1, 0,
		-0.5, 0, 0.5, 0, 1,
		0.5, 0, -0.5, 1, 0,
		0.5, 0, 0.5, 1, 1,
		-0.5, 0, 0.5, 0, 1,
	}

	numFloorTiles = 0
	numWallTiles = 0
	for y, row := range dungeon.Grid {
		for x, col := range row {
			if col == 1 {
				vertexArray = append(vertexArray, FloorTile(x, y)...)
				numFloorTiles++
				if rand.Int()%10 == 0 {
					enemies = append(enemies, &Enemy{
						X:    float64(x),
						Y:    float64(y),
						Size: 0.5,
						DPS:  0.1,
					})
				}
			}
		}
	}

	for y, row := range dungeon.Grid {
		for x, col := range row {
			if col == 1 {
				if y > 0 && dungeon.Grid[y-1][x] == 0 || y == 0 {
					vertexArray = append(vertexArray, WallTile(x, y, true, mgl32.Vec2{0, -0.5})...)
					numWallTiles++
				}
				if y < len(dungeon.Grid)-1 && dungeon.Grid[y+1][x] == 0 || y == len(dungeon.Grid)-1 {
					vertexArray = append(vertexArray, WallTile(x, y, true, mgl32.Vec2{0, 0.5})...)
					numWallTiles++
				}
				if x > 0 && dungeon.Grid[y][x-1] == 0 || x == 0 {
					vertexArray = append(vertexArray, WallTile(x, y, false, mgl32.Vec2{-0.5, 0})...)
					numWallTiles++
				}
				if x < len(row)-1 && dungeon.Grid[y][x+1] == 0 || x == len(row)-1 {
					vertexArray = append(vertexArray, WallTile(x, y, false, mgl32.Vec2{0.5, 0})...)
					numWallTiles++
				}
			}
		}
	}

	err := glfw.Init()
	check(err)

	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	platform.WHint() //platform-specific window hinting

	window, err := glfw.CreateWindow(screenWidth, screenHeight, "Dungeon", nil, nil)
	check(err)

	window.MakeContextCurrent()

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetKeyCallback(glfw.KeyCallback(KeyCallback))
	window.SetSizeCallback(glfw.SizeCallback(SizeCallback))
	window.SetMouseButtonCallback(glfw.MouseButtonCallback(MouseButtonCallback))

	fmt.Println("Initializing GL")
	err = gl.Init()
	check(err)

	fmt.Println("Loading Shaders")
	vertexShader, err := ioutil.ReadFile("shaders/shader.vert")
	check(err)

	fragmentShader, err := ioutil.ReadFile("shaders/shader.frag")
	check(err)

	program, err := newProgram(string(vertexShader)+"\x00", string(fragmentShader)+"\x00")
	check(err)

	gl.UseProgram(program)

	projection = mgl32.Perspective(mgl32.DegToRad(float32(fov)), float32(screenWidth)/float32(screenHeight), 0.001, 20.0)

	//camera := mgl32.LookAtV(mgl32.Vec3{3, 1, 0}, mgl32.Vec3{2, 1, 0}, mgl32.Vec3{0, 1, 0})
	camera := mgl32.LookAtV(mgl32.Vec3{0, 1, 0}, mgl32.Vec3{1, 1, 0}, mgl32.Vec3{0, 1, 0})

	viewProjUniform := gl.GetUniformLocation(program, gl.Str("viewProj\x00"))
	viewProj := camera.Mul4(projection)
	gl.UniformMatrix4fv(viewProjUniform, 1, false, &viewProj[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	fogDistUniform := gl.GetUniformLocation(program, gl.Str("fogDist\x00"))
	gl.Uniform1f(fogDistUniform, FOG_DISTANCE)

	fmt.Println("Loading Textures")
	floorTexture, err := newTexture("textures/floor.png")
	check(err)

	wallTexture, err := newTexture("textures/wall.png")
	check(err)

	enemyTexture, err := newTexture("textures/monster.png")
	check(err)

	bloodTexture, err := newTexture("textures/blood.png")
	check(err)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexArray)*4, gl.Ptr(vertexArray), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0, 0, 0, 1)

	previousTime := glfw.GetTime()
	lastFPS := previousTime
	fps := 0

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		time := glfw.GetTime()
		delta := time - previousTime
		previousTime = time
		gl.UseProgram(program)

		fps++
		if time-lastFPS > 1 {
			fmt.Println("FPS is ", fps)
			lastFPS = time
			fps = 0
		}

		mouseSensitivity := 0.75
		mouseX, mouseY := window.GetCursorPos()
		window.SetCursorPos(float64(screenWidth/2), float64(screenHeight/2))
		ratio := float64(screenWidth) / float64(screenHeight)
		mouseDeltaX := float64(screenWidth/2) - mouseX
		mouseDeltaY := float64(screenHeight/2) - mouseY
		player.Yaw -= mouseSensitivity * delta * mouseDeltaX
		player.Pitch += mouseSensitivity * delta * mouseDeltaY * ratio
		//fmt.Println(yaw/math.Pi*360)

		if player.Pitch > math.Pi/2 {
			player.Pitch = math.Pi / 2
		}
		if player.Pitch < -math.Pi/2 {
			player.Pitch = -math.Pi / 2
		}

		direction := mgl32.Vec2{0, 0}
		if keys[glfw.KeyW] {
			direction = direction.Add(mgl32.Vec2{0, 1})
		}
		if keys[glfw.KeyS] {
			direction = direction.Add(mgl32.Vec2{0, -1})
		}
		if keys[glfw.KeyA] {
			direction = direction.Add(mgl32.Vec2{1, 0})
		}
		if keys[glfw.KeyD] {
			direction = direction.Add(mgl32.Vec2{-1, 0})
		}
		direction = direction.Normalize()
		if direction.Len() > 0 {

			boost := 1.0
			if keys[glfw.KeyLeftShift] {
				boost = 2.0
			}

			cos := float32(math.Cos(player.Yaw - math.Pi/2))
			sin := float32(math.Sin(player.Yaw - math.Pi/2))
			rotated := mgl32.Vec2{
				direction.X()*cos - sin*direction.Y(),
				direction.X()*sin + cos*direction.Y(),
			}
			player.X += float64(rotated.X()) * player.Speed * delta * boost
			player.Y += float64(rotated.Y()) * player.Speed * delta * boost

		}
		player.CollideWithDungeon(dungeon)

		camera = mgl32.LookAt(
			float32(player.X), float32(0.25), float32(player.Y),
			float32(player.X+math.Cos(player.Yaw)), float32(0.25+math.Sin(player.Pitch)), float32(player.Y+math.Sin(player.Yaw)),
			0, 1, 0,
		)

		viewProj = projection.Mul4(camera)
		gl.UniformMatrix4fv(viewProjUniform, 1, false, &viewProj[0])

		model = mgl32.Ident4()
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(vao)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, floorTexture)
		gl.DrawArrays(gl.TRIANGLES, int32(12), int32(numFloorTiles*3*2))
		gl.BindTexture(gl.TEXTURE_2D, wallTexture)
		gl.DrawArrays(gl.TRIANGLES, int32(12+numFloorTiles*3*2), int32(numWallTiles*3*2))

		gl.BindTexture(gl.TEXTURE_2D, enemyTexture)
		closest := -1
		dist := float32(3)
		for i, enemy := range enemies {

			from := mgl32.Vec2{float32(player.X - enemy.X), float32(player.Y - enemy.Y)}
			if from.Len() < dist {
				dist = from.Len()
				closest = i
			}
			from = from.Normalize()

			enemy.X += float64(from.X()) * delta
			enemy.Y += float64(from.Y()) * delta
			if enemy.CollideWithPlayer(player) {
				player.Health -= enemy.DPS * delta
				//fmt.Println(player.Health)
			}
			for _, other := range enemies {
				if enemy != other {
					enemy.CollideWithEnemy(other)
				}
			}
			enemy.CollideWithDungeon(dungeon)
			model = mgl32.Translate3D(float32(enemy.X), 0.1, float32(enemy.Y))
			model = model.Mul4(mgl32.HomogRotate3DY(float32(math.Pi/2 - math.Atan2(enemy.Y-player.Y, enemy.X-player.X))))
			gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 6)
		}

		if time-player.LastAttack > 0.5 && keys[glfw.KeySpace] {
			if closest >= 0 {
				player.LastAttack = time
				bloods = append(bloods, NewBlood(enemies[closest].X, enemies[closest].Y))
				if closest == len(enemies) {
					enemies = enemies[:closest]
				} else {
					enemies = append(enemies[:closest], enemies[closest+1:]...)
				}
			}
		}

		gl.BindTexture(gl.TEXTURE_2D, bloodTexture)

		for _, blood := range bloods {
			model = mgl32.Translate3D(float32(blood.X), 0.01, float32(blood.Y))
			model = model.Mul4(mgl32.HomogRotate3DY(float32(blood.Angle)))
			gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
			gl.DrawArrays(gl.TRIANGLES, 6, 6)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
