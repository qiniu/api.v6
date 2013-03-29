Qiniu Resource (Cloud) Storage SDK for Golang
===

[![Build Status](https://api.travis-ci.org/qiniu/api.png?branch=master)](https://travis-ci.org/qiniu/api)

[![Qiniu Logo](http://qiniutek.com/images/logo-2.png)](http://qiniu.com/)

# Go SDK 使用指南

此 Golang SDK 适用于所有 >=go1 版本，基于 [七牛云存储官方API](http://docs.qiniutek.com/v3/api/) 构建。使用此 SDK 构建您的网络应用程序，能让您以非常便捷地方式将数据安全地存储到七牛云存储上。无论您的网络应用是一个网站程序，还是包括从云端（服务端程序）到终端（手持设备应用）的架构的服务或应用，通过七牛云存储及其 SDK，都能让您应用程序的终端用户高速上传和下载，同时也让您的服务端更加轻盈。

**目录**

- [1. 安装](#install)
- [2. 接入](#setup)
	- [2.1 配置密钥](#setup-key)
- [3. 使用](#usage)
	- [3.1 文件上传](#upload)
		- [3.1.1 生成上传授权凭证](#make-upload-token)
		- [3.1.2 使用凭证上传文件](#upload-put)
	- [3.2 文件下载](#download)
		- [3.2.1 公有资源下载](#public-download)
		- [3.2.2 私有资源下载](#private-download)
	- [3.3 文件管理](#rs)
		- [3.3.1 查看单个文件属性信息](#rs-stat)
		- [3.3.2 复制单个文件](#rs-copy)
		- [3.3.3 移动单个文件](#rs-move)
		- [3.3.4 删除单个文件](#rs-delete)
		- [3.3.5 批量操作](#batch)
			- [3.3.5.1 批量获取文件属性信息](#batch-stat)
			- [3.3.5.2 批量复制文件](#batch-copy)
			- [3.3.5.3 批量移动文件](#batch-move)
			- [3.3.5.4 批量删除文件](#batch-delete)
			- [3.3.5.5 高级批量操作](#batch-advanced)
	- [3.4 云处理](#fop)
		- [3.4.1 图像](#fop-image)
			- [3.4.1.1 查看图像属性](#fop-image-info)
			- [3.4.1.2 查看图片EXIF信息](#fop-exif)
			- [3.4.1.3 生成图片预览](#fop-image-view)
- [4. 贡献代码](#contribution)
- [5. 许可证](#license)

<a name=install></a>
## 1. 安装
在命令行下执行

	go get github.com/qiniu/bytes
	go get github.com/qiniu/rpc
	go get github.com/qiniu/api

<a name=setup></a>
## 2. 接入
<a name=setup-key></a>
### 2.1 配置密钥

要接入七牛云存储，您需要拥有一对有效的 Access Key 和 Secret Key 用来进行签名认证。可以通过如下步骤获得：

1. [开通七牛开发者帐号](https://dev.qiniutek.com/signup)
2. [登录七牛开发者自助平台，查看 Access Key 和 Secret Key](https://dev.qiniutek.com/account/keys) 。

在获取到 Access Key 和 Secret Key 之后，您可以在您的程序中调用如下两行代码进行初始化对接, 要确保`ACCESS_KEY` 和 `SECRET_KEY` 在调用所有七牛API服务之前均已赋值：

```{go}
import . "github.com/qiniu/api/conf"

func main() {
	ACCESS_KEY = "YOUR_APP_ACCESS_KEY"
	SECRET_KEY = "YOUR_APP_SECRET_KEY"
}
```
<a name=usage></a>
## 3. 使用
<a name=upload></a>
### 3.1 文件上传
**注意**：如果您只是想要上传已存在您电脑本地或者是服务器上的文件到七牛云存储，可以直接使用七牛提供的 [qrsync](/v3/tools/qrsync/) 上传工具。如果是需要通过您的网站或是移动应用(App)上传文件，则可以接入使用此 SDK，详情参考如下文档说明。一般包含两个步骤: 生成上传凭证 -> 上传文件
<a name=make-upload-token></a>
#### 3.1.1 生成上传授权凭证
上传文件需要提供上传授权凭证来验证身份, 通过实现 `github.com/api/rs/PutPolicy` 来生成Token, 具体代码如下

```{go}
import "github.com/api/rs"

func main() {
	// 填充ACCESS_KEY和SECRET_KEY, 参考 配置密钥
	
	policy := rs.PutPolicy {
	
		// [string] 必须, 指定授权上传的bucket
		Scope: bucketName, 
		
		// [int] 表示有效时间为3600秒, 即一个小时
		Expires: 3600,
		
		// [string] 用于设置文件上传成功后，七牛云存储服务端要回调客户方的业务服务器地址
		CallbackUrl: "http://example.com",
		
		// [string] 用于设置文件上传成功后，七牛云存储服务端向客户方的业务服务器发送回调请求的 `Content-Type`
		CallbackBodyType: "application/x-www-form-urlencoded",
		
		// [string] 客户方终端用户（End User）的ID，该字段可以用来标示一个文件的属主，这在一些特殊场景下（比如给终端用户上传的图片打上名字水印）非常有用。
		Customer: "",
		
		// [string] 用于设置文件上传成功后，执行指定的预转指令。
		// 参考 http://docs.qiniutek.com/v3/api/io/#uploadToken-asyncOps
		AsyncOps: "",
		
		// [uint16] 可选值 0 或者 1，缺省为 0 。值为 1 表示 callback 传递的自定义数据中允许存在转义符号 `$(VarExpression)
		// 参考 http://docs.qiniutek.com/v3/api/words/#VarExpression
		Escape: 0,
		
		// [uint16] 可选值 0 或者 1, 缺省为 0. 值为 1 表示在服务端自动识别上传文件类型.
		DetectMime: 0,
	}
	
	// 生成 uploadToken, string类型
	token := policy.Token()
}
```

> **特殊参数说明:**  
> 当 `escape` 的值为 `1` 时，常见的转义语法如下：

> - 若 `callbackBodyType` 为 `application/json` 时，一个典型的自定义回调数据（[CallbackParams](http://docs.qiniutek.com/v3/api/io/#CallbackParams)）为：

>     `{foo: "bar", w: $(imageInfo.width), h: $(imageInfo.height), exif: $(exif)}`

> - 若 `callbackBodyType` 为 `application/x-www-form-urlencoded` 时，一个典型的自定义回调数据（[CallbackParams](http://docs.qiniutek.com/v3/api/io/#CallbackParams)）为：

>     `foo=bar&w=$(imageInfo.width)&h=$(imageInfo.height)&exif=$(exif)`


<a name=upload-put></a>
#### 3.1.2 使用凭证上传文件
通过 `io` 可以上传数据到七牛的服务器
然后可以通过 `io.Put()` 上传数据了, 以上传一个文件为例:

```{go}
import (
	"os"
	"github.com/qiniu/api/io"
)

func main() {
	// 1. 配置 ACCESS_KEY 和 SECRET_KEY
	// 2. 生成 uploadToken
	
	// 获得文件的reader, 赋值给f
	filePath := "./test.txt"
	f, _ := os.Open(filePath)
	defer f.Close()
	
	// 声明 PutExtra
	extra := &io.PutExtra {
		// [string] 必选, 指定上传的目标仓库
		Bucket: "bucketName",
		
		// [string] 可选。在 uptoken 没有指定 DetectMime 时，用户客户端可自己指定 MimeType
		MimeType: "",
		
		// [string] 可选。用户自定义 Meta，不能超过 256 字节
		CustomMeta: "",
		
		// [string] 当 uptoken 指定了 CallbackUrl，则 CallbackParams 必须非空
		CallbackParams: "",
	}
	key := "fileKey"
	
	ret, err := io.Put(logger, uploadToken, key, f, extra) // PutRet, error
	if err != nil {
		// 上传失败
		return
	}
	// 当上传成功后, 得到的hash值
	hash := ret.Hash
}
```

<a name=download></a>
### 3.2 文件下载
七牛云存储上的资源下载分为 公有资源下载 和 私有资源下载 。

私有（private）是 Bucket（空间）的一个属性，一个私有 Bucket 中的资源为私有资源，私有资源不可匿名下载。

新创建的空间（Bucket）缺省为私有，也可以将某个 Bucket 设为公有，公有 Bucket 中的资源为公有资源，公有资源可以匿名下载。

<a name=public-download></a>
#### 3.2.1 公有资源下载
如果在给bucket绑定了域名的话, 可以通过以下地址访问. 可以到[七牛云存储开发者自助网站](https://dev.qiniutek.com/buckets)绑定域名, 域名可以使用自己一级域名的或者是由七牛提供的二级域名(`<bucket>.qiniutek.com`).

	[GET] http://<domain>/<key>

注意，尖括号不是必需，代表替换项。

<a name=private-download></a>
#### 3.2.2 私有资源下载
私有资源只能通过临时下载授权凭证(downloadToken)下载，并作为URL的参数 `token` 的值存在于URL.  

	[GET] http://<domain>/<key>?token=<downloadToken>

注意，尖括号不是必需，代表替换项。  
`downloadToken` 可以使用 SDK 提供的如下方法生成：

```{go}
import (
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/io"
)

func main() {
	// 配置 ACCESS_KEY 和 SECRET_KEY
	policy := rs.GetPolicy {
		// [string] 用于设置可匹配的下载链接
		// 参考: http://docs.qiniutek.com/v3/api/io/#download-token-pattern
		Scope: "",
		
		// 用于设置上传 URL 的有效期, 单位为秒
		Expires: 3600,
	}
	
	// downloadToken
	downloadToken := policy.Token()
	
	// 生成下载连接, sourceUrl 为资源原有下载链接
	downloadUrl := io.GetUrl(sourceUrl, downloadToken)
}
```

<a name=rs></a>
### 3.3 文件管理
文件管理包括对存储在七牛云存储上的文件进行查看、复制、移动和删除处理。  
该节调用的函数第一个参数都为 `logger`, 用于记录log, 如果无需求, 可以设置为nil. 具体接口可以查阅 `github.com/qiniu/rpc`

<a name=rs-stat></a>
#### 3.3.1 查看单个文件属性信息
```{go}
import "github.com/qiniu/api/rs"

func main() {
	// 配置 ACCESS_KEY 和 SECRET_KEY
	rsClient := rs.New()
	entry, err := rsClient.Stat(logger, bucket, key) // rs.Entry, error
}
```
附 `rs.Entry` 结构:
```{go}
type Entry struct {
	Hash     string // 文件的特征值，可以看做是基版本号
	Fsize    int64 // 表示文件总大小，单位是 Byte
	PutTime  int64 // 上传时间，单位是 百纳秒
	MimeType string // 文件的 mime-type
	Customer string
}
```

<a name=rs-copy></a>
#### 3.3.2 复制单个文件
```{go}
import "github.com/qiniu/api/rs"

func main() {
	// 配置 ACCESS_KEY 和 SECRET_KEY
	rsClient := rs.New()
	// 返回值 error, 操作成功时err为nil
	err := rsClient.Copy(logger, bucketSrc, keySrc, bucketDest, keyDest) 
	
}
```

<a name=rs-move></a>
#### 3.3.3 移动单个文件
```{go}
import "github.com/qiniu/api/rs"

func main() {
	// 配置 ACCESS_KEY 和 SECRET_KEY
	rsClient := rs.New()
	// 返回值 error, 操作成功时err为nil
	err := rsClient.Move(logger, bucketSrc, keySrc, bucketDest, KeyDest) 
}
```

<a name=rs-delete></a>
#### 3.3.4 删除单个文件
```{go}
import "github.com/qiniu/api/rs"

func main() {
	// 配置 ACCESS_KEY 和 SECRET_KEY
	rsClient := rs.New()
	err := rsClient.Delete(logger, bucket, key) // 返回值 error, 操作成功时err为nil
}
```

<a name=batch></a>
#### 3.3.5 批量操作
当您需要一次性进行多个操作时, 可以使用批量操作.
<a name=batch-stat></a>
##### 3.3.5.1 批量获取文件属性信息
```{go}
import "github.com/qiniu/api/rs"

func main() {
	// 配置 ACCESS_KEY 和 SECRET_KEY
	rsClient = rs.New()
	entryPath1 := rs.EntryPath {
		Bucket: bucket1,
		Key: key1,
	}
	entryPath2 := rs.EntryPath {
		Bucket: bucket2,
		Key: key2,
	}
	
	rets, err := rsClient.BatchStat(logger, entryPath1, entryPath2) 
	// []rs.BatchStatResult, error
}
```
附 `rs.BatchStatResult` 结构
```{go}
type BatchStatResult struct {
	Data Entry // 文件的属性
	Error string // 错误信息, 成功时为""
	Code int // 结果状态代码, 2开头的为成功
}
```

<a name=batch-copy></a>
##### 3.3.5.2 批量复制文件
```{go}
import "github.com/qiniu/api/rs"

func main() {
	// 配置 ACCESS_KEY 和 SECRET_KEY
	rsClient = rs.New()
	// 每个复制操作都含有源文件和目标文件
	entryPair1 := rs.EntryPathPair {
		Src: rs.EntryPath {
			Bucket: bucket1,
			Key: key1,
		},
		Dest: rs.EntryPath {
			Bucket: bucket2,
			Key: key2,
		},
	}
	
	entryPair2 := rs.EntryPathPair {
		Src: rs.EntryPath {
			Bucket: bucket3,
			Key: key3,
		},
		Dest: rs.EntryPath {
			Bucket: bucket4,
			Key: key4,
		},
	}
	
	rets, err := rsClient.BatchCopy(logger, entryPair1, entryPair2) 
	// []rs.BatchResult, error
}
```
附 `rs.BatchItemRet` 结构
```{go}
type BatchItemRet struct {
	Error string // 错误信息, 成功时为""
	Code int // 结果状态代码, 2开头的为成功
}
```

<a name=batch-move></a>
##### 3.3.5.3 批量移动文件
```{go}
import "github.com/qiniu/api/rs"

func main() {
	// 配置 ACCESS_KEY 和 SECRET_KEY
	rsClient = rs.New()
	entryPair1 := rs.EntryPathPair {
		Src: rs.EntryPath {
			Bucket: bucket1,
			Key: key1,
		},
		Dest: rs.EntryPath {
			Bucket: bucket2,
			Key: key2,
		},
	}
	
	entryPair2 := rs.EntryPathPair {
		Src: rs.EntryPath {
			Bucket: bucket3,
			Key: key3,
		},
		Dest: rs.EntryPath {
			Bucket: bucket4,
			Key: key4,
		},
	}
	
	rets, err := rsClient.BatchMove(logger, entryPair1, entryPair2) 
	// []rs.BatchResult, error
}
```

<a name=batch-delete></a>
##### 3.3.5.4 批量删除文件
```{go}
import "github.com/qiniu/api/rs"

func main() {
	// 配置 ACCESS_KEY 和 SECRET_KEY
	rsClient = rs.New()
	entryPath1 := rs.EntryPath {
		Bucket: bucket1,
		Key: key1,
	}
	entryPath2 := rs.EntryPath {
		Bucket: bucket2,
		Key: key2,
	}
	
	rets, err := rsClient.BatchDelete(logger, entryPath1, entryPath2) 
	// []rs.BatchResult, error
}
```

<a name=batch-advanced></a>
##### 3.3.5.5 高级批量操作
批量操作不仅仅支持同时进行多个相同类型的操作, 同时也支持不同的操作.
```{go}
import "github.com/qiniu/api/rs"

func main() {
	ops := []string {
		rs.URICopy(bucket, key1, bucket, key2), // 复制key1到key2
		rs.URIDelete(bucket, key1), // 删除key1
		rs.URIMove(entryURI2, entryURI1), //将key2移动到key1
	}
	
	rets := new([]rs.BatchItemRet)
	rsClient := rs.New()
	err := rsClient.Batch(logger, rets, ops) // 执行操作
	// 结果: rets, 错误: err
}
```

<a name=fop></a>
### 3.4 云处理
七牛支持在云端对图像, 视频, 音频进行个性化处理
<a name=fop-image></a>
#### 3.4.1 图像
<a name=fop-image-info></a>
#### 3.4.1.1 查看图像属性
```{go}
import "github.com/qiniu/api/fop"

func main() {
	imageUrl := "http://domain/key"
	ii := fop.ImageInfo{}
	info, err := ii.MakeRequest(logger, imageUrl) // fop.ImageInfoRet, error
}
```
附 `fop.ImageInfoRet` 结构
```{go}
type ImageInfoRet struct {
	Width int // 原始图片宽度，单位像素
	Height int // 原始图片高度，单位像素
	Format string //原始图片类型
	ColorModel string // 原始图片着色模式
}
```

<a name=fop-exif></a>
#### 3.4.1.2 查看图片EXIF信息
```{go}
import "github.com/qiniu/api/fop"

func main() {
	imageUrl := "http://domain/key"
	exif := fop.Exif{}
	info, err := exif.Call(logger, imageUrl) // fop.ExifRet, error
}
```
附 `fop.Exif` 结构
```{go}
type ExifValType struct {
	Val string
	Type int
}
type ExifRet map[string] ExifValType
```

<a name=fop-image-view></a>
#### 3.4.1.3 生成图片预览
```{go}
import "github.com/qiniu/api/fop"

func main() {
	imageUrl := "http://domain/key"
	iv := fop.ImageView{
		Mode: 1,       // 缩略模式
		Width: 200,    // Width = 0 表示不限定宽度
		Height: 200,   // Height = 0 表示不限定高度
		Quality: 80,   // 质量, 1-100
		Format: "jpg", // 输出文件格式
	}
	previewUrl := iv.MakeRequest(logger, imageUrl)
}
```

<a name=contribution></a>
## 4. 贡献代码

1. Fork
2. 创建您的特性分支 (`git checkout -b my-new-feature`)
3. 提交您的改动 (`git commit -am 'Added some feature'`)
4. 将您的修改记录提交到远程 `git` 仓库 (`git push origin my-new-feature`)
5. 然后到 github 网站的该 `git` 远程仓库的 `my-new-feature` 分支下发起 Pull Request

<a name=license></a>
## 5. 许可证

Copyright (c) 2013 qiniu.com

基于 MIT 协议发布:

* [www.opensource.org/licenses/MIT](http://www.opensource.org/licenses/MIT)
