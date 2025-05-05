package main

import (
	"log"
	"net/http"

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

// Methods - Find, Create, Update, Delete

func getCalculations(ctx echo.Context) error {
	var calculations []Calculation

	if err := db.Find(&calculations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "could not get calculations"})
	}
	return ctx.JSON(http.StatusOK, calculations)
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

	if err := db.Create(&calc).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "could not add calculation"})
	}
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

	var calc Calculation
	if err := db.Find(&calc, "id=?", id).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Could not find expression"})
	}

	calc.Expression = req.Expression
	calc.Result = result

	if err := db.Save(&calc).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "could not update calculation"})
	}

	return ctx.JSON(http.StatusOK, calc)
}

func deleteCalculations(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := db.Delete(&Calculation{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "could not delete calculation"})

	}

	return ctx.NoContent(http.StatusNoContent)
}

func main() {
	initDB()

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)
	e.PATCH("/calculations/:id", patchCalculations)
	e.DELETE("/calculations/:id", deleteCalculations)

	e.Start(":8080")

}
