package cgroups

import (
	"bufio"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type BlkioStat struct {
	// From other files
	Merged            uint64 `file:"blkio.io_merged_recursive"`
	MergedRead        uint64 `file:"blkio.io_merged_recursive" sum:"Read"`
	MergedWrite       uint64 `file:"blkio.io_merged_recursive" sum:"Write"`
	Queued            uint64 `file:"blkio.io_queued_recursive"`
	QueuedRead        uint64 `file:"blkio.io_queued_recursive" sum:"Read"`
	QueuedWrite       uint64 `file:"blkio.io_queued_recursive" sum:"Write"`
	ServiceBytes      uint64 `file:"blkio.io_service_bytes_recursive"`
	ServiceBytesRead  uint64 `file:"blkio.io_service_bytes_recursive" sum:"Read"`
	ServiceBytesWrite uint64 `file:"blkio.io_service_bytes_recursive" sum:"Write"`
	Serviced          uint64 `file:"blkio.io_serviced_recursive"`
	ServicedRead      uint64 `file:"blkio.io_serviced_recursive" sum:"Read"`
	ServicedWrite     uint64 `file:"blkio.io_serviced_recursive" sum:"Write"`
	ServiceTime       uint64 `file:"blkio.io_service_time_recursive"`
	ServiceTimeRead   uint64 `file:"blkio.io_service_time_recursive" sum:"Read"`
	ServiceTimeWrite  uint64 `file:"blkio.io_service_time_recursive" sum:"Write"`
	WaitTime          uint64 `file:"blkio.io_wait_time_recursive"`
	WaitTimeRead      uint64 `file:"blkio.io_wait_time_recursive" sum:"Read"`
	WaitTimeWrite     uint64 `file:"blkio.io_wait_time_recursive" sum:"Write"`
}

const (
	ControllerBlkio = "blkio"
)

func blkioParse(filePath string, sumField string) (uint64, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	var value uint64

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")

		if sumField == "" {
			if strings.ToLower(parts[0]) == "total" {
				v, err := strconv.ParseUint(parts[1], 10, 64)
				return v, err
			}
		} else {
			if strings.ToLower(parts[0]) == "total" {
				continue
			}

			accessType := parts[1]

			if strings.ToLower(accessType) != strings.ToLower(sumField) {
				continue
			}

			v, err := strconv.ParseUint(parts[2], 10, 64)
			if err != nil {
				continue
			}

			value += v
		}
	}

	return value, nil
}

func populateBlkioOther(cg Cgroup, stat *BlkioStat) error {
	v := reflect.ValueOf(stat).Elem()
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag

		fileName := tag.Get("file")
		if fileName == "" {
			continue
		}

		path, err := GetCgroupPath(cg, ControllerBlkio, fileName)
		if err == ErrNoCgroup {
			return err
		}

		value, err := blkioParse(path, tag.Get("sum"))
		if err != nil {
			continue
		}

		v.Field(i).SetUint(value)
	}

	return nil
}

func GetBlkioStats(cg Cgroup) (BlkioStat, error) {
	var stats BlkioStat

	err := populateBlkioOther(cg, &stats)
	if err != nil {
		return stats, err
	}

	return stats, nil
}
