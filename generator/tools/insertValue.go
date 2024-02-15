package tools

import "fmt"

func InsertValue[T any](sl *[]T, value T, maxLength int) {
	if len(*sl) < cap(*sl) {
		*sl = append(*sl, value)
	} else {
		for i := 0; i < len(*sl)-1; i++ {
			fmt.Println(i)
			(*sl)[i] = (*sl)[i+1]
		}
		(*sl)[len(*sl)-1] = value
	}
}
