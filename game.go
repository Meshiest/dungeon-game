package main

import (
    "github.com/go-gl/glfw/v3.2/glfw"
    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/mathgl/mgl32"
    "image"
    "image/draw"
    _ "image/png"
    "fmt"
    "os"
    "strings"
    "github.com/meshiest/go-dungeon/dungeon"
    "math"
    "runtime"
    "io/ioutil"
    "math/rand"
    "time"
)

var keys map[glfw.Key]bool

func FloorTile (xInt int, zInt int) ([]float32) {
  x := float32(xInt)
  z := float32(zInt)
  return []float32{
    -0.5 + x, 0, -0.5 + z, 0, 0,
     0.5 + x, 0, -0.5 + z, 1, 0,
    -0.5 + x, 0,  0.5 + z, 0, 1,
     0.5 + x, 0, -0.5 + z, 1, 0,
     0.5 + x, 0,  0.5 + z, 1, 1,
    -0.5 + x, 0,  0.5 + z, 0, 1,
  }
}

func WallTile (xInt int, zInt int, dir bool, offset mgl32.Vec2) ([]float32) {
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
    -0.5 * xAxis + x + offset.X(), 0, -0.5 * zAxis + z + offset.Y(), 1, 0,
     0.5 * xAxis + x + offset.X(), 0,  0.5 * zAxis + z + offset.Y(), 0, 0,
    -0.5 * xAxis + x + offset.X(), 1, -0.5 * zAxis + z + offset.Y(), 1, 1,
     0.5 * xAxis + x + offset.X(), 0,  0.5 * zAxis + z + offset.Y(), 0, 0,
     0.5 * xAxis + x + offset.X(), 1,  0.5 * zAxis + z + offset.Y(), 0, 1,
    -0.5 * xAxis + x + offset.X(), 1, -0.5 * zAxis + z + offset.Y(), 1, 1,
  }
}

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

type Enemy struct {
  X, Y float64
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}
 
var screenWidth = 800
var screenHeight = 600
func SizeCallback(w *glfw.Window, width, height int) {
  screenWidth = width
  screenHeight = height
  projection = mgl32.Perspective(mgl32.DegToRad(float32(fov)), float32(screenWidth)/float32(screenHeight), 0.1, 20.0)
  gl.Viewport(0, 0, int32(width), int32(height))
}

