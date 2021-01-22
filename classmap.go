package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

func colorToTag(c color.Color) byte {
	r, g, b, a := c.RGBA()
	if r > whiteThreshold && g > whiteThreshold && b > whiteThreshold {
		return white
	} else if r < blackThreshold && g < blackThreshold && b < blackThreshold {
		return black
	} else {
		fmt.Printf("r=%d g=%d b=%d a=%d\n", r, g, b, a)
		panic("目前还不支持此图有别的颜色")
	}
}

//返回上下左右是否有土地
func getDirInfo(m [][]byte, x, y int) (hasUp, hasDown, hasLeft, hasRight, hasCornerUpLeft, hasCornerUpRight, hasCornerDownLeft, hasCornerDownRight bool) {

	minX := 0
	maxX := len(m[0]) - 1

	minY := 0
	maxY := len(m) - 1

	if x < maxX && m[y][x+1] == black {
		hasRight = true
	}
	if x > minX && m[y][x-1] == black {
		hasLeft = true
	}
	if y < maxY && m[y+1][x] == black {
		hasDown = true
	}
	if y > minY && m[y-1][x] == black {
		hasUp = true
	}

	if y > minY && x > minX && m[y-1][x-1] == black {
		hasCornerUpLeft = true
	}

	if y > minY && x < maxX && m[y-1][x+1] == black {
		hasCornerUpRight = true
	}

	if y < maxY && x > minX && m[y+1][x-1] == black {
		hasCornerDownLeft = true
	}

	if y < maxY && x < maxX && m[y+1][x+1] == black {
		hasCornerDownRight = true
	}

	return
}

type classmap struct {
	wCount, hCount          int
	tilewidth, tileheight   int
	margin, spacing         int
	aliasWidth, aliasHeight int
	tilecount, columns      int
	aliasname               string
	classM                  [][]classInfo
}

const (
	white byte = 0
	black      = 1
)

//方便调试 再把tag中的值映射为某个颜色来看看对不对
var tagToColorMap = map[byte]color.RGBA{
	white: {0xff, 0xff, 0x00, 0xff},
	black: {0x00, 0xff, 0xff, 0xff},
}

//给的图并不是绝对的黑和白,加入阈值来规范下
const (
	whiteThreshold = 60000
	blackThreshold = 5000
)

func newClassmap(srcImagePath, aliasPathstring string) *classmap {
	//基本思路 读取像素大略图 =>二维像素集合
	srcImage := filenameToImage(srcImagePath)
	srcRect := srcImage.Bounds()
	wCount := srcRect.Dx()
	hCount := srcRect.Dy()
	m := make([][]byte, hCount)
	for i := 0; i < hCount; i++ {
		m[i] = make([]byte, wCount)
	}

	//从颜色生成 tag数据：就是把颜色值映射为了0 1 便于操作
	for y := srcRect.Min.Y; y < srcRect.Max.Y; y++ {
		for x := srcRect.Min.X; x < srcRect.Max.X; x++ {
			color := srcImage.At(x, y)
			tag := colorToTag(color)
			m[y-srcRect.Min.Y][x-srcRect.Min.X] = tag
		}
	}

	//将tag数据 映射为带边缘信息的 为什么不一次生成？因为它需要知道周边8个位置的信息
	c := make([][]classInfo, hCount)
	for i := 0; i < hCount; i++ {
		c[i] = make([]classInfo, wCount)
	}

	for y := 0; y < hCount; y++ {
		for x := 0; x < wCount; x++ {
			hasUp, hasDown, hasLeft, hasRight, hasCornerUpLeft, hasCornerUpRight, hasCornerDownLeft, hasCornerDownRight := getDirInfo(m, x, y)
			c[y][x] = circleInfoToClass(hasUp, hasDown, hasLeft, hasRight, hasCornerUpLeft, hasCornerUpRight, hasCornerDownLeft, hasCornerDownRight, m[y][x])
		}
	}

	aliasImage := filenameToImage(aliasPathstring)
	aliasname := aliasPathstring
	if strings.ContainsRune(aliasPathstring, '\\') {
		ss := strings.Split(aliasPathstring, "\\")
		aliasname = ss[len(ss)-1]
	}
	aliasRect := aliasImage.Bounds()
	aliasWidth := aliasRect.Dx()
	aliasHeight := aliasRect.Dy()
	tilecount := (aliasWidth / tileWidth) * (aliasHeight / tileHeight)
	columns := (aliasWidth / tileWidth)

	return &classmap{
		wCount:      wCount,
		hCount:      hCount,
		tilewidth:   tileWidth,
		tileheight:  tileHeight,
		margin:      0,
		spacing:     0,
		aliasWidth:  aliasWidth,
		aliasHeight: aliasHeight,
		tilecount:   tilecount,
		columns:     columns,
		aliasname:   aliasname,
		classM:      c,
	}
}

func (t *classmap) toPNG(filename string) {
	dstFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()
	r := image.Rect(0, 0, t.wCount, t.hCount)
	image := image.NewRGBA(r)

	for y := 0; y < t.hCount; y++ {
		for x := 0; x < t.wCount; x++ {
			image.Set(x, y, t.classM[y][x].classInfoToColor())
		}
	}

	png.Encode(dstFile, image)
}

func (t *classmap) genIds() [][]byte {

	result := make([][]byte, t.hCount)
	for i := 0; i < t.hCount; i++ {
		result[i] = make([]byte, t.wCount)
	}
	for y := 0; y < t.hCount; y++ {
		for x := 0; x < t.wCount; x++ {
			result[y][x] = t.classM[y][x].classInfoToID()
		}
	}

	return result
}

func (t *classmap) toTMX(filename string) {
	//写入tmx文件
	wCount := t.wCount
	hCount := t.hCount
	ids := t.genIds()
	strHeader := `<?xml version="1.0" encoding="UTF-8"?>
	<map version="1.4" tiledversion="1.4.3" orientation="orthogonal" renderorder="right-down" width="` + strconv.Itoa(wCount) + `" height="` + strconv.Itoa(hCount) + `" tilewidth="` + strconv.Itoa(t.tilewidth) + `" tileheight="` + strconv.Itoa(t.tileheight) + `" infinite="0" nextlayerid="2" nextobjectid="1">
	 <tileset firstgid="1" name="Terrain Tiles" tilewidth="` + strconv.Itoa(t.tilewidth) + `" tileheight="` + strconv.Itoa(t.tileheight) + `" spacing="` + strconv.Itoa(t.spacing) + `" margin="` + strconv.Itoa(t.margin) + `" tilecount="` + strconv.Itoa(t.tilecount) + `" columns="` + strconv.Itoa(t.columns) + `">
	  <image source="` + t.aliasname + `" width="` + strconv.Itoa(t.aliasWidth) + `" height="` + strconv.Itoa(t.aliasHeight) + `"/>
	 </tileset>
	 <layer id="1" name="图块层 1" width="` + strconv.Itoa(wCount) + `" height="` + strconv.Itoa(hCount) + `">
	  <data encoding="csv">
	  `

	tailer := `</data>
	</layer>
   </map>`
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	bf := bufio.NewWriter(f)
	bf.WriteString(strHeader)
	for y := 0; y < t.hCount; y++ {
		for x := 0; x < t.wCount; x++ {
			r := int(ids[y][x])
			bf.WriteString(strconv.Itoa(r))
			if x != wCount-1 || y != hCount-1 {
				bf.WriteString(",")
			}
		}
		bf.WriteString("\n")
	}
	bf.WriteString(tailer)
	bf.Flush()
}
