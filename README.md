# 自动化Docker部署运维平台

分为`auto-server`部署服务及`auto-client`跨平台命令行客户端

## 工作流程

### 自动

自动化阶段根据分支进行

#### 分支规则

*功能性分支： 以`feature/*`命名。*
> 例如新增排序功能，应从master分支check一个`feature/sort`分支。

* 具有版本唯一性，版本号为去掉`feature/`后的名字。运维平台默认不进行灰度部署（但失败时仍然会自动回滚），成功后上个版本镜像将被删除
* 默认固定部署到测试服务上
* 如果是java项目默认开启远程调试，可在auto-client上获取到当前项目的远程调试端口，远程调试端口不固定，可能需要在idea上进行修改，后续考虑在此基础上看是否可以编写idea插件进行自动化远程调试端口。其他项目自行调试，可以考虑在auto.toml中自定义远程调试端口


*test分支： 项目测试分支*
* 此分支为项目集成测试分支,此版本*所有功能性分支完成后*都应当合并到此分支并提交，有自动化平台默认部署到测试服，由测试人员进行集成测试


*master分支： 项目主分支*
* 主分支不会进行镜像构建，但仍然会持续集成
* 所有feature分支均应该从本分支check，并最终完成后合并到测试分支再合并到主分支


*其他分支： 个人开发暂存分支*
* 运维平台不进行任何操作，所有个人代码如在需要的情况下酌情将代码上传到这些分支。

#### 上线规则

为分支打上版本号tag（规则必须为 `主版本.次版本.功能版本`）时，自动上线。

#### 自动化流程

编写项目代码 => 编写.drone.yml(使用工具生成模板并修改,服务端会为drone提供默认配置文件，即可以不写) => 编写Dockerfile(使用工具生成模板并修改) => 编写auto.toml配置特殊启动参数 => 上传到git仓库 => drone平台进行持续集成并上传到docker registry（失败发送邮件） 
=> auto-server捕获到成功通知并根据Dockerfile及auto.toml进行部署（此阶段可在客户端调用auto-client进行查看） => auto-server监测启动成功 => auto-server进行nginx相应自动更改 => 部署流程结束 => auto-server进行后续运行状态监控

### 手动

手动进行版本部署

#### 手动部署/回滚

运行auto-client,如果监测到所在目录为工程项目目录（[gitlab domain]:[group/user]/\*）下。自动进入该项目版本列表，否则进入所有项目列表并选择相应项目 => 选择版本 => 选择部署测试环境或正式环境 => 运行部署 => auto-server停止已运行服务 => 部署完成

## 功能覆盖

* 涵盖从私有镜像到最终部署上线及相应监控的中台应用
* 配合Webhook(drone/gitlab)进行自动化部署
* 根据Dockerfile声明式指令/auto.toml特殊配置控制启动容器
* Dockerfile/DockerCompose/DroneYML模板支持（公司内部）
* DockerCompose项目支持
* 自动端口映射并调整nginx指向，以实现无需人工干预的部分灰度部署（用户无感升级）
* 服务器选择
* 自动化持续部署
* 自动/手动回滚版本
* 可手动发布版本
* 监控容器运行状态

## 开发

golang环境：`>= go v1.12`

格式化：请使用go-imports进行格式化(不要使用go-fmt)
