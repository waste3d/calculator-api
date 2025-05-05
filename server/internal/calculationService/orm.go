package calculationservice

import "gorm.io/gorm"

type Calculation struct {
	gorm.Model
	ID         string `gorm:"primaryKey" json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}
