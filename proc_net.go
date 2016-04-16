package cgroups

import (
	"bufio"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

type NetStat struct {
	RxBytes      uint64 `field:"0"`
	RxPackets    uint64 `field:"1"`
	RxErrors     uint64 `field:"2"`
	RxDrop       uint64 `field:"3"`
	RxFifo       uint64 `field:"4"`
	RxFrame      uint64 `field:"5"`
	RxCompressed uint64 `field:"6"`
	RxMulticast  uint64 `field:"7"`

	TxBytes      uint64 `field:"8"`
	TxPackets    uint64 `field:"9"`
	TxErrors     uint64 `field:"10"`
	TxDrop       uint64 `field:"11"`
	TxFifo       uint64 `field:"12"`
	TxFrame      uint64 `field:"13"`
	TxCompressed uint64 `field:"14"`
	TxMulticast  uint64 `field:"15"`
}

func GetNetInterfaces(cg Cgroup) ([]string, error) {
	nets := make([]string, 0, 4)
	lines, err := procNetDevLines(cg)
	if err != nil {
		return nets, err
	}

	for i := range lines {
		fields := strings.Fields(lines[i])
		if !strings.HasSuffix(fields[0], ":") {
			continue
		}

		nets = append(nets, fields[0][:len(fields[0])-1])
	}

	return nets, nil
}

func GetNetStats(cg Cgroup, intf string) (NetStat, error) {
	var stats NetStat

	err := populateNetStats(cg, intf, &stats)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

func populateNetStats(cg Cgroup, intf string, stat *NetStat) error {
	counts := make([]uint64, 0, 20)
	lines, err := procNetDevLines(cg)
	if err != nil {
		return err
	}

	for i := range lines {
		fields := strings.Fields(lines[i])
		if !strings.HasSuffix(fields[0], ":") {
			continue
		}

		net := fields[0][:len(fields[0])-1]

		if intf != "" && net != intf {
			continue
		}

		if intf == "" && net == "lo" {
			continue
		}

		for idx := 1; idx < len(fields); idx++ {
			value, err := strconv.ParseUint(fields[idx], 10, 64)
			if err != nil {
				continue
			}
			if len(counts) < idx {
				counts = append(counts, value)
			} else {
				counts[idx-1] += value
			}
		}
	}

	v := reflect.ValueOf(stat).Elem()
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag

		statFieldRaw := tag.Get("field")
		if statFieldRaw == "" {
			continue
		}

		statField, err := strconv.ParseInt(statFieldRaw, 10, 32)
		if err != nil {
			return err
		}

		if len(counts) <= int(statField) {
			continue
		}

		v.Field(i).SetUint(counts[statField])
	}

	return nil
}

func procNetDevLines(cg Cgroup) ([]string, error) {
	lines := make([]string, 0, 4)

	pids, err := GetProcs(cg)
	if err != nil || len(pids) < 1 {
		return lines, err
	}

	pid := strconv.Itoa(pids[0])

	fd, err := os.Open(path.Join("/proc", pid, "net/dev"))
	if err != nil {
		return lines, err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())
		lines = append(lines, str)
	}

	return lines, nil
}
