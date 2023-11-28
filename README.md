# picMagic

picMagic 是一个简单的图片处理的web工具

## feature

1. 提供原图大小剪切功能

## 约定

1. 当不使用样式访问的时候 (即: 无样式)，直接返回原图
2. 当样式不匹配的时候 (即: 无样不存在样式表中)，直接返回原图

为什么使用!进行样式判断:
    使用!进行样式判断 -- 更好的命中缓存 (!是uri中的一部分)

filepath中需要杜绝 ! 号的出现

存在以下问题:
1.  当路径中出现!，但是没有给出样式 (路径为 /img/avatar/q!x.png)
    这个时候就会把:
        x.png 当成样式拦截
        /img/avatar/q 当成路径
        /img/avatar/q 为一个无格式文件

## 背景

在使用 AWS 服务时，如果需要对S3内的图片进行缩略图等处理，需要使用 AWS 的 Lambda 功能，
自己编写 Lambda 函数，然后代码传到 Lambda 服务中，然后一步一步进行调试，中间还涉及到 API Gateway 等云功能。

个人使用起来感觉不是很方便，只是图片处理，每次调试都要费很长时间，所以编写了一个 picMagic 服务，
请求到 picMagic 服务，回源，根据 style 匹配，处理图片，返回。

在 picMagic 服务中间最好加一层 cdn 缓存，这样可以使 picMagic 服务的 QPS 降低，节省处理器消耗

## Usage

```shell
./bin/picMagic-linux  -configPath=./configs/config.yaml
```

## 测试

以自带配置文件案例测试,浏览器访问:

原图:

    http://127.0.0.1:20000/startops/icon/startops1.png
256样式:

    http://127.0.0.1:20000/startops/icon/startops1.png!256

不存在256x样式:

    http://127.0.0.1:20000/startops/icon/startops1.png!256x

## 图片处理

图片处理一般是对图片进行以下动作处理 (目前 picMagic 只支持更改长宽处理)

图片效果:

    亮度
    对比度
    锐化
    模糊

缩略效果:

    高度
    宽度

水印设置:

    文字水印
    图片水印
