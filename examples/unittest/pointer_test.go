package unittest_test

import (
	"fmt"
	"testing"
	"time"
)

// info 引用了map 2的指针，2被删除后 info仍然可以使用

type IUser interface {
	GetId() int
}

type UserInfo struct {
	Id int
}

func (u *UserInfo) GetId() int {
	return u.Id
}

var infos = map[int]*UserInfo{}
var info IUser

func TestXxx(t *testing.T) {
	infos[1] = &UserInfo{1}
	infos[2] = &UserInfo{2}
	info = infos[2]

	delete(infos, 2)

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		go func() {

			fmt.Println(" id ", info.(*UserInfo).Id)
		}()
	}

	select {}
}
