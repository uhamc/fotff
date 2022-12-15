# fotff

#### 介绍

fotff(find out the first fault)是为OpenHarmony持续集成设计的问题自动化问题分析工具。

为了平衡开销与收益，考虑到开发效率、资源占用等因素影响，OpenHarmony代码合入门禁（冒烟测试）只拦截部分严重基础问题（例如开机失败、关键进程崩溃、UX布局严重错乱、电话/相机基础功能不可用等）。因此，一些会影响到更细节功能、影响兼容性、系统稳定性等的问题代码将可能被合入。

fotff提供了一个框架，不断地对最新持续集成版本运行测试套，然后对其中失败用例进行分析：找到或生成在该用例上次通过的持续集成版本和本次失败的持续集成版本之间的所有中间版本，然后运用二分法的思想，找到出现该问题的第一个中间版本，从而给出引入该问题的代码提交。

#### 软件架构

```
fotff
├── .fotff    # 缓存等程序运行时产生的文件的存放目录
├── logs      # 日志存放目录
├── pkg       # 版本包管理的接口定义和特定开发板形态的具体实现
├── rec       # 测试结果记录和分析
├── tester    # 测试套的接口定义和调用测试框架的具体实现
├── utils     # 一些通用的类库
├── vcs       # 版本控制相关的包，比如manifest的处理，通过OpenAPI访问gitee查询信息的函数等
├── fotff.ini # 运行需要的必要参数配置，比如指定测试套、配置构建服务器、HTTP代理等
└── main.go   # 框架入口
```

#### 安装教程

1. 获取[GoSDK](https://golang.google.cn/dl/)并按照指引安装。
2. 在代码工程根目录执行```go build```编译。如下载依赖库出现网络问题，必要时配置GOPROXY代理。
3. 更改fotff.ini，按功能需要，选择版本包和测试套的具体实现，完成对应参数配置，并将可能涉及到的测试用例集、脚本、刷机工具等放置到对应位置。

#### 使用说明

###### 普通模式

example: ```fotff```

1. 配置好fotff.ini文件后，不指定任何命令行参数直接执行二进制，即进入普通模式。此模式下，框架会自动不断地获取最新持续集成版本，并对其运行测试套，然后对其中失败用例进行分析。
2. 分析结果在.fotff/records.json文件中记录；如果配置了邮箱信息，会发送结果到指定邮箱。

###### 对单个用例在指定区间内查找

example: ```fotff run -s pkgDir1 -f pkgDir2 -t TEST_CASE_001```

1. 配置好fotff.ini文件后，通过-s/-f/-t参数在命令行中分别指定成功版本/失败版本/测试用例名，即可对单个用例在指定区间内查找。此模式下，仅在指定的两个版本间进行二分查找，运行指定的运行测试用例。
2. 分析结果在控制台中打印，不会发送邮件。

###### 烧写指定版本包

example: ```fotff flash pkgDir```

1. 配置好fotff.ini文件后，可以指定版本包目录烧写对应版本。
2. 版本包目录指的是有fotff生成的目录。如果使用自行解压的目录，可能需要对应修改。例如dayu200的版本包，需要在目录下添加__built__空文件，以跳过构建步骤。否则会基于manifest_tag.xml重新构建。


###### tips

1. 刷机、测试具体实现可能涉及到[hdc_std](https://gitee.com/openharmony/developtools_hdc)、[xdevice](https://gitee.com/openharmony/testfwk_xdevice)，安装和配置请参考对应工具的相关页面。
2. xdevice运行需要Python运行环境，请提前安装。
3. 刷机、测试过程需要对应开发板的驱动程序，请提前安装。

#### 参与贡献

1. Fork 本仓库
2. 新建 Feat_xxx 分支
3. 提交代码
4. 新建 Pull Request

#### 相关链接

[OpenHarmony CI](http://ci.openharmony.cn/dailys/dailybuilds)

[developtools_hdc](https://gitee.com/openharmony/developtools_hdc)

[dayu200_tools](https://gitee.com/hihope_iot/docs/tree/master/HiHope_DAYU200/烧写工具及指南)

[testfwk_xdevice](https://gitee.com/openharmony/testfwk_xdevice)
