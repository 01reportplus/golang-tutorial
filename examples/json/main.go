package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name string `json:"name"`
	Age  *int   `json:"age,omitempty"`
}

func main() {
	data, err := json.Marshal(User{Name: "Ada", Age: new(0)})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
