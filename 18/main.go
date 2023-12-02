package main

import (
	"fmt"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	var cubes []Cube1x1

	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)
	for _, line := range lines {
		parts := aocutil.SplitN(line, ",", 3)
		cubes = append(cubes, Cube1x1{
			X: aocutil.Atoi[int](parts[0]),
			Y: aocutil.Atoi[int](parts[1]),
			Z: aocutil.Atoi[int](parts[2]),
		})
	}

	world := aocutil.NewSet[Cube1x1](len(cubes))
	for _, cube := range cubes {
		world.Add(cube)
	}

	part1(world)
	part2(world)
}

func part1(world aocutil.Set[Cube1x1]) {
	var sides int

	for cube := range world {
		exposed := 6
		// Check all 6 sides.
		for _, delta := range deltas {
			other := cube.Position().Add(delta)
			if world.Has(Cube1x1(other)) {
				exposed--
			}
		}
		sides += exposed
	}

	fmt.Println("part 1:", sides)
}

func part2(world aocutil.Set[Cube1x1]) {
	// Find the minimum and maximum bounds of this world.
	var min, max Pt
	for cube := range world {
		pos := cube.Position()
		min.X = aocutil.Min2(min.X, pos.X)
		min.Y = aocutil.Min2(min.Y, pos.Y)
		min.Z = aocutil.Min2(min.Z, pos.Z)
		max.X = aocutil.Max2(max.X, pos.X)
		max.Y = aocutil.Max2(max.Y, pos.Y)
		max.Z = aocutil.Max2(max.Z, pos.Z)
	}

	var sides int

	for cube := range world {
		exposed := 6
		// Check all 6 directions for each cube. We'll basically shoot a ray
		// from the cube to the edge of the world, and see if we hit a cube.
	raycast:
		for _, delta := range deltas {
			// Shoot a ray from the cube to the edge of the world.
			for ray := cube.Position().Add(delta); !ray.Lt(min) && !ray.Gt(max); ray = ray.Add(delta) {
				// If we hit a cube, then we know that this side is not exposed.
				// so we don't count it.
				if world.Has(Cube1x1(ray)) {
					exposed--
					continue raycast
				}
			}
		}
		sides += exposed
	}

	fmt.Println("part 2:", sides)
}

// Cube1x1 is a 1x1 cube.
type Cube1x1 Pt

func (c Cube1x1) Position() Pt { return Pt(c) }

type Pt struct{ X, Y, Z int }

func (pt Pt) Add(other Pt) Pt {
	return Pt{
		X: pt.X + other.X,
		Y: pt.Y + other.Y,
		Z: pt.Z + other.Z,
	}
}

func (pt Pt) Lt(other Pt) bool {
	return pt.X < other.X || pt.Y < other.Y || pt.Z < other.Z
}

func (pt Pt) Gt(other Pt) bool {
	return pt.X > other.X || pt.Y > other.Y || pt.Z > other.Z
}

var deltas = []Pt{
	{1, 0, 0},
	{-1, 0, 0},
	{0, 1, 0},
	{0, -1, 0},
	{0, 0, 1},
	{0, 0, -1},
}

func diff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
