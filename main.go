//camera & level design based on: https://github.com/gen2brain/raylib-go/blob/master/examples/core/3d_camera_first_person/main.go

package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
	gravity      = 9.81
	MAX_SPEED    = 100.0
	SPEED_SCALER = 20.0
)

type ProjectileState struct {
	position     rl.Vector3
	isActive     bool
	startPos     rl.Vector3
	direction    rl.Vector3
	speed        float32
	timeOfFlight float32
	trajectory   []rl.Vector3
	hasLanded    bool
}

func calculateTimeOfFlight(v0y float32, y0 float32) float32 {
	a := float64(gravity / 2)
	b := -float64(v0y)
	c := -float64(y0)

	return float32((-float64(b) + math.Sqrt(float64(b*b-4*a*c))) / (2 * float64(a)))
}

func calculatePosition(startPos rl.Vector3, direction rl.Vector3, speed float32, t float32) rl.Vector3 {
	return rl.NewVector3(
		startPos.X+speed*direction.X*t,
		startPos.Y+speed*direction.Y*t-0.5*gravity*t*t,
		startPos.Z+speed*direction.Z*t,
	)
}

func calculateLandingPoint(direction rl.Vector3, speed float32, startPos rl.Vector3) (rl.Vector3, float32) {
	v0y := speed * direction.Y
	timeOfFlight := calculateTimeOfFlight(v0y, startPos.Y)
	position := calculatePosition(startPos, direction, speed, timeOfFlight)
	position.Y = 0
	return position, timeOfFlight
}

func calculateTrajectory(startPos rl.Vector3, landingPoint rl.Vector3, timeOfFlight float32) (rl.Vector3, float32) {
	horizontalDistance := rl.NewVector3(landingPoint.X-startPos.X, 0, landingPoint.Z-startPos.Z)
	horizontalSpeed := rl.Vector3Length(horizontalDistance) / timeOfFlight
	v0y := (landingPoint.Y - startPos.Y + 0.5*gravity*timeOfFlight*timeOfFlight) / timeOfFlight
	totalSpeed := float32(math.Sqrt(float64(horizontalSpeed*horizontalSpeed + v0y*v0y)))

	var direction rl.Vector3
	if rl.Vector3Length(horizontalDistance) > 0 {
		direction = rl.NewVector3(
			horizontalDistance.X/rl.Vector3Length(horizontalDistance)*horizontalSpeed/totalSpeed,
			v0y/totalSpeed,
			horizontalDistance.Z/rl.Vector3Length(horizontalDistance)*horizontalSpeed/totalSpeed,
		)
	} else {
		direction = rl.NewVector3(0, v0y/totalSpeed, 0)
	}

	return direction, totalSpeed
}

func calculateIntersectionY0(start rl.Vector3, direction rl.Vector3) rl.Vector3 {
	direction = rl.Vector3Normalize(direction)

	if direction.Y == 0 {
		return rl.NewVector3(start.X-direction.X*10, start.Y, start.Z-direction.Z*10)
	}

	// rozwiązanie równania liniowego: 0=direction.Y*t+start.Y, t - zmienna parametryczna
	t := -start.Y / direction.Y
	intersection := rl.NewVector3(start.X+t*direction.X, 0, start.Z+t*direction.Z)

	if rl.Vector3DotProduct(direction, rl.Vector3Subtract(intersection, start)) > 0 {
		return rl.NewVector3(start.X-direction.X*10, start.Y-direction.Y*10, start.Z-direction.Z*10)
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
		return rl.NewVector3(start.X-direction.X*10, start.Y-direction.Y*10, start.Z-direction.Z*10)
	}
}

func calculateTrajectoryPoints(startPos rl.Vector3, direction rl.Vector3, speed float32, timeOfFlight float32) []rl.Vector3 {
	points := make([]rl.Vector3, 0)

	landingPoint := calculatePosition(startPos, direction, speed, timeOfFlight)
	totalDistance := rl.Vector3Distance(startPos, landingPoint)

	baseSteps := 30
	additionalSteps := int(totalDistance / 5)
	steps := baseSteps + additionalSteps

	if steps > 100 {
		steps = 100
	}

	deltaTime := timeOfFlight / float32(steps)

	for i := 0; i <= steps; i++ {
		t := deltaTime * float32(i)
		point := calculatePosition(startPos, direction, speed, t)
		if point.Y < 0 {
			point.Y = 0
			points = append(points, point)
			break
		}
		points = append(points, point)
	}
	return points
}

