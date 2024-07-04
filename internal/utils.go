package internal

import (
	"fmt"
	"strconv"
)

func parseInt(s string) *int {
	intVal, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf(fmt.Sprintf("Error converting string %s to integer:", s), err)
		return nil
	}
	return &intVal
}
