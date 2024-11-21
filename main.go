//camera & level design based on: https://github.com/gen2brain/raylib-go/blob/master/examples/core/3d_camera_first_person/main.go

package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1920)
	screenHeight := int32(1080)

	rl.InitWindow(screenWidth, screenHeight, "WNO projekt - 193527")
	defer rl.CloseWindow()

	rl.SetConfigFlags(rl.FlagVsyncHint)
	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(4.0, 2.0, 4.0)
	camera.Target = rl.NewVector3(0.0, 1.8, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 60.0
	camera.Projection = rl.CameraPerspective

	heights := make([]float32, 20)
	positions := make([]rl.Vector3, 20)
	colors := make([]rl.Color, 20)

	for i := 0; i < 20; i++ {
		heights[i] = float32(rl.GetRandomValue(1, 12))
		positions[i] = rl.NewVector3(float32(rl.GetRandomValue(-15, 15)), heights[i]/2, float32(rl.GetRandomValue(-15, 15)))
		colors[i] = rl.NewColor(uint8(rl.GetRandomValue(20, 255)), uint8(rl.GetRandomValue(10, 55)), 30, 255)
	}

	rl.SetTargetFPS(60)
	var time float32
	var wasMousePressed bool = false
	var gravity float32 = 9.81
	var mousePressTime float32
	var isMousePressed bool
	const MAX_SPEED float32 = 50.0
	const SPEED_SCALER float32 = 20.0
	var landingPoint rl.Vector3
	var numberSpacePressed int
	var Pause bool
	var tooLate bool

	var spherePositionRed rl.Vector3
	var isProjectileActiveRed bool
	var projectileStartPosRed rl.Vector3
	var projectileDirectionRed rl.Vector3
	var projectileSpeedRed float32
	var timeOfFlightRed float32

	var spherePositionGreen rl.Vector3
	var isProjectileActiveGreen bool
	var projectileDirectionGreen rl.Vector3
	var projectileSpeedGreen float32
	var timeOfFlightGreen float32
	var greenTime float32

	calculateTrajectory := func(startPos rl.Vector3, landingPoint rl.Vector3, timeOfFlight float32) (rl.Vector3, float32) {
		horizontalDistance := rl.NewVector3(landingPoint.X-startPos.X, 0, landingPoint.Z-startPos.Z)
		horizontalSpeed := rl.Vector3Length(horizontalDistance) / timeOfFlight
		v0y := (landingPoint.Y - startPos.Y + 0.5*gravity*timeOfFlight*timeOfFlight) / timeOfFlight
		totalSpeed := float32(math.Sqrt(
			math.Pow(float64(horizontalSpeed), 2) +
				math.Pow(float64(v0y), 2),
		))

		direction := rl.NewVector3(
			horizontalDistance.X/rl.Vector3Length(horizontalDistance),
			v0y/totalSpeed,
			horizontalDistance.Z/rl.Vector3Length(horizontalDistance),
		)

		return direction, totalSpeed
	}

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
		timeOfFlightRed = calculateTimeOfFlight(v0y, startPos.Y)
		position := calculatePosition(startPos, direction, speed, timeOfFlightRed)
		position.Y = 0
		return position
	}

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)
		if !Pause {
			time += rl.GetFrameTime()
		}
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

			projectileStartPosRed = camera.Position
			projectileDirectionRed = rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
			projectileSpeedRed = mousePressTime * SPEED_SCALER
			if projectileSpeedRed > MAX_SPEED {
				projectileSpeedRed = MAX_SPEED
			}

			spherePositionRed = projectileStartPosRed
			spherePositionGreen = projectileStartPosRed
		}

		if rl.IsKeyPressed(rl.KeySpace) && wasMousePressed {
			numberSpacePressed++
			if numberSpacePressed%2 == 1 {
				isProjectileActiveRed = true
				time = 0.0
				tooLate = false
			} else {
				timeOfFlightGreen = timeOfFlightRed - time
				if timeOfFlightGreen > 0 {
					projectileDirectionGreen, projectileSpeedGreen = calculateTrajectory(projectileStartPosRed, landingPoint, timeOfFlightGreen)

					isProjectileActiveGreen = true
					greenTime = 0
				} else {
					tooLate = true
				}
			}
		}

		if rl.IsKeyPressed(rl.KeyP) {
			Pause = !Pause
		}

		if isProjectileActiveRed {
			spherePositionRed = calculatePosition(projectileStartPosRed, projectileDirectionRed, projectileSpeedRed, time)

			if spherePositionRed.Y <= 0 {
				isProjectileActiveRed = false
				v0y := projectileSpeedRed * projectileDirectionRed.Y
				timeOfFlight := calculateTimeOfFlight(v0y, projectileStartPosRed.Y)
				spherePositionRed = calculatePosition(projectileStartPosRed, projectileDirectionRed, projectileSpeedRed, timeOfFlight)
				spherePositionRed.Y = 0
			}
		}

		if isProjectileActiveGreen {
			greenTime += rl.GetFrameTime()

			if greenTime >= timeOfFlightGreen {
				isProjectileActiveGreen = false
			} else {
				spherePositionGreen = calculatePosition(projectileStartPosRed, projectileDirectionGreen, projectileSpeedGreen, greenTime)
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.SkyBlue)

		rl.BeginMode3D(camera)

		if wasMousePressed {
			rl.DrawSphere(spherePositionRed, 0.2, rl.Red)
			rl.DrawSphere(landingPoint, 0.1, rl.Blue)
			rl.DrawSphere(spherePositionGreen, 0.2, rl.Green)
			rl.DrawCircle3D(landingPoint, 1.0, rl.NewVector3(1, 0, 0), 90, rl.Fade(rl.Blue, 0.5))
			cylinderStart := rl.NewVector3(
				projectileStartPosRed.X-projectileDirectionRed.X*0.5,
				projectileStartPosRed.Y-projectileDirectionRed.Y*0.5,
				projectileStartPosRed.Z-projectileDirectionRed.Z*0.5,
			)
			cylinderEnd := calculateIntersectionY0(cylinderStart, projectileDirectionRed)

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

		rl.DrawText(fmt.Sprintf("Time of flight: %f", timeOfFlightRed), 20, 120, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Time to arrive at destination: %f", timeOfFlightGreen), 20, 140, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Number of space pressed: %d", numberSpacePressed), 20, 160, 16, rl.White)
		if tooLate {
			rl.DrawText("You tried to shoot to late!", 20, 180, 24, rl.White)
		}

		rl.EndDrawing()
	}
}
