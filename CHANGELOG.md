## CHANGE LOG

### v0.9.1

2013-05-28 issue [#56](https://github.com/qiniu/api/pull/56)

- 修复 go get github.com/qiniu/api 失败的错误
- 遵循 [sdkspec v1.0.1](https://github.com/qiniu/sdkspec/tree/v1.0.1)
  - io.GetUrl 改为 rs.MakeBaseUrl 和 rs.GetPolicy.MakeRequest
  - rs.PutPolicy 支持 ReturnUrl, ReturnBody, CallbackBody；将 Customer 改为 EndUser
- 增加 github.com/api/url: Escape/Unescape


### v0.9.0

2013-04-08 issue [#33](https://github.com/qiniu/api/pull/33)

- 更新API文档
- 增加断点续上传resumable/io功能
- 移除bucket相关的增加/删除/列出所有bucket等管理操作，推荐到七牛云存储开发者后台中使用这些功能。
