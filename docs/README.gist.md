---
title: Go SDK 使用指南 | 七牛云存储
---

# Go SDK 使用指南

此 Golang SDK 适用于所有 >=go1 版本，基于 [七牛云存储官方API](http://docs.qiniutek.com/v3/api/) 构建。使用此 SDK 构建您的网络应用程序，能让您以非常便捷地方式将数据安全地存储到七牛云存储上。无论您的网络应用是一个网站程序，还是包括从云端（服务端程序）到终端（手持设备应用）的架构的服务或应用，通过七牛云存储及其 SDK，都能让您应用程序的终端用户高速上传和下载，同时也让您的服务端更加轻盈。

目录
----
- [安装](#install)
- [初始化](#setup)
	- [配置密钥](#setup-key)
- [资源管理接口](#rs-api)
	- [查看单个文件属性信息](#rs-stat)
	- [复制单个文件](#rs-copy)
	- [移动单个文件](#rs-move)
	- [删除单个文件](#rs-delete)
	- [批量操作](#batch)
		- [批量获取文件属性信息](#batch-stat)
		- [批量复制文件](#batch-copy)
		- [批量移动文件](#batch-move)
		- [批量删除文件](#batch-delete)
- [上传下载接口](#get-and-put-api)
	- [上传下载授权](#token)
		- [生成uptoken](#make-uptoken)
		- [生成downtoken](#make-downtoken)
	- [文件上传](#upload)
		- [普通上传](#io-upload)
		- [断点续上传](#resumable-io-upload)
	- [文件下载](#io-download)
		- [公有资源下载](#public-download)
		- [私有资源下载](#private-download)
- [数据处理接口](#fop-api)
	- [图像](#fop-image)
		- [查看图像属性](#fop-image-info)
		- [查看图片EXIF信息](#fop-exif)
		- [生成图片预览](#fop-image-view)
- [贡献代码](#contribution)
- [许可证](#license)

----

<a name=install></a>
## 1. 安装
在命令行下执行

	go get github.com/qiniu/api

<a name=setup></a>
## 2. 初始化
<a name=setup-key></a>
### 2.1 配置密钥

要接入七牛云存储，您需要拥有一对有效的 Access Key 和 Secret Key 用来进行签名认证。可以通过如下步骤获得：

1. [开通七牛开发者帐号](https://dev.qiniutek.com/signup)
2. [登录七牛开发者自助平台，查看 Access Key 和 Secret Key](https://dev.qiniutek.com/account/keys) 。

在获取到 Access Key 和 Secret Key 之后，您可以在您的程序中调用如下两行代码进行初始化对接, 要确保`ACCESS_KEY` 和 `SECRET_KEY` 在调用所有七牛API服务之前均已赋值：

```{go}
@gist(gist/conf.go#import)

func main() {
	@gist(gist/conf.go#init)
}
```

<a name=get-and-put-api></a>
## 3. 上传下载接口

<a name="io-put-flow"></a>
### 3.1 上传流程

在七牛云存储中，整个上传流程大体分为这样几步：

1. 业务服务器颁发 [uptoken（上传授权凭证）](http://docs.qiniu.com/api/put.html#uploadToken)给客户端（终端用户）
2. 客户端凭借 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken) 上传文件到七牛
3. 在七牛获得完整数据后，发起一个 HTTP 请求回调到业务服务器
4. 业务服务器保存相关信息，并返回一些信息给七牛
5. 七牛原封不动地将这些信息转发给客户端（终端用户）

需要注意的是，回调到业务服务器的过程是可选的，它取决于业务服务器颁发的 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken)。如果没有回调，七牛会返回一些标准的信息（比如文件的 hash）给客户端。如果上传发生在业务服务器，以上流程可以自然简化为：

1. 业务服务器生成 uptoken（不设置回调，自己回调到自己这里没有意义）
2. 凭借 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken) 上传文件到七牛
3. 善后工作，比如保存相关的一些信息

### 3.2 生成上传授权uptoken
uptoken是一个字符串，作为http协议Header的一部分（Authorization字段）发送到我们七牛的服务端，表示这个http请求是经过认证的。

```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)

func main() {
	@gist(gist/conf.go#init)
	@gist(gist/rs.go#put_policy)
}
```
参阅 `rs.PutPolicy`

### 3.3 上传代码
直接上传二进制流
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)
@gist(gist/io.go#import)

func main() {
	@gist(gist/conf.go#init)
	
	@gist(gist/io.go#put_policy)
	@gist(gist/io.go#put_extra)

	@gist(gist/io.go#put)
}
```

上传本地文件
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)
@gist(gist/io.go#import)

func main() {
	@gist(gist/conf.go#init)

	@gist(gist/io.go#put_policy)
	@gist(gist/io.go#put_extra)

	@gist(gist/io.go#put_file)
}
```


<a name="io-put-policy"></a>

### 3.3 上传策略

[uptoken](http://docs.qiniu.com/api/put.html#uploadToken) 实际上是用 AccessKey/SecretKey 进行数字签名的上传策略(`rs.PutPolicy`)，它控制则整个上传流程的行为。让我们快速过一遍你都能够决策啥：

* `expires` 指定 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken) 有效期（默认1小时）。一个 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken) 可以被用于多次上传（只要它还没有过期）。
* `scope` 限定客户端的权限。如果 `scope` 是 bucket，则客户端只能新增文件到指定的 bucket，不能修改文件。如果 `scope` 为 bucket:key，则客户端可以修改指定的文件。
* `callbackUrl` 设定业务服务器的回调地址，这样业务服务器才能感知到上传行为的发生。可选。
* `asyncOps` 可指定上传完成后，需要自动执行哪些数据处理。这是因为有些数据处理操作（比如音视频转码）比较慢，如果不进行预转可能第一次访问的时候效果不理想，预转可以很大程度改善这一点。
* `returnBody` 可调整返回给客户端的数据包（默认情况下七牛返回文件内容的 `hash`，也就是下载该文件时的 `etag`）。这只在没有 `callbackUrl` 时有效。
* `escape` 为真（非0）时，表示客户端传入的 `callbackParams` 中含有转义符。通过这个特性，可以很方便地把上传文件的某些元信息如 `fsize`（文件大小）、`imageInfo.width/height`（图片宽度/高度）、`exif`（图片EXIF信息）等传给业务服务器。
* `detectMime` 为真（非0）时，表示服务端忽略客户端传入的 `mimeType`，自己自行检测。

关于上传策略更完整的说明，请参考 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken)。

<a name="resumable-io-put"></a>

### 3.4 断点续上传、分块并行上传

除了基本的上传外，七牛还支持你将文件切成若干块（除最后一块外，每个块固定为4M大小），每个块可独立上传，互不干扰；每个分块块内则能够做到断点上续传。

我们先看支持了断点上续传、分块并行上传的基本样例：

上传二进制流
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)
@gist(gist/resumable_io.go#import)

func main() {
	@gist(gist/conf.go#init)
	@gist(gist/resumable_io.go#put_policy)
	@gist(gist/resumable_io.go#put_extra)
	@gist(gist/resumable_io.go#put)
}
```
参阅: `resumable.io.Put`, `resumable.io.PutExtra`, `rs.PutPolicy`

上传本地文件
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)
@gist(gist/resumable_io.go#import)

func main() {
	@gist(gist/conf.go#init)
	@gist(gist/resumable_io.go#put_policy)
	@gist(gist/resumable_io.go#put_extra)
	@gist(gist/resumable_io.go#put_file)
}
```
参阅: `resumable.io.PutFile`, `resumable.io.PutExtra`, `rs.PutPolicy`

相比普通上传，断点上续传代码没有变复杂。基本上就只是将`io.PutExtra`改为`resumable.io.PutExtra`，`io.PutFile`改为`resumable.io.PutFile`。

但实际上 `resumable.io.PutExtra` 多了不少配置项，其中最重要的是两个回调函数：`Notify` 与 `NotifyErr`，它们用来通知使用者有更多的数据被传输成功，或者有些数据传输失败。在 `Notify` 回调函数中，比较常见的做法是将传输的状态进行持久化，以便于在软件退出后下次再进来还可以继续进行断点续上传。但不传入 `Notify` 回调函数并不表示不能断点续上传，只要程序没有退出，上传失败自动进行续传和重试操作。

<a name=make-downtoken></a>
#### 3.1.2 下载授权downtoken
当想要下载私密bucket的资源时，需要提供download token，在SDK中，将不对外提供直接获取Token的接口，具体可以参照[私有资源下载](#private-download)

<a name=upload></a>
### 3.2 文件上传

为了尽可能地改善终端用户的上传体验，七牛云存储首创了客户端直传功能。一般云存储的上传流程是：

    客户端（终端用户） => 业务服务器 => 云存储服务

这样多了一次上传的流程，和本地存储相比，会相对慢一些。但七牛引入了客户端直传，将整个上传过程调整为：

    客户端（终端用户） => 七牛 => 业务服务器

客户端（终端用户）直接上传到七牛的服务器，通过DNS智能解析，七牛会选择到离终端用户最近的ISP服务商节点，速度会比本地存储快很多。文件上传成功以后，七牛的服务器使用回调功能，只需要将非常少的数据（比如Key）传给应用服务器，应用服务器进行保存即可。

**注意**：如果您只是想要上传已存在您电脑本地或者是服务器上的文件到七牛云存储，可以直接使用七牛提供的 [qrsync](/v3/tools/qrsync/) 上传工具。
文件上传有两种方式，一种是以普通方式直传文件，简称普通上传，另一种方式是断点续上传，断点续上传在网络条件很一般的情况下也能有出色的上传速度，而且对大文件的传输非常友好。

<a name=io-upload></a>
### 3.2.1 普通上传
普通上传的接口在 `github.com/qiniu/api/io` 里，如下：


<a name=resumable-io-upload></a>
### 3.2.2 断点续上传
<a name=io-download></a>
### 3.3 文件下载
七牛云存储上的资源下载分为 公有资源下载 和 私有资源下载 。

私有（private）是 Bucket（空间）的一个属性，一个私有 Bucket 中的资源为私有资源，私有资源不可匿名下载。

新创建的空间（Bucket）缺省为私有，也可以将某个 Bucket 设为公有，公有 Bucket 中的资源为公有资源，公有资源可以匿名下载。

<a name=public-download></a>
#### 3.3.1 公有资源下载
如果在给bucket绑定了域名的话，可以通过以下地址访问。

	[GET] http://<domain>/<key>

其中<domain>可以到[七牛云存储开发者自助网站](https://dev.qiniutek.com/buckets)绑定, 域名可以使用自己一级域名的或者是由七牛提供的二级域名(`<bucket>.qiniutek.com`)。注意，尖括号不是必需，代表替换项。

<a name=private-download></a>
#### 3.3.2 私有资源下载
私有资源必须通过临时下载授权凭证(downloadToken)下载，如下：

	[GET] http://<domain>/<key>?token=<downloadToken>

注意，尖括号不是必需，代表替换项。  
`downloadToken` 可以使用 SDK 提供的如下方法生成：

```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)
@gist(gist/io.go#import)

func main() {
	@gist(gist/conf.go#init)
	@gist(gist/io.go#download)
}
```
参阅: `rs.GetPolicy`, `io.GetUrl`

<a name=rs-api></a>
## 4. 资源管理接口

文件管理包括对存储在七牛云存储上的文件进行查看、复制、移动和删除处理。
该节调用的函数第一个参数都为 `logger`, 用于记录log, 如果无需求, 可以设置为nil. 具体接口可以查阅 `github.com/qiniu/rpc`

<a name=rs-stat></a>
### 4.1 查看单个文件属性信息
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)

func main() {
	@gist(gist/conf.go#init)
	@gist(gist/rs.go#stat)
}
```
参阅: `rs.Entry`, `rs.Client.Stat`


<a name=rs-copy></a>
### 4.2 复制单个文件
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)

func main() {
	@gist(gist/conf.go#init)

	@gist(gist/rs.go#copy)
}
```
参阅: `rs.Client.Copy`

<a name=rs-move></a>
### 4.3 移动单个文件
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)

func main() {
	@gist(gist/conf.go#init)

	@gist(gist/rs.go#move)
}
```
参阅: `rs.Client.Move`

<a name=rs-delete></a>
### 4.4 删除单个文件
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)

func main() {
	@gist(gist/conf.go#init)

	@gist(gist/rs.go#delete)
}
```
参阅: `rs.Client.Delete`

<a name=batch></a>
### 4.5 批量操作
当您需要一次性进行多个操作时, 可以使用批量操作.
<a name=batch-stat></a>
#### 4.5.1 批量获取文件属性信息
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)

func main() {
	@gist(gist/conf.go#init)

	@gist(gist/rs.go#entry_pathes)
	@gist(gist/rs.go#batch_stat)
}
```

参阅: `rs.EntryPath`, `rs.BatchStatItemRet`, `rs.Client.BatchStat`

<a name=batch-copy></a>
#### 4.5.2 批量复制文件
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)

func main() {
	@gist(gist/conf.go#init)

	@gist(gist/rs.go#entry_path_pairs)
	@gist(gist/rs.go#batch_copy)
}
```

参阅: `rs.BatchItemRet`, `rs.EntryPathPair`, `rs.Client.BatchCopy`

<a name=batch-move></a>
#### 4.5.3 批量移动文件
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)

func main() {
	@gist(gist/conf.go#init)

	@gist(gist/rs.go#entry_path_pairs)
	@gist(gist/rs.go#batch_move)
}
```
参阅: `rs.EntryPathPair`, `rs.Client.BatchMove`

<a name=batch-delete></a>
#### 4.5.4 批量删除文件
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)

func main() {
	@gist(gist/conf.go#init)

	@gist(gist/rs.go#entry_pathes)
	@gist(gist/rs.go#batch_delete)
}
```
参阅: `rs.EntryPath`, `rs.Client.BatchDelete`

<a name=batch-advanced></a>
#### 4.5.5 高级批量操作
批量操作不仅仅支持同时进行多个相同类型的操作, 同时也支持不同的操作.
```{go}
@gist(gist/conf.go#import)
@gist(gist/rs.go#import)

func main() {
	@gist(gist/conf.go#init)

	@gist(gist/rs.go#batch_adv)
}
```
参阅: `rs.URIStat`, `rs.URICopy`, `rs.URIMove`, `rs.URIDelete`, `rs.Client.Batch`

<a name=fop-api></a>
## 5. 数据处理接口
七牛支持在云端对图像, 视频, 音频等富媒体进行个性化处理

<a name=fop-image></a>
### 5.1 图像
<a name=fop-image-info></a>
### 5.1.1 查看图像属性
```{go}
@gist(gist/conf.go#import)
@gist(gist/fop.go#import)

func main() {
	@gist(gist/fop.go#imageurl)
	@gist(gist/fop.go#image_info)
}
```
参阅: `fop.ImageInfoRet`, `fop.ImageInfo`

<a name=fop-exif></a>
### 5.1.2 查看图片EXIF信息
```{go}
@gist(gist/conf.go#import)
@gist(gist/fop.go#import)

func main() {
	@gist(gist/fop.go#imageurl)
	@gist(gist/fop.go#exif)
}
```
参阅: `fop.Exif`, `fop.ExifRet`, `fop.ExifValType`

<a name=fop-image-view></a>
### 5.1.3 生成图片预览
```{go}
@gist(gist/conf.go#import)
@gist(gist/fop.go#import)

func main() {
	@gist(gist/fop.go#imageurl)
	@gist(gist/fop.go#image_view)
}
```
参阅: `fop.ImageView`

<a name=contribution></a>
## 6. 贡献代码

1. Fork
2. 创建您的特性分支 (`git checkout -b my-new-feature`)
3. 提交您的改动 (`git commit -am 'Added some feature'`)
4. 将您的修改记录提交到远程 `git` 仓库 (`git push origin my-new-feature`)
5. 然后到 github 网站的该 `git` 远程仓库的 `my-new-feature` 分支下发起 Pull Request

<a name=license></a>
## 7. 许可证

Copyright (c) 2013 qiniu.com

基于 MIT 协议发布:

* [www.opensource.org/licenses/MIT](http://www.opensource.org/licenses/MIT)
