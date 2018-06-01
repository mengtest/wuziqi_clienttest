// testclient project testclient.go
package main

import (
	"fmt"
)

type AI struct {
	Qipan      [15][15]int
	Myseat     int
	PlayerSeat int
}
type Point struct {
	x int
	y int
}

func (ai *AI) Evaluate(x int, y int, myseat int) int {

	p := Point{x: x, y: y}

	return ai.eva(p, myseat, 1) + ai.eva(p, myseat, 2)
}

func (ai *AI) eva(p Point, me int, plyer int) int { // me:我的代号  plyer:当前计算的player的代号
	value := 0
	numoftwo := 0
	for i := 1; i <= 8; i++ { // 8个方向
		// 活四 01111* *代表当前空位置  0代表其他空位置    下同
		if ai.getLine(p, i, -1) == plyer && ai.getLine(p, i, -2) == plyer &&
			ai.getLine(p, i, -3) == plyer && ai.getLine(p, i, -4) == plyer && ai.getLine(p, i, -5) == 0 {
			value += 300000
			if me != plyer {
				value -= 500
			}
			continue
		}
		// 死四A 21111*
		if ai.getLine(p, i, -1) == plyer && ai.getLine(p, i, -2) == plyer &&
			ai.getLine(p, i, -3) == plyer && ai.getLine(p, i, -4) == plyer &&
			(ai.getLine(p, i, -5) == 3-plyer || ai.getLine(p, i, -5) == -1) {
			value += 250000
			if me != plyer {
				value -= 500
			}
			continue
		}
		// 死四B 111*1
		if ai.getLine(p, i, -1) == plyer && ai.getLine(p, i, -2) == plyer && ai.getLine(p, i, -3) == plyer && ai.getLine(p, i, 1) == plyer {
			value += 240000
			if me != plyer {
				value -= 500
			}
			continue
		}
		// 死四C 11*11
		if ai.getLine(p, i, -1) == plyer && ai.getLine(p, i, -2) == plyer && ai.getLine(p, i, 1) == plyer && ai.getLine(p, i, 2) == plyer {
			value += 230000
			if me != plyer {
				value -= 500
			}
			continue
		}
		// 活三 近3位置 111*0
		if ai.getLine(p, i, -1) == plyer && ai.getLine(p, i, -2) == plyer && ai.getLine(p, i, -3) == plyer {
			if ai.getLine(p, i, 1) == 0 {
				value += 750
				if ai.getLine(p, i, -4) == 0 {
					value += 3150
					if me != plyer {
						value -= 300
					}
				}
			}
			if (ai.getLine(p, i, 1) == 3-plyer || ai.getLine(p, i, 1) == -1) && ai.getLine(p, i, -4) == 0 {
				value += 500
			}
			continue
		}
		// 活三 远3位置 1110*
		if ai.getLine(p, i, -1) == 0 && ai.getLine(p, i, -2) == plyer && ai.getLine(p, i, -3) == plyer && ai.getLine(p, i, -4) == plyer {
			value += 350
			continue
		}
		// 死三 11*1
		if ai.getLine(p, i, -1) == plyer && ai.getLine(p, i, -2) == plyer && ai.getLine(p, i, 1) == plyer {
			value += 600
			if ai.getLine(p, i, -3) == 0 && ai.getLine(p, i, 2) == 0 {
				value += 3150
				continue
			}
			if (ai.getLine(p, i, -3) == 3-plyer || ai.getLine(p, i, -3) == -1) && (ai.getLine(p, i, 2) == 3-plyer || ai.getLine(p, i, 2) == -1) {
				continue
			} else {
				value += 700
				continue
			}
		}
		//活二的个数
		if ai.getLine(p, i, -1) == plyer && ai.getLine(p, i, -2) == plyer && ai.getLine(p, i, -3) != 3-plyer && ai.getLine(p, i, 1) != 3-plyer {
			numoftwo++
		}
		//其余散棋
		numOfplyer := 0            // 因为方向会算两次？
		for k := -4; k <= 0; k++ { // ++++* +++*+ ++*++ +*+++ *++++
			temp := 0
			for l := 0; l <= 4; l++ {
				if ai.getLine(p, i, k+l) == plyer {
					temp++
				} else if ai.getLine(p, i, k+l) == 3-plyer || ai.getLine(p, i, k+l) == -1 {
					temp = 0
					break
				}
			}
			numOfplyer += temp
		}
		value += numOfplyer * 15
		if numOfplyer != 0 {

		}
	}
	if numoftwo >= 2 {
		value += 3000
		if me != plyer {
			value -= 100
		}
	}
	if value < 0 {
		fmt.Println("--socre:%d", value)
	}
	//
	return value
}

func (ai *AI) getLine(p Point, i int, j int) int {
	x := p.x
	y := p.y
	switch i {
	case 1:
		x = x + j

	case 2:
		x = x + j
		y = y + j

	case 3:
		y = y + j

	case 4:
		x = x - j
		y = y + j

	case 5:
		x = x - j

	case 6:
		x = x - j
		y = y - j

	case 7:
		y = y - j

	case 8:
		x = x + j
		y = y - j
	}
	if x < 0 || y < 0 || x > 14 || y > 14 { // 越界处理
		return -1
	}
	return ai.Qipan[y][x]
}
