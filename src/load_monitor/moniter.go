package monitor

import (
	"math"
	"strings"
	"time"

	json "github.com/goccy/go-json"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"

	"src/settings"
)

var Child_load_register = map[int64]*load{}
var monitor_values = []load{}

type load struct {
	id            int64
	role          string
	mem_load      uint8
	mem_total     uint64
	mem_used      uint64
	mem_avaliable uint64

	cpu_load uint8
	tasks    uint64

	disk_load       uint8
	disk_total      uint64
	disk_used       uint64
	disk_avaliable  uint64
	disk_mountpoint string
}

func monitor() load {
	result := load{}

	result.role = settings.Role

	v, _ := mem.VirtualMemory()
	result.mem_load = uint8(math.Round(v.UsedPercent))
	result.mem_total = v.Total / 1000000000
	result.mem_used = v.Used / 1000000000
	result.mem_avaliable = (v.Total - v.Used) / 1000000000

	c, _ := cpu.Percent(time.Second, false)
	cpup := 0.0
	for _, v := range c {
		cpup += v
	}

	result.cpu_load = uint8(math.Round(cpup))

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
	result.disk_load = uint8((disk_used * 100) / (disk_total * 100))

	result.disk_total = uint64(disk_total / 1000000000)

	result.disk_used = uint64(disk_used / 1000000000)

	result.disk_avaliable = uint64(disk_avaliable / 1000000000)

	result.disk_mountpoint = strings.Join(disks, ",")

	return result

}

func register_monitor_value(jsonstr string) {
	result := load{}
	json.Unmarshal([]byte(jsonstr), &result)
	monitor_values = append(monitor_values, result)
	Child_load_register[result.id] = &result
}
