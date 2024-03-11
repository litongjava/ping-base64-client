# ping-base64-client

ping-base64-client是一个部署客户端,需要配合服务端:ping-base64-webapi一起使用.
工作流程如下
- 上传文件到服务器的指定目录
- 解压文件到指定目录(可以省略)
- 在服务器上运行一条命令

## 构建

```shell
go build
```
or
```shell
go install
```

## 使用
### 方式1: 通过命令行指定参数

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

eg

```
ping-base64-client -url http:/xxxx:10405/file/upload-run/ -file target\malang-pen-api-server-1.0.0.jar -m /data/apps/webapps/malang_pen_api_server -c "docker restart malang_pen_api_server"
```

## 方式2 : 通过配置文件指定参数

```
ping-base64-client -a upload-run
```

运行过程中会读取配置配置ping-base64.toml,文件示例如下

```
[upload-run]
url = "http://xxxx/file/upload-run/"
file = "target/malang-pen-api-server-1.0.0.jar"
m = "/data/apps/webapps/malang_pen_api_server"
c = "docker restart malang_pen_api_server"
```

注: 命令行的参数优先级高于配置文件

