package unittest_test

import (
	"fmt"
	"testing"
)

type IInfo interface {
	CreateRole(name string)
	Run()
}

type BaseInfo struct {
}

func (b *BaseInfo) CreateRole(name string) {
	fmt.Println("BaseInfo CreateRole:", name)
}

func (b *BaseInfo) Run() {
	fmt.Println("BaseInfo Run ")
}

type DogInfo struct {
	BaseInfo
}

func (b *DogInfo) CreateRole(name string) {
	fmt.Println("DogInfo :", name)
	b.BaseInfo.CreateRole(name)
}

func Test_Main(t *testing.T) {
	test2()
}

func test2() {
	var role IInfo = &DogInfo{}
	role.CreateRole("dog")
	role.Run()
}

func test1() {
	var datas map[interface{}]string = make(map[interface{}]string)

	datas["123"] = "123"
	datas[1] = "1"

	fmt.Printf("data : %v\n", datas[1])
	fmt.Printf("data : %v\n", datas["123"])

	update("123")
}

func update(data interface{}) {
	v, ok := data.(int)
	if ok {
		fmt.Println("if : ", v)
	}
	switch v := data.(type) {
	case string:
		fmt.Println("switch :", v)
	}
}
