package cgroups

import (
	"testing"
)

func TestBlkioStat(t *testing.T) {
	stats, err := GetBlkioStats(Cgroup{ Cgroup: "/system.slice" })

	if err != nil {
		t.Fail()
	}

	if stats.Serviced < stats.ServicedRead+stats.ServicedWrite {
		t.Fail()
	}

	if stats.ServiceBytes < stats.ServiceBytesRead+stats.ServiceBytesWrite {
		t.Fail()
	}

	if stats.ServiceBytes <= 0 {
		t.Fail()
	}

	if stats.WaitTime <= 0 {
		t.Fail()
	}

	if stats.ServiceTime <= 0 {
		t.Fail()
	}

	if stats.ServiceTime < stats.ServiceTimeRead+stats.ServiceTimeWrite {
		t.Fail()
	}

	t.Logf("%+v\n", stats)
}

func BenchmarkBlkioStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetBlkioStats(Cgroup{ Cgroup: "/system.slice" })
		if err != nil {
			b.Fail()
		}
	}
}
