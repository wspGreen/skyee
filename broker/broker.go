package broker

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var NATS_ADDR = "nats://127.0.0.1:4222"

var (
	MaxConn       int32  = 50000
	MaxPacketSize uint32 = 1024 * 1024 * 2 // 最大2M的包
)

type BrokerMessage struct {
	Pid     uint32
	Msgtype int8
	Data    []byte
}

func NewBrokerMessage(pid uint32, msgtype int8, data []byte) *BrokerMessage {
	return &BrokerMessage{Pid: pid, Data: data, Msgtype: msgtype}
}

// func (m *BrokerMessage) Pid() uint32 {
// 	return m.pid
// }

// func (m *BrokerMessage) Data() []byte {
// 	return m.data
// }

var DefaultBrokerPack = &BrokerPack{}

type BrokerPack struct {
}

// params = source,cmd,pid,data
// data = dest|source|cmdlen|cmd|pidlen|pid|datalen|data
func (b *BrokerPack) Pack(dest uint32, params ...interface{}) ([]byte, error) {

	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.BigEndian, dest); err != nil {
		return nil, err
	}
	// source
	index := 0
	source, ok := params[index].(uint32)
	cmd := ""
	if ok {
		index++
		cmd = params[index].(string)
	} else {
		cmd = params[index].(string)
		source = 0
	}

	if err := binary.Write(dataBuff, binary.BigEndian, source); err != nil {
		return nil, err
	}

	// cmd
	cmdlen := int16(len(cmd))
	if err := binary.Write(dataBuff, binary.BigEndian, cmdlen); err != nil {
		return nil, err
	}
	if cmdlen > 0 {
		if err := binary.Write(dataBuff, binary.BigEndian, []byte(cmd)); err != nil {
			return nil, err
		}
	}

	// pid
	index++
	pid := []byte(params[index].(string))
	pidlen := int16(len(pid))
	if err := binary.Write(dataBuff, binary.BigEndian, pidlen); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.BigEndian, pid); err != nil {
		return nil, err
	}

	// data
	index++
	data := params[index].([]byte)
	datalen := uint32(len(data))
	if err := binary.Write(dataBuff, binary.BigEndian, datalen); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.BigEndian, data); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// params = {actorid, cmd, pid, data}
// ndata = dest|source|cmdlen|cmd|pid|len|data
func (b *BrokerPack) UnPack(ndata []byte) (params []interface{}, err error) {
	dataBuff := bytes.NewReader(ndata)
	// msg := &BrokerMessage{}

	// dest
	var dest uint32
	if err = binary.Read(dataBuff, binary.BigEndian, &dest); err != nil {
		return
	}

	// source
	var source uint32
	if err = binary.Read(dataBuff, binary.BigEndian, &source); err != nil {
		return
	}

	// cmd
	cmdlen := int16(0)
	if err = binary.Read(dataBuff, binary.BigEndian, &cmdlen); err != nil {
		return
	}
	var cmd string
	if cmdlen > 0 {
		bufcmd := make([]byte, cmdlen)
		if err = binary.Read(dataBuff, binary.BigEndian, &bufcmd); err != nil {
			return
		}
		cmd = string(bufcmd)
	}

	// pid
	pidlen := int16(0)
	if err = binary.Read(dataBuff, binary.BigEndian, &pidlen); err != nil {
		return
	}
	var pid string
	bufpid := make([]byte, pidlen)
	if err = binary.Read(dataBuff, binary.BigEndian, &bufpid); err != nil {
		return
	}
	pid = string(bufpid)

	// data
	var data []byte
	datalen := uint32(0)
	if err = binary.Read(dataBuff, binary.BigEndian, &datalen); err != nil {
		return
	}

	if datalen > MaxPacketSize {
		err = errors.New("too large msg data received")
		return
	}
	data = make([]byte, datalen)
	if err = binary.Read(dataBuff, binary.BigEndian, &data); err != nil {
		return
	}
	params = []interface{}{dest, source, cmd, pid, data}
	// if msg.DataLen > MaxPacketSize {
	// 	return nil, errors.New("too large msg data received")
	// }
	err = nil
	return
}
