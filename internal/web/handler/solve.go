package handler

import (
	"comp-math-6/internal/algo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SolveRequest struct {
	EquationID *int     `json:"equationId" binding:"required"`
	X0         *float64 `json:"x0" binding:"required"`
	Xn         *float64 `json:"xn" binding:"required"`
	Y0         *float64 `json:"y0" binding:"required"`
	H          *float64 `json:"h" binding:"required"`
	Epsilon    *float64 `json:"epsilon" binding:"required"`
	Methods    []string `json:"methods"`
}

func Solve() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SolveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
			return
		}

		resp, err := algo.Solve(algo.SolveInput{
			EquationID: *req.EquationID,
			X0:         *req.X0,
			Xn:         *req.Xn,
			Y0:         *req.Y0,
			H:          *req.H,
			Epsilon:    *req.Epsilon,
			Methods:    req.Methods,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func Equations() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, algo.Equations)
	}
}
