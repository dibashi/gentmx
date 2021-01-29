package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"math/rand"
	"os"
)

type classInfo byte

const (
	you classInfo = iota
	wu
	up
	down
	left
	right
	upleft
	upright
	downleft
	downright
	cornerupleft
	cornerupright
	cornerdownleft
	cornerdownright
)

//八方向有土地吗
func circleInfoToClass(hasUp, hasDown, hasLeft, hasRight, hasCornerUpLeft, hasCornerUpRight, hasCornerDownLeft, hasCornerDownRight bool, tag byte) classInfo {
	if tag == white {
		return wu
	}

	if !hasUp && !hasLeft {
		return upleft
	}
	if !hasUp && !hasRight {
		return upright
	}

	if !hasDown && !hasLeft {
		return downleft
	}
	if !hasDown && !hasRight {
		return downright
	}

	if !hasUp && hasRight && hasLeft {
		return up
	}

	if !hasDown && hasRight && hasLeft {
		return down
	}

	if !hasLeft && hasUp && hasDown {
		return left
	}

	if !hasRight && hasUp && hasDown {
		return right
	}

	//要看看角部分 通过才是有 否则要填充角
	if !hasCornerUpLeft && hasUp && hasLeft {
		return cornerupleft
	}

	if !hasCornerUpRight && hasUp && hasRight {
		return cornerupright
	}

	if !hasCornerDownLeft && hasDown && hasLeft {
		return cornerdownleft
	}

	if !hasCornerDownRight && hasDown && hasRight {
		return cornerdownright
	}

	return you
}

var classToColorMap = map[classInfo]color.RGBA{
	you: {0x00, 0x00, 0x00, 0xff},
	wu:  {0xff, 0xff, 0xff, 0xff},

	up:    {0xff, 0x00, 0x00, 0xff},
	down:  {0x00, 0xff, 0x00, 0xff},
	left:  {0x00, 0x00, 0xff, 0xff},
	right: {0x33, 0x66, 0x99, 0xff},

	upleft:    {0xff, 0xff, 0x00, 0xff},
	upright:   {0xff, 0x00, 0xff, 0xff},
	downleft:  {0x00, 0xff, 0xff, 0xff},
	downright: {0x99, 0x66, 0x33, 0xff},
}

//将classM 映射为ids 需要一个映射表来做映射
var classToIDMap = map[classInfo][]byte{}
var youWeightsTrue []float64

var tileWidth int
var tileHeight int

type gentmxjsoninfo struct {
	TileWidth, TileHeight int
	You                   []byte
	YouWeights            []byte
	YouWeightsTrue        []float64
	Wu                    []byte
	//以下是边缘 名字代表他什么方位没有东西
	Up              []byte
	Down            []byte
	Left            []byte
	Right           []byte
	Upleft          []byte
	Upright         []byte
	Downleft        []byte
	Downright       []byte
	CornerUpleft    []byte
	CornerUpright   []byte
	CornerDownleft  []byte
	CornerDownright []byte
}

//GenClassToIDMap 给定一个json文件来创建 classinfo 到 图块id的映射
func GenClassToIDMap(jsonfile string) {
	jsonf, err := os.Open(jsonfile)
	if err != nil {
		panic(err)
	}

	defer jsonf.Close()

	jsonBytes, err := ioutil.ReadAll(jsonf)
	if err != nil {
		panic(err)
	}

	var cm gentmxjsoninfo
	ej := json.Unmarshal(jsonBytes, &cm)
	if ej != nil {
		panic(ej)
	}

	tileWidth = cm.TileWidth
	tileHeight = cm.TileHeight

	fmt.Println("tileWidth:", tileWidth, "tileHeight:", tileHeight)

	classToIDMap[you] = cm.You
	classToIDMap[wu] = cm.Wu

	classToIDMap[up] = cm.Up
	classToIDMap[down] = cm.Down
	classToIDMap[left] = cm.Left
	classToIDMap[right] = cm.Right

	classToIDMap[upleft] = cm.Upleft
	classToIDMap[upright] = cm.Upright
	classToIDMap[downleft] = cm.Downleft
	classToIDMap[downright] = cm.Downright

	classToIDMap[cornerupleft] = cm.CornerUpleft
	classToIDMap[cornerupright] = cm.CornerUpright
	classToIDMap[cornerdownleft] = cm.CornerDownleft
	classToIDMap[cornerdownright] = cm.CornerDownright

	totalWeight := 0
	for _, w := range cm.YouWeights {
		totalWeight += int(w)
	}
	tmpWeight := 0
	for _, v := range cm.YouWeights {
		tmpWeight += int(v)
		cm.YouWeightsTrue = append(cm.YouWeightsTrue, float64(tmpWeight)/float64(totalWeight))
		fmt.Println(tmpWeight, totalWeight, float64(tmpWeight)/float64(totalWeight))
	}
	youWeightsTrue = cm.YouWeightsTrue
	fmt.Println(cm.YouWeightsTrue, len(cm.YouWeightsTrue))

}

func (ci classInfo) classInfoToID() byte {
	if ci == you {
		r := rand.Float64()
		for i := 0; i < len(youWeightsTrue); i++ {
			if r < youWeightsTrue[i] {
				return classToIDMap[ci][i]
			}
		}

		panic("不可能")
	}
	randIndex := rand.Intn(len(classToIDMap[ci]))
	return classToIDMap[ci][randIndex]
}
func (ci classInfo) classInfoToColor() color.RGBA {
	return classToColorMap[ci]
}
