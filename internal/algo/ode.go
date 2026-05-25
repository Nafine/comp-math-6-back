package algo

import (
	"comp-math-6/internal/numeric"
	"math"
)

type RHS func(x, y float64) float64

func gridSize(x0, xn, h float64) int {
	n := int(math.Round((xn - x0) / h))
	if n < 1 {
		n = 1
	}
	return n
}

func EulerMethod(f RHS, x0, xn, y0, h float64) []numeric.Point {
	n := gridSize(x0, xn, h)
	hAdj := (xn - x0) / float64(n)
	out := make([]numeric.Point, 0, n+1)
	x, y := x0, y0
	out = append(out, numeric.Point{X: x, Y: y})
	for i := 0; i < n; i++ {
		y = y + hAdj*f(x, y)
		x = x0 + float64(i+1)*hAdj
		out = append(out, numeric.Point{X: x, Y: y})
	}
	return out
}

func RK4Method(f RHS, x0, xn, y0, h float64) []numeric.Point {
	n := gridSize(x0, xn, h)
	hAdj := (xn - x0) / float64(n)
	out := make([]numeric.Point, 0, n+1)
	x, y := x0, y0
	out = append(out, numeric.Point{X: x, Y: y})
	for i := 0; i < n; i++ {
		k1 := hAdj * f(x, y)
		k2 := hAdj * f(x+hAdj/2, y+k1/2)
		k3 := hAdj * f(x+hAdj/2, y+k2/2)
		k4 := hAdj * f(x+hAdj, y+k3)
		y = y + (k1+2*k2+2*k3+k4)/6
		x = x0 + float64(i+1)*hAdj
		out = append(out, numeric.Point{X: x, Y: y})
	}
	return out
}

const milneMaxIter = 50

func MilneMethod(f RHS, x0, xn, y0, h, eps float64) []numeric.Point {
	n := gridSize(x0, xn, h)
	hAdj := (xn - x0) / float64(n)

	if n < 4 {
		return RK4Method(f, x0, xn, y0, h)
	}

	starter := RK4Method(f, x0, x0+3*hAdj, y0, hAdj)
	ys := make([]float64, n+1)
	xs := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		xs[i] = x0 + float64(i)*hAdj
	}
	for i := 0; i < 4; i++ {
		ys[i] = starter[i].Y
	}
	fs := make([]float64, n+1)
	for i := 0; i < 4; i++ {
		fs[i] = f(xs[i], ys[i])
	}

	for i := 4; i <= n; i++ {
		yPred := ys[i-4] + (4*hAdj/3)*(2*fs[i-3]-fs[i-2]+2*fs[i-1])
		yCorr := yPred
		for it := 0; it < milneMaxIter; it++ {
			fPred := f(xs[i], yPred)
			yCorr = ys[i-2] + (hAdj/3)*(fs[i-2]+4*fs[i-1]+fPred)
			if math.Abs(yCorr-yPred) <= eps {
				break
			}
			yPred = yCorr
		}
		ys[i] = yCorr
		fs[i] = f(xs[i], yCorr)
	}

	out := make([]numeric.Point, n+1)
	for i := 0; i <= n; i++ {
		out[i] = numeric.Point{X: xs[i], Y: ys[i]}
	}
	return out
}

const rungeMaxRefinements = 20

type SingleStepSolver func(f RHS, x0, xn, y0, h float64) []numeric.Point

func RungeRefine(method SingleStepSolver, p int, f RHS, x0, xn, y0, h0, eps float64) (curve []numeric.Point, finalH float64, refinements int) {
	denom := math.Pow(2, float64(p)) - 1
	h := h0
	curveH := method(f, x0, xn, y0, h)
	for refinements = 0; refinements < rungeMaxRefinements; refinements++ {
		curveH2 := method(f, x0, xn, y0, h/2)
		yEndH := curveH[len(curveH)-1].Y
		yEndH2 := curveH2[len(curveH2)-1].Y
		if math.Abs(yEndH-yEndH2)/denom <= eps {
			return curveH2, h / 2, refinements + 1
		}
		h /= 2
		curveH = curveH2
	}
	return curveH, h, refinements
}

func ExactCurve(eq Equation, x0, y0, xn float64, samples int) []numeric.Point {
	if samples < 2 {
		samples = 2
	}
	out := make([]numeric.Point, samples+1)
	step := (xn - x0) / float64(samples)
	for i := 0; i <= samples; i++ {
		x := x0 + float64(i)*step
		out[i] = numeric.Point{X: x, Y: eq.Exact(x0, y0, x)}
	}
	return out
}

func MaxAbsError(approx []numeric.Point, exact func(x float64) float64) float64 {
	var maxErr float64
	for _, p := range approx {
		e := math.Abs(exact(p.X) - p.Y)
		if e > maxErr {
			maxErr = e
		}
	}
	return maxErr
}
