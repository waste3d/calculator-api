package calculationservice

import (
	"fmt"
	"log"

	"github.com/Knetic/govaluate"
)

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
