package utils

import (
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

// PostgresListener PostgreSQL监听器结构体
type PostgresListener struct {
	Dsn       string
	listeners []*pq.Listener
	channels  []string
}

var PostListener *PostgresListener

type MsgCallBack func(channelName string, msg string)

func ConnectPostgresMsg(dsnVarName string, channels string, callbacks ...MsgCallBack) {
	var dsn string
	var err error
	//var dbDriverName string
	if dsn = GetConfig().String(dsnVarName, ""); dsn == "" {
		Logger.Error(err.Error())
		return
	}
	PostListener = new(PostgresListener)
	PostListener.Dsn = dsn
	channelsTxt := GetConfig().String(channels, "")
	if channelsTxt == "" {
		return
	}
	PostListener.channels = strings.Split(channelsTxt, ",")
	if len(PostListener.channels) == 0 {
		return
	}
	if len(callbacks) != len(PostListener.channels) {
		Logger.Error("callback length not match channel length")
		return
	}
	for i, channel := range PostListener.channels {
		err = PostListener.ListenToChannel(channel, callbacks[i])
		if err != nil {
			Logger.Error(err.Error())
		}
	}

}

// ListenToChannel 侦听PostgreSQL通道中的消息
func (pl *PostgresListener) ListenToChannel(channelName string, callback MsgCallBack) error {
	// 创建一个新的监听器
	listener := pq.NewListener(pl.Dsn, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("监听器错误: %v\n", err)
		}
	})
	pl.listeners = append(pl.listeners, listener)

	// 开始监听指定通道
	err := listener.Listen(channelName)
	if err != nil {
		return err
	}

	// 在后台持续处理通知
	go func() {
		for {
			select {
			case n := <-listener.Notify:
				if n != nil {
					callback(n.Channel, n.Extra)
				}
			case <-time.After(time.Minute):
				// 超时检查，确保连接存活
				go func() {
					err := listener.Ping()
					if err != nil {
						Logger.Error("listener msg via "+channelName, zap.Error(err))
						callback("", channelName)
						listener.Close()
					}
				}()
			}
		}
	}()

	return nil
}

// Close 关闭监听器连接
func (pl *PostgresListener) Close() error {
	for _, listener := range pl.listeners {
		err := listener.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
