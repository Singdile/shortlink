# URL 结构解析

## 完整 URL 示例

```
https://username:password@example.com:8080/path/to/resource?key1=value1&key2=value2#section-1
```

## 结构图解

```
  scheme    user info    host      port      path              query           fragment
    │          │          │         │         │                 │                │
    ▼          ▼          ▼         ▼         ▼                 ▼                ▼
https://username:password@example.com:8080/path/to/resource?key1=value1&key2=value2#section-1
    │          │          │                    │                       │              │
    │          └──────────┴────────────────────┴───────────────────────┘              │
    │                              authority                                          │
    │                                                                                 
    └─────────────────────────────────────────────────────────────────────────────────┘
                                       完整 URL
```

## 各部分详解

| 部分 | 名称 | 示例值 | 说明 |
|------|------|--------|------|
| `scheme` | 协议 | `https` | 协议类型，如 `http`、`https`、`ftp`、`ws`、`mailto` 等 |
| `username:password` | 用户信息 | `username:password` | HTTP 认证信息，现代 Web 中很少使用 |
| `host` | 主机 | `example.com` | 域名或 IP 地址 |
| `port` | 端口 | `8080` | 服务端口号，默认值：`http=80`、`https=443` |
| `path` | 路径 | `/path/to/resource` | 服务器上的资源路径 |
| `query` | 查询参数 | `key1=value1&key2=value2` | `?` 开头的键值对，多个参数用 `&` 分隔 |
| `fragment` | 片段标识符 | `section-1` | `#` 开头，用于页面内锚点定位，**不会发送到服务器** |

## Go URL 结构体

```go
type URL struct {
    Scheme      string    // 协议，如 "https"
    Opaque      string    // 用于非层级 URL（如 "mailto:foo@example.com"）
    User        *Userinfo // 用户信息
    Host        string    // 主机名 + 端口，如 "example.com:8080"
    Path        string    // 路径，如 "/path/to/resource"
    RawPath     string    // 编码后的路径
    ForceQuery  bool      // 是否强制包含 "?"
    RawQuery    string    // 查询参数字符串
    Fragment    string    // 片段标识符
    RawFragment string    // 编码后的片段
}
```

## Go 代码示例

### 解析 URL

```go
package main

import (
    "fmt"
    "net/url"
)

func main() {
    rawURL := "https://username:password@example.com:8080/path/to/resource?key1=value1&key2=value2#section-1"
    
    u, err := url.Parse(rawURL)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Scheme:   %s\n", u.Scheme)     // https
    fmt.Printf("User:     %s\n", u.User)       // username:password
    fmt.Printf("Host:     %s\n", u.Host)       // example.com:8080
    fmt.Printf("Path:     %s\n", u.Path)       // /path/to/resource
    fmt.Printf("RawQuery: %s\n", u.RawQuery)   // key1=value1&key2=value2
    fmt.Printf("Fragment: %s\n", u.Fragment)   // section-1
    
    // 解析 Query 参数
    query := u.Query()
    fmt.Printf("key1: %s\n", query.Get("key1")) // value1
    fmt.Printf("key2: %s\n", query.Get("key2")) // value2
}
```

### 构建 URL

```go
package main

import (
    "fmt"
    "net/url"
)

func main() {
    // 方式一：直接构造结构体
    u := &url.URL{
        Scheme:   "https",
        Host:     "example.com",
        Path:     "/api/users",
        RawQuery: "id=123&name=test",
        Fragment: "profile",
    }
    fmt.Println(u.String())
    // 输出: https://example.com/api/users?id=123&name=test#profile
    
    // 方式二：使用 Query() 构建查询参数
    u = &url.URL{
        Scheme: "https",
        Host:   "example.com",
        Path:   "/search",
    }
    query := u.Query()
    query.Set("q", "golang教程")
    query.Set("page", "1")
    u.RawQuery = query.Encode()
    fmt.Println(u.String())
    // 输出: https://example.com/search?q=golang%E6%95%99%E7%A8%8B&page=1
}
```

### URL 编解码

```go
package main

import (
    "fmt"
    "net/url"
)

func main() {
    // 编码
    encoded := url.QueryEscape("你好 世界")
    fmt.Println(encoded) // %E4%BD%A0%E5%A5%BD%20%E4%B8%96%E7%95%8C
    
    // 解码
    decoded, _ := url.QueryUnescape(encoded)
    fmt.Println(decoded) // 你好 世界
    
    // 路径编码
    path := url.PathEscape("/path/中文/文件")
    fmt.Println(path) // /path/%E4%B8%AD%E6%96%87/%E6%96%87%E4%BB%B6
}
```

## 常见 URL 格式

| URL | 说明 |
|-----|------|
| `https://example.com` | 最简形式，只有 scheme + host |
| `https://example.com/path` | 带路径 |
| `https://example.com?key=value` | 只有查询参数，无路径 |
| `https://example.com/#section` | 根路径 + fragment |
| `mailto:test@example.com` | 非 HTTP 协议 |
| `//example.com/path` | 协议相对 URL（继承当前页面协议） |
| `/path/to/resource` | 相对路径 |
| `?key=value` | 只有查询参数的相对 URL |

## 本项目中的应用

在 `pkg/urltool/urltool.go` 中：

```go
func GetBaseUrl(longurl string) (string, error) {
    urlmap, err := url.Parse(longurl)
    if err != nil {
        return "", err
    }
    if len(urlmap.Host) == 0 {
        return "", errors.New("need a valid url with host")
    }
    return path.Base(urlmap.Path), nil
}
```

用于解析长链接，提取路径的最后一部分，判断是否是短链接标识符。

## 参考资料

- [Go net/url 官方文档](https://pkg.go.dev/net/url)
- [RFC 3986 - URI 通用语法](https://datatracker.ietf.org/doc/html/rfc3986)
