package lab

import (
	"fmt"
)

type HumanAction interface {
	Eat()
	Walk()
	Fly()
}

type Human struct {
	name string
	sex  string
}

func (h Human) Eat() {
	fmt.Println("Human.Eat:", h)
}

func (h Human) Walk() {
	fmt.Println("Human.Walk:", h)
}

type SuperMan struct {
	name  string
	level int
	Human
	sex string
}

func (s SuperMan) Eat() {
	fmt.Println("SuperMan.Eat:", s)
}

func (s SuperMan) Fly() {
	fmt.Println("I believe I can fly!", s)
}

func test(h Human) {
	fmt.Println("test pass!")
}

func Lab1Command() {
	h := Human{"human_name", "human_sex"}
	h.Eat()  //调用自己的方法
	h.Walk() //调用自己的方法

	/*
	 *s := SuperMan{Human{"parent_name", "parent_sex"}, "sm_name", 99}
	 */
	s := SuperMan{"sm_name", 99, Human{"parent_name", "parent_sex"}, "sm_sex"}

	fmt.Println(s.name)       //myname
	fmt.Println(s.Human.name) //sm
	fmt.Println(s.sex)        //no
	fmt.Println(s.Human.sex)  //no

	s.sex = "new sex"
	fmt.Println(s.sex)       //no
	fmt.Println(s.Human.sex) //no

	s.Walk()      //调用父类中的方法
	s.Eat()       //调用自己重写的父类的方法
	s.Human.Eat() //调用父类中被自己重写的方法，这点比java高级
	s.Fly()       //调用自己的方法

	test(h)
	//test(s)//导致错误，方法需要Human类型，但传的是SuperMan类型
	test(s.Human) //正确
	/*
	 *var a HumanAction = &h
	 *a.Eat()
	 */
}
