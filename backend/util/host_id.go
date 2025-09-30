package util

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
)

// GenerateHostID 基于主机的MAC地址生成一个稳定的、唯一的标识符。
func GenerateHostID() (string, error) {
	// 1. 获取所有有效的MAC地址
	macAddresses, err := getMACAddresses()
	if err != nil {
		// 如果无法获取MAC地址，则返回错误
		return "", fmt.Errorf("无法生成主机ID: %w", err)
	}

	// 2. 将所有MAC地址排序，确保每次执行的顺序一致
	sort.Strings(macAddresses)

	// 3. 将排序后的MAC地址用一个固定的分隔符连接成一个字符串
	combinedInfo := strings.Join(macAddresses, "|")
	log.Printf("Generating HostID based on MACs: %s", combinedInfo)
	// 4. 使用MD5哈希算法生成最终的ID
	hash := md5.Sum([]byte(combinedInfo))
	return fmt.Sprintf("%x", hash), nil
}

// getMACAddresses 获取本机所有物理网卡的MAC地址列表。
// 它使用一个“白名单”来识别物理网卡，确保稳定性。
func getMACAddresses() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// 定义一个物理网卡名称前缀的“白名单”
	// 'en' 代表以太网和 Wi-Fi，'anpi' 代表 Apple Thunderbolt/USB4 接口
	whitelistPrefixes := []string{"en", "anpi"}

	var macs []string
	for _, iface := range interfaces {
		// 检查当前接口名称是否以白名单中的任何一个前缀开头
		isWhitelisted := false
		for _, prefix := range whitelistPrefixes {
			if strings.HasPrefix(iface.Name, prefix) {
				isWhitelisted = true
				break
			}
		}
		// 如果接口不在白名单中，或者没有硬件地址，则跳过
		if !isWhitelisted || iface.HardwareAddr == nil {
			continue
		}
		mac := iface.HardwareAddr.String()
		if mac != "" {
			macs = append(macs, mac)
		}
	}

	if len(macs) == 0 {
		return nil, fmt.Errorf("未找到任何符合白名单的物理网卡MAC地址")
	}

	return macs, nil
}
