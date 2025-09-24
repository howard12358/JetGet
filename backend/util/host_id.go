package util

import (
	"crypto/md5"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
)

// GenerateHostID 基于主机信息生成稳定的标识符
func GenerateHostID() (string, error) {
	// 收集主机信息
	var hostInfo []string

	// 获取主机名
	hostname, err := os.Hostname()
	if err == nil {
		hostInfo = append(hostInfo, hostname)
	}

	// 获取MAC地址
	macAddr, err := getMACAddress()
	if err == nil {
		hostInfo = append(hostInfo, macAddr)
	}

	// 如果没有获取到足够的信息，返回错误
	if len(hostInfo) == 0 {
		return "", fmt.Errorf("无法获取主机信息")
	}

	// 将主机信息排序以确保一致性
	sort.Strings(hostInfo)

	// 使用MD5哈希生成稳定的标识符
	hash := md5.Sum([]byte(strings.Join(hostInfo, "|")))
	return fmt.Sprintf("%x", hash), nil
}

// getMACAddress 获取主MAC地址
func getMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// 查找第一个非回环且有硬件地址的接口
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback == 0 && iface.HardwareAddr != nil {
			mac := iface.HardwareAddr.String()
			if mac != "" {
				return mac, nil
			}
		}
	}

	return "", fmt.Errorf("未找到有效的MAC地址")
}
