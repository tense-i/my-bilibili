package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	// Elasticsearch 集群配置
	Elasticsearch struct {
		DmExternal struct {
			Addresses []string
		}
		ReplyExternal struct {
			Addresses []string
		}
		ExternalPublic struct {
			Addresses []string
		}
	}

	// 搜索服务 Redis 配置 (使用不同名称避免冲突)
	SearchRedis struct {
		Host string
		Type string
		Pass string
	}
}
