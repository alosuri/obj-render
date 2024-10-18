package main

import (
	"math"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TriangleWithDp struct {
	triangle [][]float64
	dp       float64
}

func main() {
	// Variables
	var screen_h = 900
	var screen_w = 800

	var near float64 = 0.1
	var far float64 = 1000.0
	var fov float64 = 90
	var aspect_ratio float64 = float64(screen_h) / float64(screen_w)
	var fov_rad = 1 / math.Tan(fov*0.5/180*3.14159)

	rl.InitWindow(int32(screen_w), int32(screen_h), "Window")
	rl.SetTargetFPS(60)

	load_file()

	mesh := load_file()

	matrix4 := [][]float64{
		{aspect_ratio * fov_rad, 0, 0, 0},
		{0, fov_rad, 0, 0},
		{0, 0, far / (far - near), 1.0},
		{0, 0, (-far * near) / (far - near), 0.0},
	}

	camera := []float64{0.0, 0.0, 0.0}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		var thetha float64 = 1 * rl.GetTime()

		matrixRotZ := [][]float64{
			{math.Cos(thetha), math.Sin(thetha), 0, 0},
			{-math.Sin(thetha), math.Cos(thetha), 0, 0},
			{0, 0, 1.0, 0},
			{0, 0, 0, 1.0},
		}

		matrixRotX := [][]float64{
			{1.0, 0, 0, 0},
			{0, math.Cos(thetha * 0.5), math.Sin(thetha * 0.5), 0},
			{0, -math.Sin(thetha * 0.5), math.Cos(thetha * 0.5), 0},
			{0, 0, 0, 1},
		}

		trianglesToRaster := []TriangleWithDp{}

		for _, mesh := range mesh {

			meshRotationZ := [][]float64{MultiplyVector(mesh[0], matrixRotZ), MultiplyVector(mesh[1], matrixRotZ), MultiplyVector(mesh[2], matrixRotZ)}
			meshRotationX := [][]float64{MultiplyVector(meshRotationZ[0], matrixRotX), MultiplyVector(meshRotationZ[1], matrixRotX), MultiplyVector(meshRotationZ[2], matrixRotX)}

			meshRotationX[0][2] += 3
			meshRotationX[1][2] += 3
			meshRotationX[2][2] += 3

			normal := [3]float64{}
			line1 := [3]float64{}
			line2 := [3]float64{}

			line1[0] = meshRotationX[1][0] - meshRotationX[0][0]
			line1[1] = meshRotationX[1][1] - meshRotationX[0][1]
			line1[2] = meshRotationX[1][2] - meshRotationX[0][2]

			line2[0] = meshRotationX[2][0] - meshRotationX[0][0]
			line2[1] = meshRotationX[2][1] - meshRotationX[0][1]
			line2[2] = meshRotationX[2][2] - meshRotationX[0][2]

			normal[0] = line1[1]*line2[2] - line1[2]*line2[1]
			normal[1] = line1[2]*line2[0] - line1[0]*line2[2]
			normal[2] = line1[0]*line2[1] - line1[1]*line2[0]

			// var l float64 = math.Sqrt(normal[0]*normal[0] + normal[1]*normal[1] + normal[2]*normal[2])

			// if normal[2] < 0 {
			if (normal[0]*meshRotationX[0][0]-camera[0])+(normal[1]*meshRotationX[0][1]-camera[1])+(normal[2]*meshRotationX[0][2]-camera[2]) < 0.0 {
				light_direction := []float64{0.0, 0.0, -1.0}
				// var l_light float64 = math.Sqrt(light_direction[0]*light_direction[0] + light_direction[1]*light_direction[1] + light_direction[2]*light_direction[2])

				var dp float64 = normal[0]*light_direction[0] + normal[1]*light_direction[1] + normal[2]*light_direction[2]

				meshFinal := [][]float64{MultiplyVector(meshRotationX[0], matrix4), MultiplyVector(meshRotationX[1], matrix4), MultiplyVector(meshRotationX[2], matrix4)}

				meshFinal[0][0] += 1
				meshFinal[1][0] += 1
				meshFinal[2][0] += 1
				meshFinal[0][1] += 1
				meshFinal[1][1] += 1
				meshFinal[2][1] += 1

				trianglesToRaster = append(trianglesToRaster, TriangleWithDp{triangle: meshFinal, dp: dp})

			}
		}

		sort.Slice(trianglesToRaster, func(i, j int) bool {
			z1 := (trianglesToRaster[i].triangle[0][2] + trianglesToRaster[i].triangle[1][2] + trianglesToRaster[i].triangle[2][2]) / 3.0
			z2 := (trianglesToRaster[j].triangle[0][2] + trianglesToRaster[j].triangle[1][2] + trianglesToRaster[j].triangle[2][2]) / 3.0
			return z1 > z2
		})

		for _, i := range trianglesToRaster {
			rl.DrawTriangle(
				rl.NewVector2(float32(i.triangle[0][0]*0.5*float64(screen_w)), float32(i.triangle[0][1]*0.5*float64(screen_h))),
				rl.NewVector2(float32(i.triangle[1][0]*0.5*float64(screen_w)), float32(i.triangle[1][1]*0.5*float64(screen_h))),
				rl.NewVector2(float32(i.triangle[2][0]*0.5*float64(screen_w)), float32(i.triangle[2][1]*0.5*float64(screen_h))),
				rl.ColorContrast(rl.Blue, float32(i.dp)))
			rl.DrawTriangleLines(
				rl.NewVector2(float32(i.triangle[0][0]*0.5*float64(screen_w)), float32(i.triangle[0][1]*0.5*float64(screen_h))),
				rl.NewVector2(float32(i.triangle[1][0]*0.5*float64(screen_w)), float32(i.triangle[1][1]*0.5*float64(screen_h))),
				rl.NewVector2(float32(i.triangle[2][0]*0.5*float64(screen_w)), float32(i.triangle[2][1]*0.5*float64(screen_h))),
				rl.Black)
		}

		rl.EndDrawing()
	}
}

func MultiplyVector(m []float64, v [][]float64) []float64 {
	x := m[0]*v[0][0] + m[1]*v[1][0] + m[2]*v[2][0] + v[3][0]
	y := m[0]*v[0][1] + m[1]*v[1][1] + m[2]*v[2][1] + v[3][1]
	z := m[0]*v[0][2] + m[1]*v[1][2] + m[2]*v[2][2] + v[3][2]
	w := m[0]*v[0][3] + m[1]*v[1][3] + m[2]*v[2][3] + v[3][3]

	if w != 0 {
		x /= w
		y /= w
		z /= w
	}

	return []float64{x, y, z}
}
