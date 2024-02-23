// Project Name: Go Interface
// Project Description: The program draws
// the geometry (either Rectangle, Circle, or Triangle) on the screen (implemented as a logical
// device of type Display)
// Name: Abdullah Irfan
// Netid: airfan6
// Date: 12/01/2023
package main

import (
	"errors"
	"fmt"
	"os"
)

// Color represents a color in RGB format.
type Color struct {
	R, G, B int
}

// Colors represents a set of predefined colors.
var Colors = map[int]Color{
	1: {255, 0, 0},     // red
	2: {0, 255, 0},     // green
	3: {0, 0, 255},     // blue
	4: {255, 255, 0},   // yellow
	5: {255, 164, 0},   // orange
	6: {128, 0, 128},   // purple
	7: {165, 42, 42},   // brown
	8: {0, 0, 0},       // black
	9: {255, 255, 255}, // white
}

// ColorMap maps color IDs to RGB values.
var ColorMap = make(map[int]Color)

func init() {
	for id, c := range Colors {
		ColorMap[id] = c
	}
}

// Point represents a point in 2D space.
type Point struct {
	X, Y int
}

// Rectangle represents a rectangle in 2D space.
type Rectangle struct {
	LL, UR Point
	C      int
}

// draw function for Rectangle, draws the rectangle on the given screen
func (rect Rectangle) draw(scn Screen) (err error) {
	if outOfBounds(rect.LL, scn) || outOfBounds(rect.UR, scn) {
		return errors.New("geometry  out of bounds")
	}
	if colorUnknown(rect.C) {
		return errors.New("unknown color")
	}

	for x := rect.LL.X; x <= rect.UR.X; x++ {
		for y := rect.LL.Y; y <= rect.UR.Y; y++ {
			scn.drawPixel(x, y, rect.C)
		}
	}
	return nil
}

// shape function for Rectangle, returns the string representation of the shape

func (rect Rectangle) shape() string {
	return "Rectangle"
}

// Circle represents a circle in 2D space.
type Circle struct {
	CP Point
	R  int
	C  int
}

// draw function for Circle, draws the circle on the given screen
func (circ Circle) draw(scn Screen) (err error) {
	if outOfBounds(circ.CP, scn) {
		return errors.New("circle center out of bounds")
	}
	if colorUnknown(circ.C) {
		return errors.New("unknown color")
	}

	radius := circ.R
	x0, y0 := circ.CP.X, circ.CP.Y

	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			if x*x+y*y <= radius*radius {
				scn.drawPixel(x0+x, y0+y, circ.C)
			}
		}
	}
	return nil
}

// shape function for Circle, returns the string representation of the shape
func (circ Circle) shape() string {
	return "Circle"
}

// Triangle represents a triangle in 2D space.
type Triangle struct {
	Pt0, Pt1, Pt2 Point
	C             int
}

// draw function for Triangle, draws the triangle on the given screen
func (tri Triangle) draw(scn Screen) (err error) {
	if outOfBounds(tri.Pt0, scn) || outOfBounds(tri.Pt1, scn) || outOfBounds(tri.Pt2, scn) {
		return errors.New("triangle out of bounds")
	}
	if colorUnknown(tri.C) {
		return errors.New("unknown color")
	}

	y0 := tri.Pt0.Y
	y1 := tri.Pt1.Y
	y2 := tri.Pt2.Y

	// Sort the points so that y0 <= y1 <= y2
	if y1 < y0 {
		tri.Pt1, tri.Pt0 = tri.Pt0, tri.Pt1
	}
	if y2 < y0 {
		tri.Pt2, tri.Pt0 = tri.Pt0, tri.Pt2
	}
	if y2 < y1 {
		tri.Pt2, tri.Pt1 = tri.Pt1, tri.Pt2
	}

	x0, y0 := tri.Pt0.X, tri.Pt0.Y
	x1, y1 := tri.Pt1.X, tri.Pt1.Y
	x2, y2 := tri.Pt2.X, tri.Pt2.Y

	x01 := interpolate(y0, x0, y1, x1)
	x12 := interpolate(y1, x1, y2, x2)
	x02 := interpolate(y0, x0, y2, x2)

	// Concatenate the short sides
	x012 := append(x01[:len(x01)-1], x12...)

	// Determine which is left and which is right
	var xLeft, xRight []int
	m := len(x012) / 2
	if x02[m] < x012[m] {
		xLeft = x02
		xRight = x012
	} else {
		xLeft = x012
		xRight = x02
	}

	// Draw the horizontal segments
	for y := y0; y <= y2; y++ {
		for x := xLeft[y-y0]; x <= xRight[y-y0]; x++ {
			scn.drawPixel(x, y, tri.C)
		}
	}
	return nil
}

