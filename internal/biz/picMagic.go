package biz

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/qx66/picMagic/internal/conf"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Pic struct {
	Origin string
	Style  map[string]PicStyle
	logger *zap.Logger
}

type PicStyle struct {
	Height int
	Weight int
}

func NewPic(bootstrap *conf.Bootstrap, logger *zap.Logger) *Pic {
	styles := make(map[string]PicStyle)
	
	for _, style := range bootstrap.Magic.Styles {
		styles[style.Name] = PicStyle{
			Height: int(style.Height),
			Weight: int(style.Weight),
		}
	}
	
	logger.Info(
		"启动参数",
		zap.String("origin", bootstrap.Magic.Origin),
		zap.Any("styles", styles),
	)
	
	return &Pic{
		Origin: bootstrap.Magic.Origin,
		Style:  styles,
		logger: logger,
	}
}

var ProviderSet = wire.NewSet(NewPic)

/*
1. 当不使用样式访问的时候 (即: 无样式)，直接返回原图
2. 当样式不匹配的时候 (即: 无样不存在样式表中)，直接返回原图

1. 使用!进行样式判断 -- 更好的命中缓存 (!是uri中的一部分)

filepath中需要杜绝 ! 号的出现

存在以下问题:
	1. 当路径中出现!，但是没有给出样式 (路径为 /img/avatar/q!x.png)
		这个时候就会把:
			x.png 当成样式拦截
			/img/avatar/q 当成路径
			/img/avatar/q 为一个无格式文件
*/

func (pic *Pic) PicMagic(c *gin.Context) {
	filepath := c.Param("filepath")
	var realFilePath string
	var style string
	
	//
	if filepath == "/favicon.ico" {
		return
	}
	if filepath == "/" {
		c.Data(200, "text/html; charset=utf-8", []byte("请输出正确的url资源, 不支持/"))
		return
	}
	
	//
	paramSlice := strings.Split(filepath, "!")
	switch len(paramSlice) {
	case 1:
		// 没有携带样式标签
		realFilePath = filepath
	case 2:
		realFilePath = paramSlice[0]
		style = paramSlice[len(paramSlice)-1]
	default:
		realFilePath = strings.Join(paramSlice[0:len(paramSlice)-1], "!")
		style = paramSlice[len(paramSlice)-1]
	}
	
	f, err := imaging.FormatFromFilename(realFilePath)
	if err != nil {
		c.Data(200, "text/html; charset=utf-8", []byte("暂不支持该资源格式"))
		return
	}
	
	//
	rUrl, err := url.JoinPath(pic.Origin, realFilePath)
	if err != nil {
		c.Data(200, "text/html; charset=utf-8", []byte("url 异常"))
		return
	}
	
	//
	resp, err := http.Get(rUrl)
	if err != nil {
		c.Data(200, "text/html; charset=utf-8", []byte("系统异常"))
		return
	}
	
	respBody := resp.Body
	defer respBody.Close()
	
	respBodyByte, err := io.ReadAll(respBody)
	if resp.StatusCode != http.StatusOK {
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBodyByte)
		return
	}
	
	//
	resPic := &bytes.Buffer{}
	resPic.Write(respBodyByte)
	
	if style == "" {
		c.Data(200, resp.Header.Get("Content-Type"), resPic.Bytes())
		return
	}
	
	picStyle, ok := pic.Style[style]
	if !ok {
		c.Data(200, resp.Header.Get("Content-Type"), resPic.Bytes())
		return
	}
	
	img, err := imaging.Decode(resPic)
	
	thumbnail := imaging.Thumbnail(img, picStyle.Weight, picStyle.Height, imaging.Lanczos)
	destPic := &bytes.Buffer{}
	err = imaging.Encode(destPic, thumbnail, f)
	if err != nil {
		c.Data(200, "text/html; charset=utf-8", []byte("图片编码异常"))
		return
	}
	
	c.Data(200, resp.Header.Get("Content-Type"), destPic.Bytes())
	return
}
