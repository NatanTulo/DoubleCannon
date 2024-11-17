package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1920)
	screenHeight := int32(1080)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - 3d camera first person")
	defer rl.CloseWindow()

	rl.SetConfigFlags(rl.FlagVsyncHint)
	camera := rl.Camera3D{
		Position:   rl.NewVector3(4, 2, 4),
		Target:     rl.NewVector3(0, 1.8, 0),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       60,
		Projection: rl.CameraProjection(rl.CameraPerspective),
	}

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
	var gravity float32 = 9.81
	var mousePressTime float32
	var isMousePressed bool
	const MAX_SPEED float32 = 50.0
	const SPEED_SCALER float32 = 20.0
	var landingPoint rl.Vector3
	var spherePosition rl.Vector3

	var projectileStartPos rl.Vector3
	var projectileDirection rl.Vector3
	var projectileSpeed float32

	// Funkcja obliczająca czas lotu
	calculateTimeOfFlight := func(v0y float32, y0 float32) float32 {
		a := gravity / 2
		b := -v0y
		c := -y0

		return float32((-float64(b) + math.Sqrt(float64(b*b-4*a*c))) / (2 * float64(a)))
	}

	// Funkcja obliczająca pozycję w danym czasie
	calculatePosition := func(startPos rl.Vector3, direction rl.Vector3, speed float32, t float32) rl.Vector3 {
		return rl.NewVector3(
			startPos.X+speed*direction.X*t,
			startPos.Y+speed*direction.Y*t-0.5*gravity*t*t,
			startPos.Z+speed*direction.Z*t,
		)
	}

	calculateLandingPoint := func(direction rl.Vector3, speed float32, startPos rl.Vector3) rl.Vector3 {
		v0y := speed * direction.Y
		timeOfFlight := calculateTimeOfFlight(v0y, startPos.Y)
		position := calculatePosition(startPos, direction, speed, timeOfFlight)
		position.Y = 0 // Upewniamy się, że punkt lądowania jest dokładnie na poziomie ziemi
		return position
	}

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
			time = 0.0

			projectileStartPos = camera.Position
			projectileDirection = rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
			projectileSpeed = mousePressTime * SPEED_SCALER
			if projectileSpeed > MAX_SPEED {
				projectileSpeed = MAX_SPEED
			}

			spherePosition = projectileStartPos
			isProjectileActive = true
		}

		if isProjectileActive {
			spherePosition = calculatePosition(projectileStartPos, projectileDirection, projectileSpeed, time)

			if spherePosition.Y <= 0 {
				isProjectileActive = false
				// Obliczamy dokładną pozycję lądowania
				v0y := projectileSpeed * projectileDirection.Y
				timeOfFlight := calculateTimeOfFlight(v0y, projectileStartPos.Y)
				spherePosition = calculatePosition(projectileStartPos, projectileDirection, projectileSpeed, timeOfFlight)
				spherePosition.Y = 0
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.SkyBlue)

		rl.BeginMode3D(camera)

		if wasMousePressed {
			rl.DrawSphere(spherePosition, 0.2, rl.Red)
			rl.DrawSphere(landingPoint, 0.1, rl.Blue)
			rl.DrawCircle3D(landingPoint, 1.0, rl.NewVector3(1, 0, 0), 90, rl.Fade(rl.Blue, 0.5))
		}

		if isMousePressed {
			currentDirection := rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
			initialSpeed := mousePressTime * SPEED_SCALER
			if initialSpeed > MAX_SPEED {
				initialSpeed = MAX_SPEED
			}
			landingPoint = calculateLandingPoint(currentDirection, initialSpeed, camera.Position)
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
