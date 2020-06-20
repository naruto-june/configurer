## 背景
```
    开发者经常会遇到不同的环境使用不同配置，或者即使是同样的环境（同一个主机）但却需
要运行多个不同配置的实例，这就导致了部署时
的混乱局面，开发者也往往被搞得会晕头转。本项目借鉴docker镜像仓库自适应不同cpu架构而进行智能下载镜像功能，在configurer中统一
管理不同环境不同配置，旨在解决环境不同而配置不同，引发部署时的混乱问题。
```

## 目录结构说明
```
    |
    |----online: 线上环境的所有配置汇集处 (测试用例中使用）
    |        |-------kafka_consumer_conf.json        消费者配置 kafka_consumer_conf1.json 和   kafka_consumer_conf2.json 同环境下的扩展配置
    |        |-------kafka_monitor_conf.json         监控配置 kafka_monitor_conf1.json 和 kafka_monitor_conf2.json 同环境下的扩展配置
    |        |-------kafka_producer_conf.json        生产者配置 kafka_producer_conf1.json 和 kafka_producer_conf2.json  同环境下的扩展配置
    |
    |----sandbox: 测试环境的所有配置汇集处  (测试用例中使用）
    |        |-------kafka_consumer_conf.json        消费者配置 kafka_consumer_conf1.json 和   kafka_consumer_conf2.json 同环境下的扩展配置
    |        |-------kafka_monitor_conf.json         监控配置 kafka_monitor_conf1.json 和 kafka_monitor_conf2.json 同环境下的扩展配置
    |        |-------kafka_producer_conf.json        生产者配置 kafka_producer_conf1.json 和 kafka_producer_conf2.json  同环境下的扩展配置
    |----seelog.xml 不同环境的共用配置文件 (测试用例中使用）
    |
    |----configurer.json/configurer.toml/configurer.yaml 不同环境不同配置的统一管理配置处   (注意：需要开发者根据自己项目编写配置，并置于自己项目根目录下）
    |----configurer.go + parse.go  项目库的代码
    |----configurer_test.go  三种格式配置文件的测试用例
```

## 示例
```
    package main

    import "github.com/naruto-june/configurer"

    func main() {
        err := configurer.ParseConf("configurer.json")
        if nil != err {
            panic(err)
        }

        // 获取公共配置项
        msize,err := configurer.GetCommonConfItem("single_log_max_size")
        ...

        var fContent string
        //获取指定公共配置文件的具体内容 指定文件名称是无扩展名,无后缀的文件名称
        fContent,err = configurer.GetConfByFName("seelog")
        ...
        // 根据具体业务解析fContent得到应用的具体配置
        ...

        //获取个性化配置文件的具体内容 指定文件名称是无后缀的文件名称 比如：kafka_consumer_conf.json 和 kafka_consumer_conf1.json 和   kafka_consumer_conf2.json 使用 kafka_consumer_conf获取
        fContent，err = GetConfByFName("kafka_consumer_conf")
        // 根据具体业务解析fContent得到应用的具体配置
        ...
    }
```


## 联系方式
```
1643127918@qq.com
```

