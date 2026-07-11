package connect

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// client 全局的HTTP客户端
var client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
	Timeout: 2 * time.Second,
}

// CheckURL 检查url链接是否真实有效
func CheckURL(url string) (bool, error) {
	// 构建一个 HEAD 请求
	req, err := http.NewRequestWithContext(context.Background(), http.MethodHead, url, nil)

	if err != nil {
		return false, fmt.Errorf("构建测试请求失败: %w", err)
	}

	// 发起请求
	resp, err := client.Do(req)
	if err != nil { //网络层面问题，比如DNS解析失败、TCP连接被拒
		return false, err
	}
	defer resp.Body.Close()

	// 判断状态码,这里重定向的链接判定为不可用
	return resp.StatusCode == http.StatusOK, nil
}
