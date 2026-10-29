package main

import (
	"context"
	"fmt"
	"log"

	"github.com/01reportplus/golang-tutorial/internal/pipeline"
)

func main() {
	results, err := pipeline.Map(context.Background(), 3, []int{1, 2, 3, 4, 5},
		func(ctx context.Context, value int) (int, error) {
			if err := ctx.Err(); err != nil {
				return 0, err
			}
			return value * value, nil
		})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(results)
}
