package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
)

func monitor() map[string]string {
	result := map[string]string{}
	v, _ := mem.VirtualMemory()
	result["mem_load"] = strconv.FormatInt(int64(math.Round(v.UsedPercent)), 10)
	result["mem_size"] = strconv.FormatInt(int64(v.Total), 10)

	c, _ := cpu.Percent(time.Second, true)
	cpu := 0.0
	for _, v := range c {
		cpu += v
	}

	result["cpu_load"] = strconv.FormatInt(int64(math.Round(cpu)), 10)

	d, _ := disk.Partitions(false)
	disks := []string{}
	disk_percent := 0.0
	for _, p := range d {
		disks = append(disks, p.Mountpoint)
		di, _ := disk.Usage(p.Mountpoint)
		disk_percent += di.UsedPercent
	}
	disk_percent = disk_percent / float64(len(disks))
	result["disk_load"] = strconv.FormatInt(int64(math.Round(disk_percent)), 10)

	result["mountpoint"] = strings.Join(disks, ",")

	return result

}

func monitor_loop(duration time.Duration) {
	for {
		time.Sleep(duration)
		result := monitor()
		fmt.Println(result)
	}
}
