package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1920)
	screenHeight := int32(1080)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - 3d camera first person")
	defer rl.CloseWindow()

	rl.SetConfigFlags(rl.FlagVsyncHint)
	// rl.SetConfigFlags(rl.FlagBorderlessWindowedMode)
	camera := rl.Camera3D{
		Position:   rl.NewVector3(4, 2, 4),
		Target:     rl.NewVector3(0, 1.8, 0),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       60,
		Projection: rl.CameraProjection(rl.CameraPerspective),
	}

	initialTarget := camera.Target
	initialDirection := rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
	spherePosition := initialTarget

	// rl.ToggleFullscreen()
	heights := make([]float32, 20)
	positions := make([]rl.Vector3, 20)
	colors := make([]rl.Color, 20)

	for i := 0; i < 20; i++ {
		heights[i] = float32(rand.Intn(12) + 1)
		positions[i] = rl.NewVector3(float32(rand.Intn(31)-15), heights[i]/2, float32(rand.Intn(31)-15))
		colors[i] = rl.NewColor(uint8(rand.Intn(236)+20), uint8(rand.Intn(46)+10), 30, 255)
	}

	rl.SetTargetFPS(120)
	var distance float32
	speed := 0.1

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)

		rl.SetMousePosition(int(screenWidth)/2, int(screenHeight)/2)
		rl.DisableCursor()

		rl.BeginDrawing()
		rl.ClearBackground(rl.SkyBlue)

		rl.BeginMode3D(camera)

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			distance = 0.0
			initialTarget = camera.Target
			initialDirection = rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
			spherePosition = initialTarget
		}

		distance += float32(speed)
		spherePosition = rl.Vector3Add(initialTarget, rl.Vector3Scale(initialDirection, distance))

		rl.DrawSphere(spherePosition, 0.1, rl.Red)
		rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(128, 128), rl.LightGray)

		for i := 0; i < 20; i++ {
			rl.DrawCube(positions[i], 2, heights[i], 2, colors[i])
			rl.DrawCubeWires(positions[i], 2, heights[i], 2, rl.Maroon)
		}

		rl.EndMode3D()

		rl.DrawCircle(screenWidth/2, screenHeight/2, 3, rl.Black)
		// rl.DrawRectangle(10, 10, 220, 70, rl.Fade(rl.SkyBlue, 0.5))
		// rl.DrawRectangleLines(10, 10, 220, 70, rl.Blue)
		rl.DrawFPS(20, 20)

		// rl.DrawText("First person camera default controls:", 20, 20, 10, rl.Black)
		// rl.DrawText("- Move with keys: W, A, S, D", 40, 40, 10, rl.DarkGray)
		// rl.DrawText("- Mouse move to look around", 40, 60, 10, rl.DarkGray)

		rl.EndDrawing()
	}
}
