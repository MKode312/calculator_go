package main

import (
	"math"
)

type Geometry interface {
	Area() float64
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

type Rectangle struct {
	width  float64
	height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

type Rhombus struct {
	d1 float64
	d2 float64
}

func (r Rhombus) Area() float64 {
	return (r.d1 * r.d2) / 2
}

type Parallelogram struct {
	a float64
	h float64
}

func (p Parallelogram) Area() float64 {
	return p.a * p.h
}
