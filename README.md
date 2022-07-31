# 🎉LPan

![pass](https://img.shields.io/badge/building-pass-green) ![pass](https://img.shields.io/badge/checks-pass-green)

## 🎁特点

- [x] 一个**轻量级**网盘
- [x] **限速**
- [x] **断点续传**
- [x] **文件秒传**
- [x] 轻松**权限管理** 一键生成**加密二维码**
- [x] **回收站**系统，**自动清理过期链接、删除文件**

## 🚀架构

设计思路是每个哈希值不同的文件挂载在服务器，数据库存储所有文件的信息和用户与文件ID的映射

使用**虚拟目录**以减少系统开销

![image-20220731124010610](http://typora.fengxiangrui.top/1659242411.png)





## 🎨一些接口与实现

- **上传文件**

```go
fileGroup.POST("/upload", uploadfile)
```

上传文件 会从body的form-data中读取**文件与父级虚拟目录**

然后检查文件**哈希值** 与数据库中已有哈希对比 如果存在相同 直接**秒传**(直接在私人存储表添加一条记录 指向公共存储表中指定文件)

- **断点续传**

```go
fileGroup.POST("/uploadbysilce", uploadfilebysile)
```

基本流程如上 但是将from-data的文件copy到服务器本地路径时，会先生成一个**临时文件，记录传输进度**，下一次再上传时根据临时文件信息 使用seek即可实现断点续传

- **下载文件**

```go
fileGroup.GET("/download/:file_id", downloadfile)
```

先检查url确定该文件是否为分享，然后确**认点击者权限**，如果有权限可进行下载

如果点击者非VIP 速度会被**限制到500KB**

- **边下边播**

```go
fileGroup.GET("/downloadbysilce/:file_id", downloadbyslice)
```

基本流程与下载一致 但传输时会把链接升级为`Websocket` 然后**分片传输**

- **删除文件**

```go
fileGroup.DELETE("/delete/:file_id", deletefile)
```

删除文件是只针对数据库的操作，将用户私人表中的delete字段设置为删除日期即可，**三十天内可找回，三十天后字段自动删除**

- **恢复文件**

```go
fileGroup.GET("/recover/:file_id", recoverfile)
```

恢复文件同样是针对**数据库**的操作，将delete字段设置为空即可恢复（被删除的文件无法下载分享等))

- **文件改名**

```go
fileGroup.POST("/rename/:file_id", renamefile)
```

文件改名只改变**用户私有库中文件名**

- **文件路径修改**

```go
fileGroup.POST("/modify-path/:file_id", modifypath)
```

文件路径修改会修改用户私有库中`father_path`字段，只改变虚拟路径，数据返回前段时即可根据父路径 显示某个路径下所有文件与文件夹

-  **文件分享**

```go
fileGroup.POST("/share/:file_id", sharefile)
```

文件分享会生成一个二维码，它携带参数`share=true&expr_time`分享字段和过期时间（天），在期限前点击链接即可一键入库，然后下载

- **加密分享**

```go
fileGroup.POST("/share/secret/:file_id", sharefilewithsecret)
```

会生成一个加密链接二维码

- **解密重定向链接**

```go
r.GET("/secret/:val", JWTAuthMiddleware(), RateLimitMiddleware(time.Millisecond*100, 2048), func(c *gin.Context) 
```

加密的链接会有一个前缀，所有加密的链接都由这个接口处理，接口拿到密文后解密重定向