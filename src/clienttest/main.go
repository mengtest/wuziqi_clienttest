// testclient project testclient.go

package main

import (
	"encoding/json"
	"flag"
	"net/url"
	"strconv"
	"sync"
	"time"

	"fmt"

	"dq/datamsg"

	"math/rand"

	"github.com/gorilla/websocket"
)

//var addr = flag.String("addr", "www.game5868.top:443", "http service address")

var addr = flag.String("addr", "127.0.0.1:1117", "http service address")

func main() {

	//	score1 := make([]ScorePoint, 5)

	//	pushScore(&score1, ScorePoint{Score: 5, X: 10, Y: 10})
	//	pushScore(&score1, ScorePoint{Score: 15, X: 11, Y: 10})
	//	pushScore(&score1, ScorePoint{Score: 25, X: 12, Y: 10})
	//	pushScore(&score1, ScorePoint{Score: 35, X: 13, Y: 10})
	//	pushScore(&score1, ScorePoint{Score: 51, X: 14, Y: 10})
	//	pushScore(&score1, ScorePoint{Score: 14, X: 15, Y: 10})
	//	pushScore(&score1, ScorePoint{Score: 32, X: 16, Y: 10})
	//	pushScore(&score1, ScorePoint{Score: 18, X: 17, Y: 10})

	//	for k, v := range score1 {
	//		fmt.Println("----")
	//		fmt.Println(k)
	//		fmt.Println(v.Score)
	//		fmt.Println(v.X)
	//		fmt.Println("----")
	//	}
	//	time.Sleep(time.Millisecond * 20000)

	fmt.Println("start!!")
	var waitg sync.WaitGroup
	for j := 0; j < 1; j++ {
		waitg.Add(1)

		go func() {

			client(strconv.Itoa(j))
			waitg.Done()
		}()

		time.Sleep(time.Millisecond * 100)
	}

	waitg.Wait()
	fmt.Println("over!!")

}

type ScorePoint struct {
	Score int
	X     int
	Y     int
}

