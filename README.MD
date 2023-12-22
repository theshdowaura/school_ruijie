# 简介

本脚本实现了锐捷 ePortal Web 认证的登录，需要 curl 指令的支持，通过 curl 发送认证所需的 POST 数据包来实现 Web 认证。

如果你所在的学校也是使用锐捷进行校园网认证，并且开启了 Web 的认证方式（如下图），那么应该是可以使用本脚本进行认证的。

![ruijie_web](images/ruijie_web.png)

项目中的脚本文件以 **_jmu** 结尾命名的仅适用于集美大学，其他学校请使用 **_general** 结尾命名的脚本文件。



欢迎大家发起 Pull request 提交各自学校的认证脚本！

**脚本依赖于 curl，所以使用之前请确保系统中已经安装了 curl。**

# 脚本工作原理

这个 Shell script 主要工作原理如下：

1. 向 [Captive Portal Server](https://en.wikipedia.org/wiki/Captive_portal) 发送 GET 请求，检查返回的 HTTP 状态码是否为 204：
   1. 如果为 204 则说明当前已经能正常联网，退出脚本，停止后续认证流程；
   2. 如果不是 204 则说明当前未进行锐捷认证（一般来说，在未认证的情况下请求任意页面，认证服务器都会返回一个 HTTP 302 或者 301 状态码，用于将你重定向到 ePortal 认证页面），则将 HTTP 重定向所转向的页面的 URL 作为认证页面 URL。
2. 构造进行锐捷 ePortal 认证所需的 HTTP 请求，并通过 curl 发送该请求。
3. 认证成功。

# 基本用法
** 请注意 linux与macos请授予执行权限 **
```shell
chmod +x ruijie_{platform}_{Architecture}
```
shell脚本:
```shell
./ruijie_template.sh username password code
```
windows amd64平台使用:
```cmd
ruijie_windows_amd64.exe username password code
```
macos amd64平台使用:
```cmd
./ruijie_darwin_amd64 username password code
```
- username 参数为认证用户名。
- password 参数为认证密码。
- code 参数为运营商密码（无则忽略）

## 图形化客户端
在fyne分支上请自行编译客户端,master分支仅发布内核

# 局限性

在锐捷 ePortal Web 认证的 POST 数据包中发现了 `validcode` 参数，为验证码。我所在的学校并没有强制在认证的时候要求输入验证码，但是经过测试发现：如果多次认证失败（例如用户名或密码错误）或者多次重复刷新认证页面，则认证过程需要输入页面上显示的验证码。

如果你所在的学校需要在认证过程中输入验证码，则不适合使用这个 shell script 进行认证，需要使用 Python 等脚本语言实现验证码识别功能再进行认证。




# 后记

锐捷的 ePortal Web 认证上了之后，向 `InterFace.do?method=getOnlineUserInfo` 这个接口中 Post 提交 `userIndex` 参数会明文显示用户的密码(部分学校已修复)，感觉一点安全性都没有。（集美大学在认证系统升级之后已修复该问题）

我所在的学校在未认证的情况下，53端口(DNS)似乎是可以正常通信的(我这里不支持隧道)。如果你所在的学校也有类似的情况，也许可以尝试搭建个53端口的代理实现免认证上网。

# 友情链接
https://github.com/callmeliwen/RuijiePortalLoginTool

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=theshdowaura/school_ruijie&type=Date)](https://star-history.com/#theshdowaura/school_ruijie&Date)
