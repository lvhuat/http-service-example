手册
---------------



#### 返回的数据格式固定

成功
```
{
    "result":false,
    "mcode":"PARAMETER_ERROR",
    "message":"invalid parameter,'xxx' is empty",
    "timestamp":1545471685000
}
```

失败
```
{
    "result":true,
    "data":{
        "userId":"xiaoming_2018",
        "birthday":"1995-06-11",
        "email":"demo@lwork.com",
        "mobile":"+861888888888"
        "address":"四川省成都市XXXX",
        "timestamp":1545471685000
    }
}
```

目录和主要文件
------------
- common 存放一些工具和一些常量  
- conf 配置解析  
- dto 数据库模型和操作  
- httpserver http服务  
- invokes 内部其他微服务的接口调用  
- pb 接口传输模型  
- rcontxt 接口上下文,包含了数据的解析,接口参数逻辑无关的判断  
- service 逻辑整合  
- testcmd 集成测试命令行工具  
- app.toml 配置文件  
- Dockerfile docker文件  
- main.go  