func drawTrajectory(points []rl.Vector3, color rl.Color) {
	if len(points) < 2 {
		return
	}

	for i := 0; i < len(points)-1; i++ {
		rl.DrawLine3D(points[i], points[i+1], color)
		if i%5 == 0 {
			rl.DrawSphere(points[i], 0.05, color)
		}
	}
}

func main() {
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
		positions[i] = rl.NewVector3(float32(rl.GetRandomValue(-128, 128)), heights[i]/2, float32(rl.GetRandomValue(-128, 128)))
		colors[i] = rl.NewColor(uint8(rl.GetRandomValue(20, 255)), uint8(rl.GetRandomValue(10, 55)), 30, 255)
	}

	rl.SetTargetFPS(60)

	var (
		time               float32
		mousePressTime     float32
		isMousePressed     bool
		wasMousePressed    bool
		wasMouseReleased   bool
		landingPoint       rl.Vector3
		numberSpacePressed int
		Pause              bool
		tooLate            bool
		greenTime          float32
		showTrajectory     bool
	)

	RedP := ProjectileState{}
	GreenP := ProjectileState{}

	handleInput := func() {
		if rl.IsKeyDown(rl.KeyLeftShift) {
			camera.Position.Y += 0.01
		} else if rl.IsKeyDown(rl.KeyLeftControl) {
			camera.Position.Y -= 0.01
		}

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			mousePressTime = 0.0
			isMousePressed = true
			wasMousePressed = true
			showTrajectory = false
			RedP.hasLanded = false
			GreenP.hasLanded = false
		}

		if isMousePressed {
			mousePressTime += rl.GetFrameTime()
			currentDirection := rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
			initialSpeed := mousePressTime * SPEED_SCALER
			if initialSpeed > MAX_SPEED {
				initialSpeed = MAX_SPEED
			}
			landingPoint, RedP.timeOfFlight = calculateLandingPoint(currentDirection, initialSpeed, camera.Position)
			RedP.trajectory = calculateTrajectoryPoints(camera.Position, currentDirection, initialSpeed, RedP.timeOfFlight)
		}

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			isMousePressed = false
			wasMouseReleased = true

			RedP.startPos = camera.Position
			RedP.direction = rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
			RedP.speed = mousePressTime * SPEED_SCALER
			if RedP.speed > MAX_SPEED {
				RedP.speed = MAX_SPEED
			}

			RedP.position = RedP.startPos
			GreenP.position = RedP.startPos
			RedP.trajectory = calculateTrajectoryPoints(RedP.startPos, RedP.direction, RedP.speed, RedP.timeOfFlight)
		}

		if rl.IsKeyPressed(rl.KeySpace) && wasMousePressed {
			numberSpacePressed++
			if numberSpacePressed%2 == 1 {
				RedP.isActive = true
				time = 0.0
				tooLate = false
			} else {
				GreenP.timeOfFlight = RedP.timeOfFlight - time
				showTrajectory = true
				if GreenP.timeOfFlight > 0 {
					GreenP.direction, GreenP.speed = calculateTrajectory(RedP.startPos, landingPoint, GreenP.timeOfFlight)
					GreenP.isActive = true
					greenTime = 0
					GreenP.trajectory = calculateTrajectoryPoints(RedP.startPos, GreenP.direction, GreenP.speed, GreenP.timeOfFlight)
				} else {
					tooLate = true
				}
			}
		}

		if rl.IsKeyPressed(rl.KeyP) {
			Pause = !Pause
		}
	}

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)
		if !Pause {
			time += rl.GetFrameTime()
		}
		rl.DisableCursor()

		handleInput()

		if RedP.isActive {
			RedP.position = calculatePosition(RedP.startPos, RedP.direction, RedP.speed, time)

			if RedP.position.Y <= 0 {
				RedP.isActive = false
				v0y := RedP.speed * RedP.direction.Y
				timeOfFlight := calculateTimeOfFlight(v0y, RedP.startPos.Y)
				RedP.position = calculatePosition(RedP.startPos, RedP.direction, RedP.speed, timeOfFlight)
				RedP.position.Y = 0
			}
		}

		if GreenP.isActive {
			greenTime += rl.GetFrameTime()

			if greenTime >= GreenP.timeOfFlight {
				GreenP.isActive = false
				GreenP.position = landingPoint
			} else {
				newPos := calculatePosition(RedP.startPos, GreenP.direction, GreenP.speed, greenTime)
				if newPos.Y <= 0 {
					GreenP.isActive = false
					GreenP.position = landingPoint
				} else {
					GreenP.position = newPos
				}
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.SkyBlue)

		rl.BeginMode3D(camera)

		if wasMousePressed {
			rl.DrawSphere(landingPoint, 0.1, rl.Blue)
			rl.DrawCircle3D(landingPoint, 1.0, rl.NewVector3(1, 0, 0), 90, rl.Fade(rl.Blue, 0.5))

			var cylinderStart, cylinderEnd rl.Vector3
			if numberSpacePressed%2 == 0 && GreenP.isActive && GreenP.timeOfFlight > 0 {
				cylinderStart = rl.NewVector3(
					RedP.startPos.X-GreenP.direction.X*0.5,
					RedP.startPos.Y-GreenP.direction.Y*0.5,
					RedP.startPos.Z-GreenP.direction.Z*0.5,
				)
				cylinderEnd = calculateIntersectionY0(cylinderStart, GreenP.direction)
			} else if RedP.timeOfFlight > 0 {
				cylinderStart = rl.NewVector3(
					RedP.startPos.X-RedP.direction.X*0.5,
					RedP.startPos.Y-RedP.direction.Y*0.5,
					RedP.startPos.Z-RedP.direction.Z*0.5,
				)
				cylinderEnd = calculateIntersectionY0(cylinderStart, RedP.direction)
			}

			if rl.Vector3Distance(cylinderStart, cylinderEnd) > 0 {
				rl.DrawCylinderEx(cylinderStart, cylinderEnd, 0.2, 0.2, 10, rl.DarkGray)
				if rl.Vector3Distance(cylinderStart, cylinderEnd) > 9.9 {
					rl.DrawCylinderEx(cylinderEnd, rl.NewVector3(cylinderEnd.X, -0.1, cylinderEnd.Z), 0.2, 0.2, 10, rl.Black)
				}
			}

			if RedP.timeOfFlight > 0 && len(RedP.trajectory) > 1 {
				drawTrajectory(RedP.trajectory, rl.Fade(rl.Red, 0.5))
			}
			if showTrajectory && GreenP.timeOfFlight > 0 && len(GreenP.trajectory) > 1 {
				drawTrajectory(GreenP.trajectory, rl.Fade(rl.Green, 0.5))
			}
		}

		if wasMouseReleased {
			rl.DrawSphere(RedP.position, 0.2, rl.Red)
			rl.DrawSphere(GreenP.position, 0.2, rl.Green)
		}

		rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(256, 256), rl.LightGray)
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
		rl.DrawText(fmt.Sprintf("Time of flight: %f", RedP.timeOfFlight), 15, 120, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Time to arrive at destination: %f", GreenP.timeOfFlight), 15, 140, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Number of space pressed: %d", numberSpacePressed), 15, 160, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Parameters of RedP - vector Y: %.3f, speed: %.3f", RedP.direction.Y, RedP.speed), 15, 200, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Parameters of GreenP - vector Y: %.3f, speed: %.3f", GreenP.direction.Y, GreenP.speed), 15, 220, 16, rl.White)
		if tooLate {
			rl.DrawText("You tried to shoot to late!", 15, 180, 24, rl.White)
		}

		rl.EndDrawing()
	}
}
