package algo

import "math"

type Equation struct {
	ID        int                                  `json:"id"`
	Expr      string                               `json:"expr"`
	ExactExpr string                               `json:"exactExpr"`
	F         func(x, y float64) float64           `json:"-"`
	Exact     func(x0, y0, x float64) float64      `json:"-"`
}

var Equations = []Equation{
	{
		ID:        1,
		Expr:      "y' = y",
		ExactExpr: "y = y0 * exp(x - x0)",
		F:         func(x, y float64) float64 { return y },
		Exact:     func(x0, y0, x float64) float64 { return y0 * math.Exp(x-x0) },
	},
	{
		ID:        2,
		Expr:      "y' = x + y",
		ExactExpr: "y = (y0 + x0 + 1) * exp(x - x0) - x - 1",
		F:         func(x, y float64) float64 { return x + y },
		Exact: func(x0, y0, x float64) float64 {
			return (y0+x0+1)*math.Exp(x-x0) - x - 1
		},
	},
	{
		ID:        3,
		Expr:      "y' = -2*x*y",
		ExactExpr: "y = y0 * exp(x0^2 - x^2)",
		F:         func(x, y float64) float64 { return -2 * x * y },
		Exact:     func(x0, y0, x float64) float64 { return y0 * math.Exp(x0*x0-x*x) },
	},
	{
		ID:        4,
		Expr:      "y' = y * cos(x)",
		ExactExpr: "y = y0 * exp(sin(x) - sin(x0))",
		F:         func(x, y float64) float64 { return y * math.Cos(x) },
		Exact: func(x0, y0, x float64) float64 {
			return y0 * math.Exp(math.Sin(x)-math.Sin(x0))
		},
	},
}

func FindEquation(id int) (Equation, bool) {
	for _, eq := range Equations {
		if eq.ID == id {
			return eq, true
		}
	}
	return Equation{}, false
}
