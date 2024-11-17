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
	var time float32
	var isProjectileActive bool
	var wasMousePressed bool = false
	var initialSpeed float32
	var gravity float32 = 9.81
	var mousePressTime float32
	var isMousePressed bool
	const MAX_SPEED float32 = 50.0
	const SPEED_SCALER float32 = 20.0

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)
		time += rl.GetFrameTime()

		rl.SetMousePosition(int(screenWidth)/2, int(screenHeight)/2)
		rl.DisableCursor()

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			mousePressTime = 0.0
			isMousePressed = true
		}

		if isMousePressed {
			mousePressTime += rl.GetFrameTime()
		}

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			isMousePressed = false
			wasMousePressed = true
			initialSpeed = mousePressTime * SPEED_SCALER
			if initialSpeed > MAX_SPEED {
				initialSpeed = MAX_SPEED
			}
			time = 0.0
			initialTarget = camera.Target
			initialDirection = rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
			spherePosition = initialTarget
			isProjectileActive = true
		}
		// Obliczanie pozycji kuli zgodnie z trajektorią
		if isProjectileActive {
			spherePosition = rl.NewVector3(
				initialTarget.X+initialSpeed*initialDirection.X*time,
				initialTarget.Y+initialSpeed*initialDirection.Y*time-0.5*gravity*time*time,
				initialTarget.Z+initialSpeed*initialDirection.Z*time,
			)

			// Zatrzymanie ruchu, jeśli kula spadnie poniżej płaszczyzny
			if spherePosition.Y <= 0 {
				isProjectileActive = false
				spherePosition.Y = 0
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.SkyBlue)

		rl.BeginMode3D(camera)

		if wasMousePressed {
			rl.DrawSphere(spherePosition, 0.1, rl.Red)
		}
		rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(128, 128), rl.LightGray)

		for i := 0; i < 20; i++ {
			rl.DrawCube(positions[i], 2, heights[i], 2, colors[i])
			rl.DrawCubeWires(positions[i], 2, heights[i], 2, rl.Maroon)
		}

		rl.EndMode3D()

		rl.DrawCircle(screenWidth/2, screenHeight/2, 3, rl.Black)

		if mousePressTime*SPEED_SCALER < MAX_SPEED {
			rl.DrawRectangle(10, 50, int32(220.0*mousePressTime*SPEED_SCALER/MAX_SPEED), 50, rl.Fade(rl.Black, 0.5))
		} else {
			rl.DrawRectangle(10, 50, 220, 50, rl.Fade(rl.Red, 0.5))
		}
		rl.DrawRectangleLines(10, 50, 220, 50, rl.Blue)

		rl.DrawFPS(20, 20)

		rl.EndDrawing()
	}
}
