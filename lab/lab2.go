package lab

import (
	"fmt"
	"reflect"
)

type T struct {
	A int
	B string
}

func Lab2Command() {
	t := T{23, "skidoo"}
	tt := reflect.TypeOf(t)
	fmt.Printf("t type:%v\n", tt)
	ttp := reflect.TypeOf(&t)
	fmt.Printf("t type:%v\n", ttp)
	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}
