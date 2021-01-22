package main

import (
	"fmt"
	"os"
)

//核心任务:根据黑白图 生成tmx文件
//白色为河流
//黑色分为周边有无白色：
//若黑周边无白色 则其代表地面
//若黑周边有白色 则其代表边缘
//总结出三类元素
//河流 地面 边缘
//每类元素 都有一个图片集合 来提供内容
//河流 从河流 id中随便找个图 目前可能就一张图
//地面 从地面 id中随便找个图 目前先随便
//边缘 这个就会比较麻烦 需要知道具体是什么边缘 根据 8方向的像素 是否为白色 来决定其 到底选择哪个边缘图

func init() {
	if len(os.Args) != 5 {
		fmt.Println("必须输入 框架像素图,土块图集文件,土块id映射的json文件格式, 以及生成的tmx文件名。如: \ngentmx.exe 地图x02.jpg tile图集.png tileID.json tmx名字.tmx")
		os.Exit(1)
	}

	//一定要先创建classinfo到土块id的映射先
	GenClassToIDMap(os.Args[3])
}

func main() {

	//将颜色图 映射为classinfo 二维数组
	cm := newClassmap(os.Args[1], os.Args[2])
	cm.toPNG("classView.png")
	fmt.Println(os.Args[4])
	cm.toTMX(os.Args[4])

}

//上面是命令行版本的实现,但是公司没人会用 还是自己跑吧 下面是源码修改版本

// func main() {
// 	//一定要先创建classinfo到土块id的映射先
// 	GenClassToIDMap("tileID1.json")
// 	//将颜色图 映射为classinfo 二维数组
// 	cm := newClassmap("地图x02.jpg", "橡皮擦.jpg")
// 	cm.toPNG("classView.png")
// 	cm.toTMX("map1.tmx")
// }
