# beego-api-1

 一个简单的beego API demo，用作前后端分离的后端脚手架
 包含基本的日志、redis、数据库基本功能，用户、资源、角色model
 
使用方法：
go get github.com/tyrantkemp/beego-api-1 
或者直接 git clone 到$GOPATH/src/github.com/下

项目根目录下 glide install 
执行根目录下 api.sql 本地数据库建立对应表（对应的数据库信息见conf/app.conf）

第一次运行：bee run -gendoc=true -downdoc=true  初始化api文档并下载swagger文件

如果数据库新增表，想要项目自动生成对应的model,controller,router,可在根目录下直接<br>
bee generate appcode -tables="tablename" -conn="root:root@tcp(127.0.0.1:3306)/beegoApi" -level=1<br>
其中（-level:  [1 | 2 | 3], 1 = models; 2 = models,controllers; 3 = models,controllers,router）<br>
如果要更新全部model <br>
bee generate appcode -conn="root:root@tcp(127.0.0.1:3306)/beegoApi" -level=3<br>

