package exercises

const (
	PI = 3.14
)

type Circle struct {
	Radius float64
}

func (c Circle) Circumference() float64 {
	return 2 * PI * c.Radius
}

func (c Circle) Area() float64 {
	return PI * c.Radius * c.Radius
}
