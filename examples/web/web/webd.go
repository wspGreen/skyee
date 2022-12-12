package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/frame"
	"github.com/wspGreen/skyee/iface"
	"github.com/wspGreen/skyee/scontext"
	"github.com/wspGreen/skyee/slog"
)

type ClientConfigData map[string]map[string]map[string]string

type ConfigRespone struct {
	Head         int
	IsSuccess    int
	ClientConfig map[string]string
}

var allconfig ClientConfigData

type ConfigReq struct {
	Ver    string `json:"ConfigVersion"`
	Gameid string `json:"GameID"`
	Head   int    `json:"Head"`
}

type Webd struct {
	a int
}

var Web = NewWeb()

func NewWeb() *Webd {
	return &Webd{a: 0}
}

func (d *Webd) Init(a iface.IActor) {
	skyee.RegisterProtocol(a.GetId(), &scontext.Protocal{
		Name: "socket",
		Type: frame.PTYPE_SOCKET,
		Pack: func(params []interface{}) (rawParams []interface{}) {
			return params
		},
		UnPack: func(rawparams []interface{}) []interface{} {
			return rawparams
		},

		F: func(session uint32, source uint32, cmd string, params []interface{}) []interface{} {
			a.FunCall(cmd, params)
			return nil
		},
	})

	// skyee.Dispatch(
	// 	a.GetId(),
	// 	"cmd",
	// 	func(session uint32, source uint32, cmd string, params []interface{}) []interface{} {
	// 		ret := a.FunCall(cmd, params)

	// 		return skyee.Ret(ret)
	// 	},
	// )

	if allconfig != nil {
		return
	}

	bytes, err := ioutil.ReadFile("./config/AllClientConfig.txt")
	if err != nil {
		slog.Fatal("读取json文件失败 %v", err)
		return
	}
	// c := &ConfigRespone{}
	err = json.Unmarshal(bytes, &allconfig)
	if err != nil {
		slog.Fatal("解析数据失败 %v", err)
		return
	}

	slog.Info("Success Load ClientConfig:%v", allconfig)
}

func (d *Webd) OnServerMessage() {

}

func (d *Webd) OnClientRequest(w http.ResponseWriter, r *http.Request) {

	resp := ConfigRespone{Head: 9, IsSuccess: 0}
	resp.ClientConfig = make(map[string]string)

	defer func() {
		// log.Println("Req :", r.RequestURI)
		// log.Println("Req :", resp)
		jsonbt, _ := json.Marshal(resp)
		w.Write(jsonbt)
	}()

	sign := r.URL.Query().Get("key")
	if sign != "123" {
		// w.Write([]byte("fail"))
		return
	}

	// var conf ConfigReq
	// if err := json.NewDecoder(r.Body).Decode(&conf); err != nil {
	// 	log.Println(err)
	// }

	body, _ := ioutil.ReadAll(r.Body)

	var conf ConfigReq
	if err := json.Unmarshal(body, &conf); err != nil {
		log.Println(err)
		// w.Write([]byte("fail"))
		return
	}

	// jsonbt, _ := json.Marshal(conf)
	// log.Println("Req :", r.RequestURI, sign, string(jsonbt))

	resp.IsSuccess = 1
	clientconf := d.getconfig(conf.Gameid, conf.Ver)
	if clientconf != nil {
		resp.ClientConfig = clientconf
	}
	// fmt.Fprintf(w, "Home Page")
}

func (d *Webd) getconfig(gameid string, ver string) map[string]string {
	gameinfo := allconfig[gameid]
	if gameinfo == nil {
		return nil
	}

	return gameinfo[ver]

}
