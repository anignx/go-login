package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"my.service/go-login/package/config"
)

type HttpClient struct {
	ServerMap map[string]*http.Client
}

func New() *HttpClient {
	return &HttpClient{}
}

func (h *HttpClient) InitHttpClient() error {
	h.ServerMap = make(map[string]*http.Client)
	for _, s := range config.Conf.ServerClient {
		timeout := time.Duration(s.ReadTimeout) * time.Millisecond
		httpClient := &http.Client{
			Timeout: timeout,
			// 返回301、302重定向时，不会自动发起重定向访问
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// 不校验https证书
					InsecureSkipVerify: true,
				},
				MaxConnsPerHost:     300,
				MaxIdleConns:        150,
				MaxIdleConnsPerHost: 75,
				IdleConnTimeout:     10 * time.Second,
			},
		}
		h.ServerMap[s.ServiceName] = httpClient
	}
	return nil
}

var proxy = New()

func InitHttp(serviceName string) *http.Client {
	if _, err := proxy.ServerMap[serviceName]; err {
		fmt.Println("serviceName is not exist", serviceName)
		return nil
	}
	return proxy.ServerMap[serviceName]
}

func HttpPostWithRetry(client *http.Client, url, body string, retryTime int) (result []byte, err error) {
	if retryTime < 0 {
		return nil, errors.New("HttpPostWithRetry retryTime is " + strconv.Itoa(retryTime))
	}
	for i := 0; i < retryTime; i++ {
		if result, err = HttpPostClient(client, url, body); err == nil {
			break
		}
	}
	return
}

func HttpGetWithRetry(client *http.Client, url string, retryTime int) (result []byte, err error) {
	if retryTime < 0 {
		return nil, errors.New("HttpGetWithRetry retryTime is " + strconv.Itoa(retryTime))
	}
	for i := 0; i < retryTime; i++ {
		if result, err = HttpGetClient(client, url); err == nil {
			break
		}
	}
	return
}

func HttpPostClient(client *http.Client, url, body string) ([]byte, error) {
	resp, err := client.Post(url, gin.MIMEPOSTForm, strings.NewReader(body))
	// 重定向时，resp不为nil，如果CheckRedirect返回非ErrUseLastResponse错误，err也不为nil
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func HttpGetClient(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	// 重定向时，resp不为nil，如果CheckRedirect返回非ErrUseLastResponse错误，err也不为nil
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// 非服务发现可以使用下面的函数进行http请求
var httpClient = &http.Client{
	Timeout: 3 * time.Second,
	// 返回301、302重定向时，不会自动发起重定向访问
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			// 不校验https证书
			InsecureSkipVerify: true,
		},
		MaxConnsPerHost:     300,
		MaxIdleConns:        150,
		MaxIdleConnsPerHost: 75,
		IdleConnTimeout:     10 * time.Second,
	},
}

func HttpPost(url, body string) ([]byte, error) {
	resp, err := httpClient.Post(url, gin.MIMEPOSTForm, strings.NewReader(body))
	// 重定向时，resp不为nil，如果CheckRedirect返回非ErrUseLastResponse错误，err也不为nil
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func HttpGet(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	// 重定向时，resp不为nil，如果CheckRedirect返回非ErrUseLastResponse错误，err也不为nil
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}
