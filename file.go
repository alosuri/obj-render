package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func load_file() [][][]float64 {
	file, err := os.Open("./monkey.obj")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	verts := [][]float64{}
	faces := [][][]float64{}

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "v") {
			parts := strings.Fields(scanner.Text())

			if len(parts) >= 4 {
				x, err1 := strconv.ParseFloat(parts[1], 64)
				y, err2 := strconv.ParseFloat(parts[2], 64)
				z, err3 := strconv.ParseFloat(parts[3], 64)

				if err1 != nil || err2 != nil || err3 != nil {
					log.Println("Error parsing vertex coordinates:", err1, err2, err3)
					continue
				}

				verts = append(verts, []float64{x, y, z})
			}
		}

		if strings.HasPrefix(scanner.Text(), "f") {
			parts := strings.Fields(scanner.Text())

			x, err1 := strconv.ParseInt(parts[1], 10, 16)
			y, err2 := strconv.ParseInt(parts[2], 10, 16)
			z, err3 := strconv.ParseInt(parts[3], 10, 16)

			if err1 != nil || err2 != nil || err3 != nil {
				log.Println("Error parsing vertex coordinates:", err1, err2, err3)
				continue
			}

			faces = append(faces, [][]float64{verts[x-1], verts[y-1], verts[z-1]})
		}
	}
	return faces
}