var vertexArray []float32
var numFloorTiles, numWallTiles int
var yaw, pitch, fov, positionX, positionY float64
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

  yaw = 0
  pitch = 0
  positionX = 0
  positionY = 0
  fov = 90.0
  fmt.Println("Generating Dungeon...")
  dungeon := dungeon.NewDungeon(50, 200)
  fmt.Println("Generated!")
  room := dungeon.Rooms[0]
  positionX = float64(room.Y + room.Height/2)
  positionY = float64(room.X + room.Width/2)
  //dungeon.Print()

  enemies := []Enemy{}

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
    -0.5, 0,  0.5, 0, 1,
     0.5, 0, -0.5, 1, 0,
     0.5, 0,  0.5, 1, 1,
    -0.5, 0,  0.5, 0, 1,
  }

  numFloorTiles = 0
  numWallTiles = 0
  for y, row := range(dungeon.Grid) {
    for x, col := range(row) {
      if col == 1 {
        vertexArray = append(vertexArray, FloorTile(x, y)...)
        numFloorTiles ++
        if rand.Int() % 10 == 0 {
          enemies = append(enemies, Enemy{X: float64(x), Y: float64(y)})
        }
      }
    }
  }

  for y, row := range(dungeon.Grid) {
    for x, col := range(row) {
      if col == 1 {
        if y > 0 && dungeon.Grid[y-1][x] == 0 || y == 0 {
          vertexArray = append(vertexArray, WallTile(x, y, true, mgl32.Vec2{0, -0.5})...)
          numWallTiles ++
        }
        if y < len(dungeon.Grid) - 1 && dungeon.Grid[y+1][x] == 0 || y == len(dungeon.Grid) - 1 {
          vertexArray = append(vertexArray, WallTile(x, y, true, mgl32.Vec2{0, 0.5})...)
          numWallTiles ++
        }
        if x > 0 && dungeon.Grid[y][x-1] == 0 || x == 0 {
          vertexArray = append(vertexArray, WallTile(x, y, false, mgl32.Vec2{-0.5, 0})...)
          numWallTiles ++
        }
        if x < len(row) - 1 && dungeon.Grid[y][x+1] == 0 || x == len(row) - 1 {
          vertexArray = append(vertexArray, WallTile(x, y, false, mgl32.Vec2{0.5, 0})...)
          numWallTiles ++
        }
      }
    }
  }

  err := glfw.Init()
  check(err)

  defer glfw.Terminate()

  glfw.WindowHint(glfw.Resizable, glfw.True)

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

  program, err := newProgram(string(vertexShader) + "\x00", string(fragmentShader) + "\x00")
  check(err)

  gl.UseProgram(program)

  projection = mgl32.Perspective(mgl32.DegToRad(float32(fov)), float32(screenWidth)/float32(screenHeight), 0.01, 20.0)
  projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
  gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

  //camera := mgl32.LookAtV(mgl32.Vec3{3, 1, 0}, mgl32.Vec3{2, 1, 0}, mgl32.Vec3{0, 1, 0})
  camera := mgl32.LookAtV(mgl32.Vec3{0, 1, 0}, mgl32.Vec3{1, 1, 0}, mgl32.Vec3{0, 1, 0})
  cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
  gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

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
  gl.BufferData(gl.ARRAY_BUFFER, len(vertexArray) * 4, gl.Ptr(vertexArray), gl.STATIC_DRAW)

  vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
  gl.EnableVertexAttribArray(vertAttrib)
  gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

  texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
  gl.EnableVertexAttribArray(texCoordAttrib)
  gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

  gl.Enable(gl.DEPTH_TEST)
  gl.DepthFunc(gl.LESS)
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

    fps ++
    if time - lastFPS > 1 {
      fmt.Println("FPS is ",fps)
      lastFPS = time
      fps = 0
    }


    mouseSensitivity := 0.75
    mouseX, mouseY := window.GetCursorPos()
    window.SetCursorPos(float64(screenWidth/2), float64(screenHeight/2))
    ratio := float64(screenWidth)/float64(screenHeight)
    mouseDeltaX := float64(screenWidth/2) - mouseX
    mouseDeltaY := float64(screenHeight/2) - mouseY
    yaw -= mouseSensitivity * delta * mouseDeltaX
    pitch += mouseSensitivity * delta * mouseDeltaY * ratio
    //fmt.Println(yaw/math.Pi*360)

    if pitch > math.Pi/2 {
      pitch = math.Pi/2
    }
    if pitch < -math.Pi/2 {
      pitch = -math.Pi/2
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
      speed := 2.0
      
      if keys[glfw.KeyLeftShift] {
        speed *= 2
      }

      cos := float32(math.Cos(yaw - math.Pi/2))
      sin := float32(math.Sin(yaw - math.Pi/2))
      rotated := mgl32.Vec2{
        direction.X() * cos - sin * direction.Y(),
        direction.X() * sin + cos * direction.Y(),
      }
      positionX += float64(rotated.X()) * speed * delta
      positionY += float64(rotated.Y()) * speed * delta
      x := int(math.Floor(positionX+0.5))
      y := int(math.Floor(positionY+0.5))

      if x >= 0 && y >= 0 && x < len(dungeon.Grid) && y < len(dungeon.Grid) {
        if (y == 0 || dungeon.Grid[y-1][x] == 0) && positionY - float64(y) < -0.4 {
          positionY = float64(y)-0.4
        }
        if (y == len(dungeon.Grid) - 1 || dungeon.Grid[y+1][x] == 0) && positionY - float64(y) > 0.4 {
          positionY = float64(y)+0.4
        }
        if (x == 0 || dungeon.Grid[y][x-1] == 0) && positionX - float64(x) < -0.4 {
          positionX = float64(x)-0.4
        }
        if (x == len(dungeon.Grid) - 1 || dungeon.Grid[y][x+1] == 0) && positionX - float64(x) > 0.4 {
          positionX = float64(x)+0.4
        }
      }
      //fmt.Println(x, y, dungeon.Grid[y][x])
    }

    camera = mgl32.LookAt(
        float32(positionX), float32(0.25), float32(positionY),
        float32(positionX + math.Cos(yaw)), float32(0.25 + math.Sin(pitch)), float32(positionY + math.Sin(yaw)),
        0, 1, 0,
      )
    gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

    model = mgl32.Ident4()
    gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

    gl.BindVertexArray(vao)
    gl.ActiveTexture(gl.TEXTURE0)
    gl.BindTexture(gl.TEXTURE_2D, floorTexture)
    gl.DrawArrays(gl.TRIANGLES, int32(12), int32(numFloorTiles*3*2))
    gl.BindTexture(gl.TEXTURE_2D, wallTexture)
    gl.DrawArrays(gl.TRIANGLES, int32(12 + numFloorTiles*3*2), int32(numWallTiles*3*2))


    gl.BindTexture(gl.TEXTURE_2D, enemyTexture)
    for _, enemy := range(enemies) {
      model = mgl32.Translate3D(float32(enemy.X), 0.1, float32(enemy.Y))
      model = model.Mul4(mgl32.HomogRotate3DY(float32(math.Pi/2-math.Atan2(enemy.Y-positionY, enemy.X-positionX))))
      //model = mgl32.LookAt(float32(enemy.X), 0.1, float32(enemy.Y), float32(positionX), 0.1, float32(positionY), 0, 1, 0)
      gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
      gl.DrawArrays(gl.TRIANGLES, 0, 6)
    }

    gl.BindTexture(gl.TEXTURE_2D, bloodTexture)

    window.SwapBuffers()
    glfw.PollEvents()
  }
}

