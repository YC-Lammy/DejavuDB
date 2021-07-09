package main

import (
	"encoding/json"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
)

var monitor_values = []map[string]string{}

func monitor() map[string]string {
	result := map[string]string{}

	result["role"] = role

	v, _ := mem.VirtualMemory()
	result["mem_load"] = strconv.FormatInt(int64(math.Round(v.UsedPercent)), 10)
	result["mem_total"] = strconv.FormatInt(int64(v.Total/1000000000), 10)
	result["mem_used"] = strconv.FormatInt(int64(v.Used/1000000000), 10)
	result["mem_avaliable"] = strconv.FormatInt(int64((v.Total-v.Used)/1000000000), 10)

	c, _ := cpu.Percent(time.Second, true)
	cpup := 0.0
	for _, v := range c {
		cpup += v
	}

	result["cpu_load"] = strconv.FormatInt(int64(math.Round(cpup)), 10)

	d, _ := disk.Partitions(true)
	disks := []string{}
	disk_total := 0
	disk_used := 0
	disk_avaliable := 0
	for _, p := range d {
		disks = append(disks, p.Mountpoint)
		di, _ := disk.Usage(p.Mountpoint)
		disk_total += int(di.Total)               //1073741824 // GB as unit
		disk_used += int(di.Used)                 //1073741824
		disk_avaliable += int(di.Total - di.Used) //1073741824
	}
	result["disk_load"] = strconv.FormatInt(int64((float64(disk_used)/float64(disk_total))*100.0), 10)

	result["disk_total"] = strconv.FormatInt(int64(disk_total/1000000000), 10)

	result["disk_used"] = strconv.FormatInt(int64(disk_used/1000000000), 10)

	result["disk_avaliable"] = strconv.FormatInt(int64(disk_avaliable/1000000000), 10)

	result["disk_mountpoint"] = strings.Join(disks, ",")

	result["tasks"] = strconv.FormatInt(int64(len(process_query)), 10)

	return result

}

func register_monitor_value(jsonstr string) {
	result := map[string]string{}
	json.Unmarshal([]byte(jsonstr), &result)
	monitor_values = append(monitor_values, result)
}

func get_all_monitor(conn net.Conn) {
	send_to_all_router([]byte("monitor"))
	send_to_all_shard([]byte("monitor"))
}
