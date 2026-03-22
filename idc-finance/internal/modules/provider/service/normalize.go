package service

import (
	"strconv"
	"strings"

	orderDomain "idc-finance/internal/modules/order/domain"
)

func buildConfigurationSelections(remote map[string]any) []orderDomain.ServiceConfigSelection {
	items := make([]orderDomain.ServiceConfigSelection, 0, 6)
	appendItem := func(code, name, value, label string) {
		if strings.TrimSpace(value) == "" {
			return
		}
		if strings.TrimSpace(label) == "" {
			label = value
		}
		items = append(items, orderDomain.ServiceConfigSelection{
			Code:       code,
			Name:       name,
			Value:      value,
			ValueLabel: label,
			PriceDelta: 0,
		})
	}

	appendItem("cpu", "CPU 核数", pickString(remote, "cpu", "cpucores"), pickString(remote, "cpu", "cpucores"))
	appendItem("memory", "内存规格", pickString(remote, "memory"), pickString(remote, "memory"))
	appendItem("system", "操作系统", pickString(remote, "operate_system", "os_name", "os"), pickString(remote, "operate_system", "os_name", "os"))
	appendItem("node", "节点", pickString(remote, "node_name"), pickString(remote, "node_name"))
	appendItem("bandwidth", "带宽规格", firstNonEmptyString(
		pickString(remote, "out_bw"),
		pickString(recordMap(remote["default_bw_group"]), "out_bw"),
	), firstNonEmptyString(
		pickString(remote, "out_bw"),
		pickString(recordMap(remote["default_bw_group"]), "out_bw"),
	))
	appendItem("traffic", "流量限额", pickString(remote, "traffic_quota"), pickString(remote, "traffic_quota"))
	return items
}

func buildResourceSnapshot(remote map[string]any) orderDomain.ServiceResourceSnapshot {
	systemDisk := 0
	dataDisk := 0
	for _, disk := range recordArray(remote["disk"]) {
		size := pickInt(disk, "size", "size_gb")
		if firstNonEmptyString(pickString(disk, "type", "disk_type"), "") == "system" {
			systemDisk += size
		} else {
			dataDisk += size
		}
	}

	if systemDisk == 0 {
		systemDisk = pickInt(remote, "system_disk_size")
	}
	if dataDisk == 0 {
		dataDisk = pickInt(remote, "data_disk_size")
	}

	memoryGB := pickInt(remote, "memory")
	cpuCores := firstNonZeroInt(pickInt(remote, "cpu"), pickInt(remote, "cpucores"))
	bandwidth := firstNonZeroInt(
		pickInt(remote, "out_bw"),
		pickInt(recordMap(remote["default_bw_group"]), "out_bw"),
		pickInt(remote, "in_bw"),
	)

	ipv6 := ""
	if ipv6Items := recordArray(remote["ipv6"]); len(ipv6Items) > 0 {
		ipv6 = pickString(ipv6Items[0], "ipaddress", "ip", "address")
	}

	return orderDomain.ServiceResourceSnapshot{
		RegionName:      firstNonEmptyString(pickString(remote, "area_name"), pickRegion(remote)),
		ZoneName:        firstNonEmptyString(pickString(remote, "node_name"), ""),
		Hostname:        firstNonEmptyString(pickString(remote, "hostname"), "mofang-instance"),
		OperatingSystem: firstNonEmptyString(pickString(remote, "operate_system", "os_name", "os"), "-"),
		LoginUsername:   firstNonEmptyString(pickString(remote, "osuser"), "root"),
		PasswordHint:    "敏感密码不落本地，请通过重置密码或控制台获取",
		SecurityGroup:   firstNonEmptyString(pickString(remote, "security_name"), "-"),
		CPUCores:        cpuCores,
		MemoryGB:        memoryGB,
		SystemDiskGB:    systemDisk,
		DataDiskGB:      dataDisk,
		BandwidthMbps:   bandwidth,
		PublicIPv4:      firstNonEmptyString(pickString(remote, "mainip", "main_ip"), firstIPFromRemote(remote)),
		PublicIPv6:      ipv6,
	}
}

func firstIPFromRemote(remote map[string]any) string {
	for _, item := range recordArray(remote["ip"]) {
		if value := pickString(item, "ipaddress", "ip", "address"); value != "" {
			return value
		}
	}
	for _, networkItem := range recordArray(remote["network"]) {
		for _, ipItem := range recordArray(networkItem["ipaddress"]) {
			if value := pickString(ipItem, "ipaddress", "ip", "address"); value != "" {
				return value
			}
		}
	}
	return ""
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func firstNonZeroInt(values ...int) int {
	for _, value := range values {
		if value != 0 {
			return value
		}
	}
	return 0
}

func parseIntString(value string) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0
	}
	return parsed
}
