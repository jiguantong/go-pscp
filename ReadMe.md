## 一个简单的传输工具, 所有传输到服务器的任务都可通过该工具简化为一条命令(pscp)
将部署时xftp及xshell等的操作全部在开发工具控制台完成, 效率提升至少20倍
-------------------------
<br>例:<br>

前端打包部署流程
- 之前:
  ### yarn build => 打开xftp => 找到部署目录 => 打开打包好的文件夹所在目录 => 拖动上传 => 回到webstorm
- 现在:
  ### yarn build && pscp.exe
-------------------------
tomcat项目部署流程
- 之前:
  ### build => 打开xftp => 找到部署目录 => 打开打包好的文件夹所在目录 => 拖动上传 => 切换到xshell => shutdown tomcat => start tomcat =>tailf log => 回IDEA 
- 现在:
  ### pscp.exe

## 使用说明
将pscp.exe和pscp.yml放到项目根目录, 配置pscp.yml, 即可在jetbrains工具的控制台使用pscp命令

## 配置说明
 ```
################################################################################################
# localdir 后不带 "/" 表示将该目录直接放到 { remotedir } 中
# 例 1: remotedir="/opt/project" localdir="../go-pscp" 该目录下有 1.txt文件, 部署后1.txt文件路径为
# /opt/project/go-pscp/1.txt
#-----------------------------------------------------------------------------------------------
# localdir 后带 "/" 表示将该目录下所有文件部署到 { remotedir }中
# 例 2: remotedir="/opt/project" localdir="../go-pscp/" 该目录下有 1.txt文件, 部署后1.txt文件路径为
# /opt/project/1.txt
################################################################################################

# 服务器IP地址
ip: 47.102.196.137
# 服务器用户
user: root
# 服务器密码
password: 123456
# ssh端口
port: 22
# 服务器部署目录
remotedir: "/opt/project"
# 本地待部署项目目录
localdir: "../pscp/"
# 传输完成后执行该命令
# 示例: 重启tomcat并跟踪日志  !!! 如果仅需传输文件, 注释该行即可
cmd: "source /etc/profile&&/opt/tomcat/bin/shutdown.sh&&/opt/tomcat/bin/startup.sh&&tail -f /opt/tomcat/logs/catalina.out"

```
### build
`go build -ldflags="-w -s"`
`upx.exe -9 -k "jscp.exe"`
