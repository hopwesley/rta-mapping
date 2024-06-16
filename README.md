
# rta-mapping

## config.json 文件配置介绍

```json
{
    "use_ssl": false, 
    "srv_port": "8801",
    "ssl_cert_file": "", 
    "ssl_key_file": "",
    "redis": { 
        "addr": "localhost:6379",
        "password": "123"
    },
    "mysql": {
        "user_name": "marketing",
        "password": "",
        "host": "",
        "port": "",
        "database": "marketing",
        "limit": 20000000
    }
}
```

### 配置说明

- **use_ssl**: 是否使用SSL证书独立对外提供服务，布尔类型。
- **srv_port**: 对外提供服务的端口号，字符串类型。
- **ssl_cert_file**: 如果使用SSL证书，则提供证书存储目录，字符串类型。
- **ssl_key_file**: SSL证书的key的存储目录，字符串类型。
- **redis**: Redis配置，包含以下字段：
    - **addr**: Redis地址，字符串类型。
    - **password**: Redis密码，字符串类型。
- **mysql**: MySQL配置，包含以下字段：
    - **user_name**: 用户名，字符串类型。
    - **password**: 密码，字符串类型。
    - **host**: 主机地址，字符串类型。
    - **port**: 端口号，字符串类型。
    - **database**: 数据库名称，字符串类型。
    - **limit**: 限制，整数类型。



## 启动方式及命令介绍
### - **编译版本**:
- **编译proto 接口文件**
```
make pbs
```
- **编译可执行文件**
```
make mac/linux
```

### - **版本查看**:   
```
 rtaHint.(mac/linux) -v

==================================================

Version:        fd576e6f8f77d85c6a2db0bd071eb069e9fea457
Build:          2024-06-16_15:58:41_+0800
Commit:         fd576e6f8f77d85c6a2db0bd071eb069e9fea457

==================================================
```

### - **根据配置文件启动**:

```
 ./rtaHint.(mac/linux) -c config.json
```
## API 接口介绍
### 参考测试文件：[srv_test.go](srv_test.go)
### - **对外接口，device与rta的关系**:
- **详细描述输入参数**
[rta_api.proto](common/rta_api.proto)

 - **api url： /rta_hint**
 - 参考
```
req := &common.Req{
    Device: &common.Device{
        ImeiMd5:      "15d35cced5fb....0fbfe76ac626df6",
        Oaid:         "0hnMDI....u7XlLZFY",
        AndroidIdMd5: "puQAjHRIDGN....77eZdm4I6NQJQcD",
    },
    ReqId:  "xxx-xxx-xxx",
    RtaIds: []int64{10003, 10004, 10005},
}
	
var response = &Rsp{
    StatusCode: HitSuccess,
    BidType:    &wrapperspb.Int32Value{Value: BidTypeOk},
    UserInfos:  uis,
    ReqId:      request.ReqId,
}
```

### 修改 Rta 与 user 映射
- **api url： /rta_update**
- **入参**:

```
type RtaUpdateItem struct {
  RtaID   int64 `json:"rta_id"`
  UserIDs []int `json:"user_ids"`
  IsDel   bool  `json:"is_del"`
}
```

- **结果**:

```
type JsonResponse struct {
  Success bool   `json:"success"`
  Code    int    `json:"code"`
  Msg     string `json:"msg"`
}
```
### 修改 ID Mapping

- **api url： /id_map_update**
- **入参**:

```
type JsonRequest struct {
	UserID       int    `json:"user_id"`
	IMEIMD5      string `json:"imei_md5"`
	OAID         string `json:"oaid"`
	IDFA         string `json:"idfa"`
	AndroidIDMD5 string `json:"android_id_md5"`
}
```

- **结果**:

```
type JsonResponse struct {
  Success bool   `json:"success"`
  Code    int    `json:"code"`
  Msg     string `json:"msg"`
}
```

### 修改 RTA 命令

- **入参**: 详细描述输入参数。
- **结果**: 详细描述返回结果。

