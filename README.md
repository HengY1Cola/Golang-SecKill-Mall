# 仿商城秒杀系统(Golang练手项目)

**寒假无聊🥱，为了准备实习**，想着除了了日常的比赛等项目～

去慕课网上找了找不到Go的实战，**年前**已经搞定了`微服务网关`实战

**年后**顺便也把这个做了。 姊妹篇地址：[用Go实现微服务企业网关(强烈推荐要掌握)](https://github.com/HengY1Sky/Golang-Gate-Way)

这个项目还是挺**推荐**的。有`ThinkPhp`的基础，`Iris`相当于就是Go版本了🚀

---

这个项目是前后端一起写，返回的大多数`Html`（当然接口也是可以的）

主要分为了`前台`，`后台`，`秒杀`这个**三个模块**，大多数比较基础但是还是能学到思想

下面简单说一下**技术栈**，**架构图**(重新手画)，**目录**，**学习经验**，**如何使用**。

##  技术栈/架构图

> 其实里面的分布式验证可以用Redis实现 但是只要 一致性哈希 这个功能完全可以用代码层面实现
>
> SLB 业务不大的情况下 完全可以自己配置个组 但是在一个内网做均衡会快很多

|        技术         |                           相关文档                           |
| :-----------------: | :----------------------------------------------------------: |
| 消息队列： RabbitMQ | [自己总结了下->快速上手](https://blog.csdn.net/weixin_51485807/article/details/122761910) |
|    MVC框架：Iris    |          [中文文档](https://www.topgoer.com/Iris/)           |
|   静态加速： CDN    | [腾讯云CDN的详细配置](https://cloud.tencent.com/developer/article/1462593?from=15425) |
|    负载均衡：SLB    |       [阿里云SLB](https://www.aliyun.com/product/slb)        |
|    数据库：Mysql    | [Mysql5.7中文文档](https://www.docs4dev.com/docs/zh/mysql/5.7/reference/) |
|  数据库映射：GORM   |    [GORM中文文档](https://gorm.io/zh_CN/docs/index.html)     |
|    压力测试：Wrk    | [性能测试神器 wrk 使用教程](https://segmentfault.com/a/1190000023212126) |

![](https://raw.githubusercontent.com/HengY1Sky/GoSecKillMall/main/GoFramework.webp)

##  教学目录/框架目录

*教学地址：https://coding.imooc.com/class/chapter/347.html

```
├── backend
│   ├── main.go # 后台启动文件
│   └── web
├── common # 公共库
├── fronted
│   ├── main.go # 前端启动
│   ├── middleware # 前端中间件
│   ├── productMain.go # CDN入口
│   └── web # 逻辑
├── rabbitmq # rabbitmq
├── datamodels # 基础模型
├── repositories # 面向SQL
├── services # 服务
├── consumer.go # 抢购队列消费
├── getOne.go # 判断是否超卖
├── go.mod
├── go.sum
└── validate.go
```

```
├── 秒杀系统需求整理&系统设计
├── 环境搭建之初识RabbitMQ
├── 环境搭建之Iris框架入门
├── 后台管理功能开发之商品管理功能开发
├── 后台管理功能开发之订单功能开发
├── 秒杀前台功能开发 之用户注册登录功能开发
├── 秒杀前台功能开发之商品展示及数据控制功能开发
├── 秒杀系统分析&前端优化
├── 服务端性能优化之实现cookie验证
├── 服务端性能优化之分布式验证实现
├── 服务端性能优化解决超卖&引入消息队列
└── 秒杀安全优化
```

## 经验/更新

1. **第三章**的知识点我总结成了一篇博客 👉 [点击这里查看](https://blog.csdn.net/weixin_51485807/article/details/122761910)
2. **第五章**，**第六章**的代码就是换了个皮，意义不大，也就没补充功能。
3. **第五章**，**第六章**主要是后台代码，**第七章，第八章**是前台代码。
4. **第九章**生成静态文件使后面的访问不会打到数据库与CDN技术
5. **第十章**就讲了个将AES做了个加密
6. **第十一章**用一致性哈希与抽出单独的validate方便横向拓展(偏**干货**) [简了解一致性哈希](https://segmentfault.com/a/1190000021199728)
7. **第十二章**实现消息队列以及完成了基本的抢购逻辑
8. **第十三章**主要是给了个例子横向拓展秒杀验证逻辑

---

1. 使用的是`iris/v12`最新的框架，并且使用了包管理，开箱即用。
2. 将原生的SQL写法更换成了Gorm，更加规范了。
3. 完善抢购逻辑，将整个抢购思路形成闭环。

## 快速开始

> 当然这个只是个学习库，谈不上进行使用
>
> 但是我还是简单写一下如何快速布局，别进去全是红的找不到文件就行了😊

- SQL文件直接到`Navicat`**创建自己的库**然后**创建查询**然后**执行**就好了

```bash
$ export GO111MODULE=on && export GOPROXY=https://goproxy.cn # 开启mod模块以及换源
$ git clone https://github.com/HengY1Sky/GoSecKillMall
$ cd GoSecKillMall
$ go mod tidy
# 进入Goland开始操作吧
```

