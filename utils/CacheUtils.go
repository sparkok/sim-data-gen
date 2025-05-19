package utils
import (
	"github.com/patrickmn/go-cache"
	"time"
)


func CreateCache()  {
	var expSeconds = int(GetConfig().Int64("cache.exp_seconds",300))
	theCache = cache.New(time.Second*time.Duration(expSeconds), time.Duration(expSeconds*2)*time.Second)
}
var theCache *cache.Cache
func GetCache() *cache.Cache {
	return theCache
}
