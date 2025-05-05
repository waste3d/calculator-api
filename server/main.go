package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Calculation{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}

type Calculation struct {
	gorm.Model
	ID         string `gorm:"primaryKey" json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

func calculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		log.Fatal("err", err)
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		log.Fatal("err", err)
	}

	return fmt.Sprintf("%v", result), err
}

func getCalculations(ctx echo.Context) error {
	var calculations []Calculation

	if err := db.Find(&calculations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "could not get calculations"})
	}
}

func postCalculations(ctx echo.Context) error {
	var req CalculationRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	result, err := calculateExpression(req.Expression)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: req.Expression,
		Result:     result,
	}

	calculations = append(calculations, calc)
	return ctx.JSON(http.StatusOK, calc)
}

func patchCalculations(ctx echo.Context) error {
	id := ctx.Param("id")

	var req CalculationRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	result, err := calculateExpression(req.Expression)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	for i, calculation := range calculations {
		if calculation.ID == id {
			calculations[i].Expression = req.Expression
			calculations[i].Result = result

			return ctx.JSON(http.StatusOK, calculations[i])
		}
	}
	return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "calculation not found"})
}

func deleteCalculations(ctx echo.Context) error {
	id := ctx.Param("id")

	for i, calculation := range calculations {
		if calculation.ID == id {
			calculations = append(calculations[:i], calculations[i+1:]...)
			return ctx.NoContent(http.StatusNoContent)
		}
	}
	return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "calculation not found"})
}

func main() {

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)
	e.PATCH("/calculations/:id", patchCalculations)
	e.DELETE("/calculations/:id", deleteCalculations)

	e.Start(":8080")

}
