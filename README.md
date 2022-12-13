# fotff

#### 介绍

fotff(find out the first fault)是为OpenHarmony持续集成设计的问题自动化问题分析工具。

为了平衡开销与收益，考虑到开发效率、资源占用等因素影响，OpenHarmony代码合入门禁（冒烟测试）只拦截部分严重基础问题（例如开机失败、关键进程崩溃、UX布局严重错乱、电话/相机基础功能不可用等）。因此，一些会影响到更细节功能、影响兼容性、系统稳定性等的问题代码将可能被合入。

fotff提供了一个框架，不断地对最新持续集成版本运行测试套，然后对其中失败用例进行分析：找到或生成在该用例上次通过的持续集成版本和本次失败的持续集成版本之间的所有中间版本，然后运用二分法的思想，找到出现该问题的第一个中间版本，从而给出引入该问题的代码提交。

#### 软件架构

```
fotff
├── pkg       # 版本包管理的接口定义和特定开发板形态的具体实现
├── rec       # 测试结果记录和分析
├── tester    # 测试套的接口定义和调用测试框架的具体实现
├── utils     # 一些通用的类库
├── vcs       # 版本控制相关的包，比如manifest的处理，通过OpenAPI访问gitee查询元数据的函数等
├── fotff.ini # 运行需要的必要参数配置，比如指定测试套、配置构建服务器、HTTP代理等
└── main.go   # 框架入口
```

#### 安装教程

1. 获取[GoSDK](https://golang.google.cn/dl/)并按照指引安装。
2. 在代码工程根目录执行```go build```编译。如下载依赖库出现网络问题，必要时配置GOPROXY代理。
3. 更改fotff.ini，按功能需要，选择版本包和测试套的具体实现，完成对应参数配置，并将可能涉及到的测试用例集、脚本、刷机工具等放置到对应位置。
4. fotff所有参数均通过ini文件管理，执行二进制不需要加其他命令行参数。

#### 使用说明

1. 分析结果在.fotff/records.json文件中记录；如果配置了邮箱信息，会发送结果到指定邮箱。
2. 刷机、测试具体实现可能涉及到[hdc_std](https://gitee.com/openharmony/developtools_hdc)、[xdevice](https://gitee.com/openharmony/testfwk_xdevice)，安装和配置请参考对应工具的相关页面。
3. xdevice运行需要Python运行环境，请提前安装。
4. 刷机、测试过程需要对应开发板的驱动程序，请提前安装。

#### 参与贡献

1. Fork 本仓库
2. 新建 Feat_xxx 分支
3. 提交代码
4. 新建 Pull Request

#### 相关链接

[OpenHarmony CI](http://ci.openharmony.cn/dailys/dailybuilds)

[developtools_hdc](https://gitee.com/openharmony/developtools_hdc)

[hihope](https://gitee.com/hihope_iot/docs/tree/master/HiHope_DAYU200/烧写工具及指南)

[testfwk_xdevice](https://gitee.com/openharmony/testfwk_xdevice)
