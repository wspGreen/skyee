package profile

import (
	"io/ioutil"

	"github.com/bitly/go-simplejson"
	"github.com/wspGreen/skyee/slog"
)

func LoadFile(profilepath string) *simplejson.Json {
	bytes, err := ioutil.ReadFile(profilepath)
	if err != nil {
		slog.Fatal("读取json文件失败 %v", err)
		return nil
	}
	js, err := simplejson.NewJson(bytes)

	if err != nil {
		slog.Fatal("解析数据失败 %v", err)
		return nil
	}
	return js
}
