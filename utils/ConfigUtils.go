package utils

import (
	"fmt"
	"github.com/beego/beego/v2/core/config/env"
	"github.com/beego/beego/v2/server/web"
	"strconv"
)

var (
	_ Config = (*BaseConfig)(nil)
	_ Config = (*EnvConfig)(nil)
	_ Config = (*FileConfig)(nil)
)

type Config interface {
	SetNextConfig(config Config) *BaseConfig
	String(key string, defaultValue string) string
	innerString(key string) (string, error)
	Int64(key string, defaultValue int64) int64
	Int(key string, defaultValue int) int
	Float(key string, defaultValue float64) float64
	Bool(key string, defaultValue bool) bool
}
type BaseConfig struct {
	NextConfig Config
	self       Config
}

func (this *FileConfig) setSelf() Config {
	this.self = this
	return this
}
func (this *EnvConfig) setSelf() Config {
	this.self = this
	return this
}
func (this *BaseConfig) Bool(key string, defaultValue bool) bool {
	if val, err := strconv.ParseBool(this.String(key, fmt.Sprintf("%t", defaultValue))); err == nil {
		return val
	}
	return defaultValue
}
func (this *BaseConfig) Int64(key string, defaultValue int64) int64 {
	if val, err := strconv.ParseInt(this.String(key, fmt.Sprintf("%d", defaultValue)), 10, 16); err == nil {
		return val
	}
	return defaultValue
}
func (this *BaseConfig) Int(key string, defaultValue int) int {
	if val, err := strconv.ParseInt(this.String(key, fmt.Sprintf("%d", defaultValue)), 10, 16); err == nil {
		return int(val)
	}
	return defaultValue
}
func (this *BaseConfig) Float(key string, defaultValue float64) float64 {
	if val, err := strconv.ParseFloat(this.String(key, fmt.Sprintf("%f", defaultValue)), 64); err == nil {
		return val
	}
	return defaultValue
}

type EnvConfig struct {
	BaseConfig
}

func (this *BaseConfig) innerString(key string) (string, error) {
	panic("to do")
}
func (this *BaseConfig) String(key string, defaultValue string) string {
	value, err := this.self.innerString(key)
	if err == nil && len(value) != 0 {
		return value
	}
	if this.NextConfig == nil {
		return defaultValue
	}
	return this.NextConfig.String(key, defaultValue)
}

func (this *BaseConfig) SetNextConfig(config Config) *BaseConfig {
	this.NextConfig = config
	return this
}

type FileConfig struct {
	BaseConfig
}

var config Config

func GetConfig() Config {
	if config != nil {
		return config
	}
	//config = new(EnvConfig)
	//config.SetNextConfig(new(FileConfig).SetNextConfig(nil)).setSelf(config)
	//return config
	envConfig := new(EnvConfig).
		setSelf().
		SetNextConfig(new(FileConfig).
			setSelf().
			SetNextConfig(nil))
	return envConfig
}
func (this *FileConfig) innerString(key string) (string, error) {
	return web.AppConfig.String(key)
}
func (this *EnvConfig) innerString(key string) (string, error) {
	return env.MustGet(key)
}
