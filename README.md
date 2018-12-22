手册
---------------

一些固定的规则
-------------------

#### 返回的数据格式固定

成功
```
{
    "result":false,
    "mcode":"PARAMETER_ERROR",
    "message":"invalid parameter,'xxx' is empty",
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
    }
}
```

协议传输的对象放在pb/xxxpb中,使用protobuf来生成对象,如果不需要omitempty,则可以重新编译proto-gen-go或使用脚本去除