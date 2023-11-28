package main

import (
	"bytes"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/qx66/picMagic/internal/biz"
	"github.com/qx66/picMagic/internal/conf"
	"github.com/qx66/picMagic/pkg/middleware"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "configPath", "", "-configPath")
}

type app struct {
	pic *biz.Pic
}

func newApp(pic *biz.Pic) *app {
	return &app{
		pic: pic,
	}
}

func main() {
	flag.Parse()
	
	logger, _ := zap.NewProduction(zap.Fields(zap.String("service", "picMagic")))
	defer logger.Sync()
	
	if configPath == "" {
		logger.Error("configPath 参数为空")
		return
	}
	
	//
	f, err := os.Open(configPath)
	defer f.Close()
	if err != nil {
		logger.Error(
			"加载配置文件失败",
			zap.String("configPath", configPath),
			zap.Error(err),
		)
		return
	}
	
	//
	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		logger.Error(
			"加载配置文件copy内容失败",
			zap.Error(err),
		)
		return
	}
	
	//
	var bootstrap conf.Bootstrap
	err = yaml.Unmarshal(buf.Bytes(), &bootstrap)
	if err != nil {
		logger.Error(
			"序列化配置失败",
			zap.Error(err),
		)
		return
	}
	
	//
	app, err := initApp(&bootstrap, logger)
	if err != nil {
		logger.Error(
			"初始化 picMagic 服务失败",
			zap.Error(err),
		)
		panic(err)
	}
	
	gin.SetMode(gin.DebugMode)
	route := gin.New()
	route.Use(middleware.Recording(logger))
	route.GET("/*filepath", app.pic.PicMagic)
	
	err = route.Run(":20000")
	if err != nil {
		logger.Error(
			"启动 picMagic 服务失败",
			zap.Error(err),
		)
		panic(err)
	}
}
