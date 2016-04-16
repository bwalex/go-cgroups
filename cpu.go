package cgroups

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type CpuStat struct {
	// From cpuacct.stat
	UserTimeUs       uint64 /* in microseconds */
	SystemTimeUs     uint64 /* in microseconds */

	// From cpu.stat
	ThrottledTimeUs  uint64 /* in microseconds */
	Periods          uint64
	ThrottledPeriods uint64

	// Derived
	ThrottledPct     float64
}

const (
	ControllerCpu = "cpu"
)

var ticksPerSec = uint64(sysconfClockTicks())

func populateCpuStat(cg Cgroup, stat *CpuStat) error {
	path, err := GetCgroupPath(cg, ControllerCpu, "cpu.stat")
	if err == ErrNoCgroup {
		return err
	}

	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		value, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}

		switch parts[0] {
		case "nr_periods":
			stat.Periods = value
		case "nr_throttled":
			stat.ThrottledPeriods = value
		case "throttled_time":
			stat.ThrottledTimeUs = ticksToUs(value)
		}
	}

	if stat.Periods == 0 {
		stat.ThrottledPct = 0
	} else {
		stat.ThrottledPct = 100.0 * float64(stat.ThrottledPeriods) / float64(stat.Periods)
	}

	return nil
}

func populateCpuacctStat(cg Cgroup, stat *CpuStat) error {
	path, err := GetCgroupPath(cg, ControllerCpu, "cpuacct.stat")
	if err == ErrNoCgroup {
		return err
	}

	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		value, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}

		switch parts[0] {
		case "user":
			stat.UserTimeUs = ticksToUs(value)
		case "system":
			stat.SystemTimeUs = ticksToUs(value)
		}
	}
	return nil
}

func GetCpuStats(cg Cgroup) (CpuStat, error) {
	var stats CpuStat

	err := populateCpuStat(cg, &stats)
	if err != nil {
		return stats, err
	}

	err = populateCpuacctStat(cg, &stats)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

func ticksToUs(ticks uint64) uint64 {
	return ticks * 1000 * 1000 / ticksPerSec
}
