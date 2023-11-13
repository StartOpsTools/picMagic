package main

import (
	"fmt"
	"github.com/ExportersTools/picMagic/internal/biz"
	"github.com/gin-gonic/gin"
)

func main() {
	style := make(map[string]biz.PicStyle)
	style["128"] = biz.PicStyle{
		Weight: 128,
		Height: 128,
	}
	
	style["256"] = biz.PicStyle{
		Weight: 256,
		Height: 256,
	}
	
	style["360"] = biz.PicStyle{
		Weight: 360,
		Height: 360,
	}
	
	style["400"] = biz.PicStyle{
		Weight: 400,
		Height: 400,
	}
	
	style["512"] = biz.PicStyle{
		Weight: 512,
		Height: 512,
	}
	
	style["360x203"] = biz.PicStyle{
		Weight: 360,
		Height: 203,
	}
	
	pic := &biz.Pic{
		Origin: "https://startops-static.oss-cn-hangzhou.aliyuncs.com",
		Style:  style,
	}
	
	gin.SetMode(gin.DebugMode)
	route := gin.New()
	
	route.GET("/*filepath", pic.PicMagic)
	
	err := route.Run(":20000")
	if err != nil {
		errMessage := fmt.Sprint("Start Gateway Server Error,", err)
		panic(errMessage)
	}
}
