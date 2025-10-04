package util

import (
	"crypto/md5"
	"fmt"
	"github.com/denisbrodbeck/machineid"
)

// GenerateHostID 基于硬件信息生成一个稳定的、唯一的标识符。
func GenerateHostID() (string, error) {
	// 如果你不需要应用隔离，只想获取原始ID的哈希值，也可以这样做：
	originalID, err := machineid.ID()
	if err != nil {
		return "", fmt.Errorf("the host ID cannot be generated: %w", err)
	}
	hash := md5.Sum([]byte(originalID))
	id := fmt.Sprintf("%x", hash)
	return id, nil
}
