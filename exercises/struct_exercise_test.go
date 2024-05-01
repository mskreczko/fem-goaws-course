package exercises

import "testing"

func testCircumference(t *testing.T) {
	circle := Circle{Radius: 3}

	circumference := circle.Circumference()

	if circumference != 2*PI*3 {
		t.Fatal("Invalid result")
	}
}

func testArea(t *testing.T) {
	circle := Circle{Radius: 3}

	area := circle.Area()

	if area != PI*3*3 {
		t.Fatal("Invalid result")
	}
}
