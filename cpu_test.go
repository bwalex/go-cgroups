package cgroups

import (
	"testing"
)

func TestCpuStat(t *testing.T) {
	stats, err := GetCpuStats(Cgroup{ Cgroup: "/system.slice" })

	if err != nil {
		t.Fail()
	}

	if stats.UserTimeUs < 1 {
		t.Fail()
	}

	if stats.SystemTimeUs < 1 {
		t.Fail()
	}

	if stats.ThrottledPct < 0 {
		t.Fail()
	}

	t.Logf("%+v\n", stats)
}

func BenchmarkCpuStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetCpuStats(Cgroup{ Cgroup: "/system.slice" })
		if err != nil {
			b.Fail()
		}
	}
}