func pushScore(score *[]ScorePoint, cursocre ScorePoint) {
	for k := 0; k < len(*score); k++ {
		if cursocre.Score > (*score)[k].Score || (rand.Intn(4) == 0 && cursocre.Score == (*score)[k].Score) {
			for k1 := len(*score) - 1; k1 > k; k1-- {
				if k1-1 >= 0 {
					(*score)[k1] = (*score)[k1-1]
				}
			}
			(*score)[k] = cursocre
			break
		}
	}
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

	//clientId, _ := strconv.Atoi(id)
	//nengli := rand.Intn(10) + 2
	//nengli := 2

	var myUid = -1
	var myInfo = datamsg.MsgGame5GPlayerInfo{}
	var gameInfo = datamsg.MsgGame5GInfo{}

	var myai = AI{Myseat: 1, PlayerSeat: 2}

	var isfirst = true

	var gonum = 0 //步数

	//	type AI struct {
	//	Qipan      [][]int
	//	Myseat     int
	//	PlayerSeat int
	//}

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

						//myai.Myseat = myInfo.SeatIndex
					}
				}

				gameInfo = h2.GameInfo

				for y := 0; y < 15; y++ {
					for x := 0; x < 15; x++ {
						myai.Qipan[y][x] = gameInfo.QiPan[y][x] + 1
					}
				}

				//myai.Qipan = gameInfo.QiPan
				//

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

					x := -1
					y := -1
					//
					if isfirst == true {

						sleeptime := rand.Intn(2) + 2
						time.Sleep(time.Second * time.Duration(sleeptime))

						x = rand.Intn(11) + 2
						y = rand.Intn(11) + 2

						if gameInfo.QiPan[y][x] < 0 {
							lock.Lock()
							c.WriteMessage(websocket.TextMessage, CS_DoGame5G(x, y))
							lock.Unlock()
						}
						isfirst = false

					} else {
						//前面几步快速下棋
						if gonum < 6 {
							sleeptime := rand.Intn(2) + 2
							time.Sleep(time.Second * time.Duration(sleeptime))
						} else if gonum < 12 {
							sleeptime := rand.Intn(3) + 2
							time.Sleep(time.Second * time.Duration(sleeptime))
						} else if gonum < 24 {
							sleeptime := rand.Intn(4) + 2
							time.Sleep(time.Second * time.Duration(sleeptime))
						} else if gonum < 40 {
							sleeptime := rand.Intn(5) + 2
							time.Sleep(time.Second * time.Duration(sleeptime))
						} else {
							sleeptime := rand.Intn(10) + 2
							time.Sleep(time.Second * time.Duration(sleeptime))
						}

						for {
							//maxScore := -1

							//score1 := make([]int, 5)
							score1 := make([]ScorePoint, 5)
							scoreMaxNum := 0
							for y1 := 0; y1 < 15; y1++ {
								for x1 := 0; x1 < 15; x1++ {
									if myai.Qipan[y1][x1] == 0 {
										score := myai.Evaluate(x1, y1, myInfo.SeatIndex+1)

										pushScore(&score1, ScorePoint{Score: score, X: x1, Y: y1})
										scoreMaxNum++

										//										if score > maxScore {
										//											maxScore = score
										//											x = x1
										//											y = y1
										//										} else if score == maxScore {
										//											if rand.Intn(4) == 0 {
										//												maxScore = score
										//												x = x1
										//												y = y1
										//											}
										//										}
									}
								}
							}

							fmt.Println("score0:", score1[0].Score, "-----score2:", score1[2].Score, "-----score4:", score1[4].Score)

							//if rand.Intn(nengli) == 0 && scoreMaxNum >= 5 {
							if scoreMaxNum >= 5 {
								//								x = score1[rand.Intn(4)+1].X
								//								y = score1[rand.Intn(4)+1].Y
								//r := rand.Intn(2) + 2
								if score1[0].Score < 1000 {
									r11 := 4
									x = score1[r11].X
									y = score1[r11].Y
									fmt.Println("score1 do:", score1[r11].Score)
								} else if rand.Intn(7) == 0 {
									r11 := 4
									x = score1[r11].X
									y = score1[r11].Y
									fmt.Println("score2 do:", score1[r11].Score)
								} else {
									x = score1[0].X
									y = score1[0].Y
								}

							} else {
								x = score1[0].X
								y = score1[0].Y
							}

							if gameInfo.QiPan[y][x] < 0 {
								lock.Lock()
								c.WriteMessage(websocket.TextMessage, CS_DoGame5G(x, y))
								lock.Unlock()
								break
							}
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
				myai.Qipan[h2.Y][h2.X] = h2.GameSeatIndex + 1
				isfirst = false
				gonum++

			} else if h1.MsgType == "SC_GameOver" {
				h2 := &datamsg.SC_GameOver{}
				err := json.Unmarshal([]byte(h1.JsonData), h2)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				isfirst = true
				gonum = 0
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

var names = []string{"粉笔超人", "走", "灯火阑珊", "wing", "v爱阅读", "华云", "小朋友", "jessica", "稳稳", "舒静模",
	"何晓芳", "自醉自演", "妞妈", "小雅", "小鹅子", "孟茜茜", "深蓝李子", "蒲公英", "小丸子", "patient",
	"烂菜叶", "燕妮", "木木妈", "珍惜现在", "叶赫", "小甘甘", "临西", "报单手", "岁月如歌", "一缕阳光",
	"缘来如此", "淡定", "美好生活", "浮云", "荟", "灿烂人生", "静心倾听", "小玉", "小雅", "敏敏",
	"痛", "十一先生", "无所谓", "ablajan", "顺其自然", "草原花", "蓝心", "小柠檬", "静心好", "似水流年"}

//快速登录
func CS_MsgQuickLogin(id string) []byte {
	data := msgBase("Login", "CS_MsgQuickLogin")

	jd := &datamsg.CS_MsgQuickLogin{}
	jd.Platform = "android"
	jd.MachineId = "android_" + id

	idint, _ := strconv.Atoi(id)
	jd.Name = names[idint%50]
	jd.Avatar = "resources/head/head" + strconv.Itoa(idint%50+1) + ".jpg"

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