// shape function for Triangle, returns the string representation of the shape
func (tri Triangle) shape() string {
	return "Triangle"
}

// Screen is an interface representing a 2D screen.
type Screen interface {
	initialize(maxX, maxY int)
	getMaxXY() (int, int)
	drawPixel(x, y, c int) error
	getPixel(x, y int) (int, error)
	clearScreen()
	screenShot(f string) error
}

// Display represents a logical device with a matrix of colors.
type Display struct {
	maxX, maxY int
	matrix     [][]int
}

// initialize function for Display, initializes the display with the given dimensions
func (d *Display) initialize(maxX, maxY int) {
	d.maxX = maxX
	d.maxY = maxY
	d.matrix = make([][]int, maxY)
	for i := range d.matrix {
		d.matrix[i] = make([]int, maxX)
	}
	d.clearScreen()
}

// getMaxXY function for Display, returns the maximum X and Y dimensions of the display
func (d *Display) getMaxXY() (int, int) {
	return d.maxX, d.maxY
}

// drawPixel function for Display, draws a pixel with the specified color at the given coordinates
func (d *Display) drawPixel(x, y, c int) error {
	if x < 0 || x >= d.maxX || y < 0 || y >= d.maxY {
		return errors.New("pixel out of bounds")
	}
	d.matrix[x][y] = c
	return nil

}

// getPixel function for Display, gets the color of the pixel at the specified coordinates
func (d *Display) getPixel(x, y int) (int, error) {
	if x < 0 || x >= d.maxX || y < 0 || y >= d.maxY {
		return 0, errors.New("pixel out of bounds")
	}
	return d.matrix[y][x], nil
}

// clearScreen function for Display, clears the screen by setting all pixels to white color
func (d *Display) clearScreen() {
	for i := range d.matrix {
		for j := range d.matrix[i] {
			d.matrix[i][j] = 9 // Set to white color
		}
	}
}

// screenShot function for Display, takes a screenshot and saves it to a PPM file
func (d *Display) screenShot(f string) error {
	file, err := os.Create(f + ".ppm")
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "P3\n%d %d\n255\n", d.maxX, d.maxY)
	for _, row := range d.matrix {
		for _, pixel := range row {
			color := ColorMap[pixel]
			fmt.Fprintf(file, "%d %d %d ", color.R, color.G, color.B)
		}
		fmt.Fprintln(file)
	}

	return nil
}

// outOfBounds function checks if a given point is outside the bounds of the screen
func outOfBounds(p Point, s Screen) bool {
	x, y := s.getMaxXY()
	return p.X < 0 || p.X >= x || p.Y < 0 || p.Y >= y
}

// colorUnknown function checks if a given color value is unknown or not in the Colors map
func colorUnknown(c int) bool {
	_, exists := Colors[c]
	return !exists
}

// interpolate function performs linear interpolation between two points in 1D space
func interpolate(l0, d0, l1, d1 int) (values []int) {
	a := float64(d1-d0) / float64(l1-l0)
	d := float64(d0)

	count := l1 - l0 + 1
	for ; count > 0; count-- {
		values = append(values, int(d))
		d = d + a
	}
	return
}