func newTexture(file string) (uint32, error) {
  imgFile, err := os.Open(file)
  if err != nil {
    return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
  }
  img, _, err := image.Decode(imgFile)
  if err != nil {
    return 0, err
  }

  rgba := image.NewRGBA(img.Bounds())
  if rgba.Stride != rgba.Rect.Size().X*4 {
    return 0, fmt.Errorf("unsupported stride")
  }
  draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

  var texture uint32
  gl.GenTextures(1, &texture)
  gl.ActiveTexture(gl.TEXTURE0)
  gl.BindTexture(gl.TEXTURE_2D, texture)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
  gl.TexImage2D(
    gl.TEXTURE_2D,
    0,
    gl.RGBA,
    int32(rgba.Rect.Size().X),
    int32(rgba.Rect.Size().Y),
    0,
    gl.RGBA,
    gl.UNSIGNED_BYTE,
    gl.Ptr(rgba.Pix))

  return texture, nil
}


func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
  vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
  if err != nil {
    return 0, err
  }

  fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
  if err != nil {
    return 0, err
  }

  program := gl.CreateProgram()

  gl.AttachShader(program, vertexShader)
  gl.AttachShader(program, fragmentShader)
  gl.LinkProgram(program)

  var status int32
  gl.GetProgramiv(program, gl.LINK_STATUS, &status)
  if status == gl.FALSE {
    var logLength int32
    gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

    log := strings.Repeat("\x00", int(logLength+1))
    gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

    return 0, fmt.Errorf("failed to link program: %v", log)
  }

  gl.DeleteShader(vertexShader)
  gl.DeleteShader(fragmentShader)

  return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
  shader := gl.CreateShader(shaderType)

  csources, free := gl.Strs(source)
  gl.ShaderSource(shader, 1, csources, nil)
  free()
  gl.CompileShader(shader)

  var status int32
  gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
  if status == gl.FALSE {
    var logLength int32
    gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

    log := strings.Repeat("\x00", int(logLength+1))
    gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

    return 0, fmt.Errorf("failed to compile %v: %v", source, log)
  }

  return shader, nil
}