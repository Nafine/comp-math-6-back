package algo

import (
	"comp-math-6/internal/numeric"
	"fmt"
)

type SolveInput struct {
	EquationID int
	X0         float64
	Xn         float64
	Y0         float64
	H          float64
	Epsilon    float64
	Methods    []string
}

type TableRow struct {
	I     int     `json:"i"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Exact float64 `json:"exact"`
	Error float64 `json:"error"`
}

type ODEResult struct {
	Method      string          `json:"method"`
	Table       []TableRow      `json:"table"`
	Curve       []numeric.Point `json:"curve"`
	StepUsed    float64         `json:"stepUsed"`
	MaxError    float64         `json:"maxError"`
	Refinements int             `json:"refinements,omitempty"`
	Note        string          `json:"note,omitempty"`
}

type SolveResponse struct {
	Equation   Equation        `json:"equation"`
	ExactCurve []numeric.Point `json:"exactCurve"`
	Results    []ODEResult     `json:"results"`
}

const exactCurveSamples = 200

var supportedMethods = map[string]bool{
	"euler": true,
	"rk4":   true,
	"milne": true,
}

func Solve(in SolveInput) (SolveResponse, error) {
	eq, ok := FindEquation(in.EquationID)
	if !ok {
		return SolveResponse{}, fmt.Errorf("unknown equationId=%d", in.EquationID)
	}
	if in.Xn <= in.X0 {
		return SolveResponse{}, fmt.Errorf("xn must be greater than x0")
	}
	if in.H <= 0 {
		return SolveResponse{}, fmt.Errorf("step h must be positive")
	}
	if in.Epsilon <= 0 {
		return SolveResponse{}, fmt.Errorf("epsilon must be positive")
	}
	if (in.Xn-in.X0)/in.H > 100000 {
		return SolveResponse{}, fmt.Errorf("too many grid nodes; increase h or shrink interval")
	}

	methods := in.Methods
	if len(methods) == 0 {
		methods = []string{"euler", "rk4", "milne"}
	}
	for _, m := range methods {
		if !supportedMethods[m] {
			return SolveResponse{}, fmt.Errorf("unsupported method %q", m)
		}
	}

	resp := SolveResponse{
		Equation:   eq,
		ExactCurve: ExactCurve(eq, in.X0, in.Y0, in.Xn, exactCurveSamples),
		Results:    make([]ODEResult, 0, len(methods)),
	}

	exact := func(x float64) float64 { return eq.Exact(in.X0, in.Y0, x) }

	for _, m := range methods {
		var res ODEResult
		switch m {
		case "euler":
			curve, hUsed, refs := RungeRefine(EulerMethod, 1, eq.F, in.X0, in.Xn, in.Y0, in.H, in.Epsilon)
			res = ODEResult{
				Method:      "Euler",
				Curve:       curve,
				StepUsed:    hUsed,
				Table:       buildTable(curve, exact),
				MaxError:    MaxAbsError(curve, exact),
				Refinements: refs,
				Note:        "single-step, p=1; step refined by Runge rule",
			}
		case "rk4":
			curve, hUsed, refs := RungeRefine(RK4Method, 4, eq.F, in.X0, in.Xn, in.Y0, in.H, in.Epsilon)
			res = ODEResult{
				Method:      "Runge-Kutta 4",
				Curve:       curve,
				StepUsed:    hUsed,
				Table:       buildTable(curve, exact),
				MaxError:    MaxAbsError(curve, exact),
				Refinements: refs,
				Note:        "single-step, p=4; step refined by Runge rule",
			}
		case "milne":
			curve := MilneMethod(eq.F, in.X0, in.Xn, in.Y0, in.H, in.Epsilon)
			res = ODEResult{
				Method:   "Milne (predictor-corrector)",
				Curve:    curve,
				StepUsed: in.H,
				Table:    buildTable(curve, exact),
				MaxError: MaxAbsError(curve, exact),
				Note:     "multi-step; error vs exact solution",
			}
		}
		resp.Results = append(resp.Results, res)
	}

	return resp, nil
}

func buildTable(curve []numeric.Point, exact func(x float64) float64) []TableRow {
	rows := make([]TableRow, len(curve))
	for i, p := range curve {
		ex := exact(p.X)
		rows[i] = TableRow{I: i, X: p.X, Y: p.Y, Exact: ex, Error: ex - p.Y}
	}
	return rows
}
