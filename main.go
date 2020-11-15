package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	//以项目为分组，内部接口全面压缩整合
	//先来个考勤项目
	attendance := router.Group("/attendance")
	{
		//依据所带参数是latest/last-week/last-month，转向不同的处理函数或干脆合并，分情形输出，减少接口
		attendance.POST("/rec",rec) 
		//与上面采取类似的思路
		attendance.POST("/stat",stat)
		attendance.POST("/add", add)
	}	
	/*楚留香项目
	crh:=router.Group("/crh")
	{
		crh.POST("/list", podList)
		crh.POST("/detail", podDetail)
		crh.POST("/container_list", podContainerList)
		crh.POST("/container_detail", podContainerDetail)
	}
	//仙剑四项目
	pal4:=router.Group("/pal4")
	{
		pal4.POST("/list", svcList)
		pal4.POST("/detail", svcDetail)
		pal4.POST("/pods", svcPods)
	}
	//股票项目
	stock:=router.Group("/stock")
	{
		stock.POST("/list", nsList)
		stock.POST("/detail", nsDetail)
		stock.POST("/deploys", nsDeploys)
		stock.POST("/svcs", nsSvcs)
		stock.POST("/pods", nsPods)
	}*/
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
