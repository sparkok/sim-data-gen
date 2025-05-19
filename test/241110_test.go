package test

import (
	"bytes"
	"context"
	auto "sim_data_gen/auto"
	_ "sim_data_gen/controllers"
	. "sim_data_gen/services"
	"sim_data_gen/simulator"
	. "sim_data_gen/utils"
	"fmt"
	"github.com/beego/beego/v2/core/config/env"
	beegoLogs "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

// setEnv 设置环境变量，并返回一个恢复函数
func setEnv() {
	// 定义需要设置的环境变量
	envVars := map[string]string{
		"dsn":                "../../db/minedb-176.db",
		"dbDriverName":       "sqlite",
		"simulator.workMode": "2",
	}

	for key, value := range envVars {
		env.Set(key, value)
	}
}

func SetEnvs(t *testing.T) {
	// 设置环境变量
	setEnv()
}

// httpTest 封装 HTTP POST 请求的逻辑
func recalculate() error {
	// 定义请求的 URL
	url := "http://127.0.0.1:8888/sim_data_gen/setofboomgroups/recalculate"

	// 定义请求的 JSON 数据
	jsonData := `{
		"date": "2024-11-10",
		"reloadData": false
	}`

	// 创建一个 bytes.Buffer 来存储 JSON 数据
	jsonBody := bytes.NewBuffer([]byte(jsonData))

	// 发起 POST 请求
	resp, err := http.Post(url, "application/json", jsonBody)
	if err != nil {
		return fmt.Errorf("HTTP POST request failed: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	//打印响应体（可选）
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %v", err)
	}
	Logger.Info("Response body:", zap.String("body", string(body)))
	return nil
}
func Test241110(t *testing.T) {
	RunTest(func() {
		if recalculate() != nil {
			t.Fatalf("recalculate failed")
		}
	})
}
func RunTest(testMethod func()) {

	// 设置环境变量
	setEnv()
	ConnectDB("dsn", "dbDriverName")
	InitRemarkRender()
	//暂时注释掉
	//ConnectPostgresMsg("dsn", "dbChannels", GetProductionsContents().OnProductChange)
	AutoMigrate("autoMigrate", auto.AutoMigrateDbList)
	StartProductionsContentService("productions")
	beegoLogs.SetLevel(beegoLogs.LevelInfo)
	simulator.StartSimulator("simulator")

	// 使用 WaitGroup 来等待 HTTP 请求完成
	var wg sync.WaitGroup
	wg.Add(1)

	var testWg sync.WaitGroup
	testWg.Add(1)

	// 启动 Web 服务
	go func() {
		defer wg.Done()
		web.Run()

	}()

	// 等待 Web 服务启动
	time.Sleep(3 * time.Second)
	// 调用封装的 HTTP 请求函数
	testMethod()
	//测试方法结束
	testWg.Done()
	// 创建一个 context 用于 Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Second)
	defer cancel()

	//等待测试方法结束
	testWg.Wait()
	web.BeeApp.Server.Shutdown(ctx)

	// 等待 Web 服务结束
	wg.Wait()
}
