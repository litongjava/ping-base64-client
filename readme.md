# ping-base64-client
使用http进行远程部署的客户端
## 构建
```shell
go build
```
## 使用
需要配合服务端:ping-base64-webapi一起使用
```
ping-base64-client.exe -url http://192.168.3.9:10405/file/upload-run/ -file {localfile} -m {savePath} -d {unzipPath} -c "{cmd}"
```
```shell
-url url:服务器全url
-file localfile:本地文件通常是一个压缩包,上传到服务器上
-m savePath:上传文件保存路径
-d unzipPath:上传文件解压路径
-c cmd:解压完成执行的命令
```

