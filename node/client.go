package node

import (
	"time"
)

type nodeInfo struct {
	conn []string
}
type connect struct {
	url  string
	name string
}
type rec struct {
	index       int
	dur         time.Duration
	blockNumber int64
}

func Init(connInfo []string) *nodeInfo {
	return &nodeInfo{
		conn: connInfo,
	}
}

func (node *nodeInfo) BTCConnect() *connect {
	return nil
}
func (node *nodeInfo) ETHConnect() *connect {
	var conn connect
	// var r rec
	// flag := true
	// for i := 0; i < len(node.conn); i++ {
	// 	start := time.Now()
	// 	blockNum, _ := eth.Ethrpc_blockNumber(node.conn[i])
	// 	// fmt.Println("blockNum:", blockNum)
	// 	if blockNum > 0 {
	// 		elapsed := time.Since(start)
	// 		// fmt.Println("elapsed:", elapsed)
	// 		if flag {
	// 			r.index = i
	// 			r.dur = elapsed
	// 			r.blockNumber = blockNum
	// 			conn.url = node.conn[i]
	// 			conn.name = "ETH"
	// 			flag = false
	// 		} else {
	// 			if r.blockNumber <= blockNum {
	// 				// fmt.Println("#1")
	// 				if elapsed < r.dur {
	// 					// fmt.Println("#2")
	// 					r.index = i
	// 					r.dur = elapsed
	// 					r.blockNumber = blockNum
	// 					conn.url = node.conn[i]
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	return &connect{
		url:  conn.url,
		name: conn.name,
	}
}
func (node *nodeInfo) DDMXConnect() *connect {
	var conn connect
	// var r rec
	// flag := true
	// for i := 0; i < len(node.conn); i++ {
	// 	start := time.Now()
	// 	blockNum, _ := ddmx.Ddmxrpc_blockNumber(node.conn[i])
	// 	// fmt.Println("blockNum:", blockNum)
	// 	if blockNum > 0 {
	// 		elapsed := time.Since(start)
	// 		// fmt.Println("elapsed:", elapsed)
	// 		if flag {
	// 			r.index = i
	// 			r.dur = elapsed
	// 			r.blockNumber = blockNum
	// 			conn.url = node.conn[i]
	// 			conn.name = "DDMX"
	// 			flag = false
	// 		} else {
	// 			if r.blockNumber <= blockNum {
	// 				// fmt.Println("#1")
	// 				if elapsed < r.dur {
	// 					// fmt.Println("#2")
	// 					r.index = i
	// 					r.dur = elapsed
	// 					r.blockNumber = blockNum
	// 					conn.url = node.conn[i]
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	return &connect{
		url:  conn.url,
		name: conn.name,
	}
}
func (node *nodeInfo) GPEConnect() *connect {
	var conn connect
	// var r rec
	// flag := true
	// for i := 0; i < len(node.conn); i++ {
	// 	start := time.Now()
	// 	blockNum, _ := gpe.Gperpc_blockNumber(node.conn[i])
	// 	// fmt.Println("blockNum:", blockNum)
	// 	if blockNum > 0 {
	// 		elapsed := time.Since(start)
	// 		// fmt.Println("elapsed:", elapsed)
	// 		if flag {
	// 			r.index = i
	// 			r.dur = elapsed
	// 			r.blockNumber = blockNum
	// 			conn.url = node.conn[i]
	// 			conn.name = "GPE"
	// 			flag = false
	// 		} else {
	// 			if r.blockNumber <= blockNum {
	// 				// fmt.Println("#1")
	// 				if elapsed < r.dur {
	// 					// fmt.Println("#2")
	// 					r.index = i
	// 					r.dur = elapsed
	// 					r.blockNumber = blockNum
	// 					conn.url = node.conn[i]
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	return &connect{
		url:  conn.url,
		name: conn.name,
	}
}
func (node *nodeInfo) USDTConnect() *connect {
	return nil
}
func (node *nodeInfo) RegConnect() *connect {
	if len(node.conn) <= 0 {
		return &connect{
			url:  "",
			name: "reg",
		}
	}
	return &connect{
		url:  node.conn[0],
		name: "reg",
	}
}
func (node *nodeInfo) PayConnect() *connect {
	if len(node.conn) <= 0 {
		return &connect{
			url:  "",
			name: "pay",
		}
	}
	return &connect{
		url:  node.conn[0],
		name: "pay",
	}
}
func (node *nodeInfo) PwdConnect() *connect {
	if len(node.conn) <= 0 {
		return &connect{
			url:  "",
			name: "pwd",
		}
	}
	return &connect{
		url:  node.conn[0],
		name: "pwd",
	}
}

func (con *connect) URL() string {
	return con.url
}
func (con *connect) Name() string {
	return con.name
}
