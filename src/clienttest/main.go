// testclient project testclient.go

package main

import (
	"encoding/json"
	"flag"
	"net/url"
	"strconv"
	"sync"
	"time"

	//"os"
	//"os/signal"
	//"io"
	"fmt"
	//	"net"
	//"dq/rpc"
	//"time"
	//"net/rpc/jsonrpc"
	"dq/datamsg"

	"math/rand"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "www.game5868.top:443", "http service address")

//var addr = flag.String("addr", "127.0.0.1:1117", "http service address")

func main() {

	fmt.Println("start!!")
	var waitg sync.WaitGroup
	for j := 0; j < 5000; j++ {
		waitg.Add(1)

		go func() {

			client(strconv.Itoa(j))
			waitg.Done()
		}()

		time.Sleep(time.Millisecond * 20)
	}

	waitg.Wait()
	fmt.Println("over!!")

}

//ModeType  string
//	Uid       int
//	MsgId     int
//	MsgType   string
//	ConnectId int
//	JsonData  string

func client(id string) {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/connect"}
	fmt.Println("start game to ", id)

	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

	lock := new(sync.RWMutex)

	lock.Lock()
	c.WriteMessage(websocket.TextMessage, CS_MsgQuickLogin(id))
	lock.Unlock()

	go func() {
		for {
			lock.Lock()
			c.WriteMessage(websocket.TextMessage, CS_Heart())
			lock.Unlock()
			time.Sleep(time.Second * 3)
		}
	}()

	var myUid = -1
	var myInfo = datamsg.MsgGame5GPlayerInfo{}
	var gameInfo = datamsg.MsgGame5GInfo{}
	//var qizipos = make([]int, 0)
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		h1 := &datamsg.MsgBase{}
		err = json.Unmarshal(data, h1)
		if err != nil {
			fmt.Println("--error")
			break
		} else {

			//登录成功
			if h1.MsgType == "SC_LoginResponse" {
				time.Sleep(time.Second * 2)
				lock.Lock()
				c.WriteMessage(websocket.TextMessage, CS_QuickGame())
				lock.Unlock()

			} else if h1.MsgType == "SC_NewGame" { //游戏创建好了 等待加入
				h2 := &datamsg.SC_NewGame{}
				err := json.Unmarshal([]byte(h1.JsonData), h2)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				lock.Lock()
				c.WriteMessage(websocket.TextMessage, CS_GoIn(h2.GameId))
				lock.Unlock()

			} else if h1.MsgType == "SC_GameInfo" {
				h2 := &datamsg.SC_GameInfo{}
				err := json.Unmarshal([]byte(h1.JsonData), h2)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				for _, v := range h2.PlayerInfo {
					if v.Uid == myUid {
						myInfo = v
					}
				}

				gameInfo = h2.GameInfo

			} else if h1.MsgType == "SC_MsgHallInfo" {
				h2 := &datamsg.SC_MsgHallInfo{}
				err := json.Unmarshal([]byte(h1.JsonData), h2)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				myUid = h2.PlayerInfo.Uid

			} else if h1.MsgType == "SC_ChangeGameTurn" {
				h2 := &datamsg.SC_ChangeGameTurn{}
				err := json.Unmarshal([]byte(h1.JsonData), h2)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				//轮到我了
				if h2.GameSeatIndex == myInfo.SeatIndex {
					time.Sleep(time.Second * 2)
					x := -1
					y := -1
					for {
						x = rand.Intn(15)
						y = rand.Intn(15)
						if gameInfo.QiPan[y][x] < 0 {
							lock.Lock()
							c.WriteMessage(websocket.TextMessage, CS_DoGame5G(x, y))
							lock.Unlock()
							break
						}
					}

				}

			} else if h1.MsgType == "SC_DoGame5G" {

				h2 := &datamsg.SC_DoGame5G{}
				err := json.Unmarshal([]byte(h1.JsonData), h2)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				gameInfo.QiPan[h2.Y][h2.X] = h2.GameSeatIndex

			} else if h1.MsgType == "SC_GameOver" {
				h2 := &datamsg.SC_GameOver{}
				err := json.Unmarshal([]byte(h1.JsonData), h2)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				time.Sleep(time.Second * 2)
				lock.Lock()
				c.WriteMessage(websocket.TextMessage, CS_QuickGame())
				lock.Unlock()

			} else if h1.MsgType == "SC_LoginResponse" {

			}

		}

		//fmt.Println("recv: %s", message)
	}

	c.Close() //关闭连接
}

func msgBase(modeType string, msgType string) *datamsg.MsgBase {
	data := &datamsg.MsgBase{}
	data.ModeType = modeType
	data.MsgType = msgType
	return data
}

//快速登录
func CS_MsgQuickLogin(id string) []byte {
	data := msgBase("Login", "CS_MsgQuickLogin")

	jd := &datamsg.CS_MsgQuickLogin{}
	jd.Platform = "ios"
	jd.MachineId = "android_" + id

	jdbytes, _ := json.Marshal(jd)
	data.JsonData = string(jdbytes)

	data1, err1 := json.Marshal(data)
	if err1 == nil {
		return data1
	}

	return []byte("")
}

//快速游戏
func CS_QuickGame() []byte {
	data := msgBase("Hall", "CS_QuickGame")

	data1, err1 := json.Marshal(data)
	if err1 == nil {
		return data1
	}

	return []byte("")
}

//心跳
func CS_Heart() []byte {
	data := msgBase("Hall", "CS_Heart")

	data1, err1 := json.Marshal(data)
	if err1 == nil {
		return data1
	}

	return []byte("")
}

//进入游戏
func CS_GoIn(id int) []byte {
	data := msgBase("Game5G", "CS_GoIn")

	jd := &datamsg.CS_GoIn{}
	jd.GameId = id

	jdbytes, _ := json.Marshal(jd)
	data.JsonData = string(jdbytes)
	data1, err1 := json.Marshal(data)
	if err1 == nil {
		return data1
	}

	return []byte("")
}

//走棋
func CS_DoGame5G(x int, y int) []byte {
	data := msgBase("Game5G", "CS_DoGame5G")

	jd := &datamsg.CS_DoGame5G{}
	jd.X = x
	jd.Y = y

	jdbytes, _ := json.Marshal(jd)
	data.JsonData = string(jdbytes)
	data1, err1 := json.Marshal(data)
	if err1 == nil {
		return data1
	}

	return []byte("")
}
