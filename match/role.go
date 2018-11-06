package match

// import (
// 	tg "github.com/hiank/think/net/protobuf/grpc"
// )

// Role 玩家结构
type Role interface {

	GetToken() string 						//NOTE: 用于获取 玩家对应的token，服务中唯一
	GetId() string							//NOTE: 服务器唯一，玩家id
	GetKey() int 							//NOTE: 用于获取匹配关键字
	// GetResponseChan() chan *tg.Response		//NOTE: 用于获取 通知战局 的chan
}