package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/wspGreen/skyee/log"
)

type ClientConfigData map[string]map[string]map[string]string

type ConfigRespone struct {
	Head         int
	IsSuccess    int
	ClientConfig map[string]string
}

var allconfig ClientConfigData

func init() {
	bytes, err := ioutil.ReadFile("./config/AllClientConfig.txt")
	if err != nil {
		log.Fatal("读取json文件失败", err)
		return
	}
	// c := &ConfigRespone{}
	err = json.Unmarshal(bytes, &allconfig)
	if err != nil {
		log.Fatal("解析数据失败", err)
		return
	}

	log.Println("Success Load ClientConfig:", allconfig)
}

type ConfigReq struct {
	Ver    string `json:"ConfigVersion"`
	Gameid string `json:"GameID"`
	Head   int    `json:"Head"`
}

var Web = NewWeb()

type Webd struct {
}

func NewWeb() *Webd {
	return &Webd{}
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
