package unittest_test

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"testing"
)

type User struct {
	Userid int
	Name   string
}

func (u *User) Task0() int {

	fmt.Println("Call Task0 arg is 0 !!!! ")
	return 666
}

func (u *User) Task1(a int) {

	fmt.Println("Call Task1 arg is 1 !!!! ")

}

func Show1(data interface{}) {
	fmt.Println("=============================== ")
	fmt.Println("[data info] ")
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	fmt.Println("Type ", t) // *unittest_test.User
	fmt.Println("Value ", v)
	fmt.Println("Kind ", t.Kind())

}

func Show2(data interface{}) {
	fmt.Println("=============================== ")
	fmt.Println("[方法索引调用]")
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	fmt.Println("method Num ", t.NumMethod())
	// var m reflect.Method
	for i := 0; i < v.NumMethod(); i++ {

		fmt.Printf("method Name:%s Type:%s\n", t.Method(i).Name, v.Method(i).Type())
		m := v.Method(i)
		// 通过反射调用方法传递的参数必须是 []reflect.Value 类型
		var args = []reflect.Value{reflect.ValueOf(1)}
		paramenum := v.Method(i).Type().NumIn()
		var ret []reflect.Value
		if paramenum > 0 {
			ret = m.Call(args)
		} else {
			ret = m.Call(nil)
		}

		if len(ret) > 0 {
			intf := ret[0].Interface()
			fmt.Printf("%s Return -> %s \n", t.Method(i).Name, intf)
		}
	}
}

func Show3(data interface{}) {
	fmt.Println("=============================== ")
	fmt.Println("[方法名调用]")
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	m := v.MethodByName("Task0") // var m reflect.Value
	m.Call(nil)

	m1, ok := t.MethodByName("Task0") // var m1 reflect.Method
	if !ok {
		panic("method no exist")
	}
	var args = []reflect.Value{v}
	m1.Func.Call(args)
}

func Start() {
	user01 := &User{
		Userid: 123,
		Name:   "mark",
	}
	// reflect.Struct
	// fmt.Println("name ", reflect.TypeOf())
	Show1(user01)
	Show2(user01)
	Show3(user01)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("Handle message panic: %+v\n%s", err, debug.Stack()))
		}
	}()
}

func reflectNew() {

	createObj((*User)(nil))
}

func createObj(u interface{}) *User {
	typ := reflect.TypeOf(u)
	if typ.Kind() == reflect.Ptr { //指针类型获取真正type需要调用Elem
		typ = typ.Elem()
	}
	fVal := reflect.New(typ)
	fVal.Elem().Field(0).SetInt(20)
	fVal.Elem().Field(1).SetString("Number")
	user := fVal.Elem().Interface().(User)
	fmt.Printf("userid:%v name:%v \n", user.Userid, user.Name)
	return &user
}

func TestRef(t *testing.T) {
	// Start()
	reflectNew()

}
