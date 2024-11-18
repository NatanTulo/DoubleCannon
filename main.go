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

	calculateTimeOfFlight := func(v0y float32, y0 float32) float32 {
		a := gravity / 2
		b := -v0y
		c := -y0

		return float32((-float64(b) + math.Sqrt(float64(b*b-4*a*c))) / (2 * float64(a)))
	}

	calculatePosition := func(startPos rl.Vector3, direction rl.Vector3, speed float32, t float32) rl.Vector3 {
		return rl.NewVector3(
			startPos.X+speed*direction.X*t,
			startPos.Y+speed*direction.Y*t-0.5*gravity*t*t,
			startPos.Z+speed*direction.Z*t,
		)
	}

	calculateIntersectionY0 := func(start rl.Vector3, direction rl.Vector3) rl.Vector3 {
		direction = rl.Vector3Normalize(direction)

		if direction.Y == 0 {
			return rl.NewVector3(
				start.X-direction.X*10,
				start.Y,
				start.Z-direction.Z*10,
			)
		}

		t := -start.Y / direction.Y
		intersection := rl.NewVector3(
			start.X+t*direction.X,
			0,
			start.Z+t*direction.Z,
		)

		toIntersection := rl.Vector3Subtract(intersection, start)
		dotProduct := rl.Vector3DotProduct(direction, toIntersection)

		if dotProduct > 0 {
			return rl.NewVector3(
				start.X-direction.X*10,
				start.Y-direction.Y*10,
				start.Z-direction.Z*10,
			)
		}

		distance := rl.Vector3Distance(start, intersection)

		if direction.Y < 0 {
			if distance > 10 {
				scale := 10 / distance
				return rl.NewVector3(
					start.X+(intersection.X-start.X)*scale,
					start.Y+(intersection.Y-start.Y)*scale,
					start.Z+(intersection.Z-start.Z)*scale,
				)
			}
			return intersection
		} else {
			return rl.NewVector3(
				start.X-direction.X*10,
				start.Y-direction.Y*10,
				start.Z-direction.Z*10,
			)
		}
	}

	calculateLandingPoint := func(direction rl.Vector3, speed float32, startPos rl.Vector3) rl.Vector3 {
		v0y := speed * direction.Y
		timeOfFlight := calculateTimeOfFlight(v0y, startPos.Y)
		position := calculatePosition(startPos, direction, speed, timeOfFlight)
		position.Y = 0
		return position
	}

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)
		time += rl.GetFrameTime()

		rl.SetMousePosition(int(screenWidth)/2, int(screenHeight)/2)
		rl.DisableCursor()

		if rl.IsKeyDown(rl.KeyLeftShift) {
			camera.Position.Y += 0.01
		} else if rl.IsKeyDown(rl.KeyLeftControl) {
			camera.Position.Y -= 0.01
		}

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			mousePressTime = 0.0
			isMousePressed = true
		}

		if isMousePressed {
			mousePressTime += rl.GetFrameTime()
			currentDirection := rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
			initialSpeed := mousePressTime * SPEED_SCALER
			if initialSpeed > MAX_SPEED {
				initialSpeed = MAX_SPEED
			}
			landingPoint = calculateLandingPoint(currentDirection, initialSpeed, camera.Position)
		}

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			isMousePressed = false
			wasMousePressed = true

			projectileStartPos = camera.Position
			projectileDirection = rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
			projectileSpeed = mousePressTime * SPEED_SCALER
			if projectileSpeed > MAX_SPEED {
				projectileSpeed = MAX_SPEED
			}

			spherePosition = projectileStartPos
		}

		if rl.IsKeyPressed(rl.KeySpace) && wasMousePressed {
			isProjectileActive = true
			time = 0.0
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
			cylinderStart := rl.NewVector3(
				projectileStartPos.X-projectileDirection.X*0.5,
				projectileStartPos.Y-projectileDirection.Y*0.5,
				projectileStartPos.Z-projectileDirection.Z*0.5,
			)
			cylinderEnd := calculateIntersectionY0(cylinderStart, projectileDirection)

			rl.DrawCylinderEx(cylinderStart, cylinderEnd, 0.2, 0.2, 10, rl.DarkGray)
			if rl.Vector3Distance(cylinderStart, cylinderEnd) > 9.9 {
				rl.DrawCylinderEx(cylinderEnd, rl.NewVector3(cylinderEnd.X, -0.1, cylinderEnd.Z), 0.2, 0.2, 10, rl.Black)
			}
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
