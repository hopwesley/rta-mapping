
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

### 修改 ID Mapping

- **入参**: 详细描述输入参数。
- **结果**: 详细描述返回结果。

### 修改 RTA 命令

- **入参**: 详细描述输入参数。
- **结果**: 详细描述返回结果。

