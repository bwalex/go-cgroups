package cgroups

import (
	"testing"
)

func TestMemoryStat(t *testing.T) {
	stats, err := GetMemoryStats(Cgroup{ Cgroup: "/system.slice" })

	if err != nil {
		t.Fail()
	}

	if stats.RSS < 1024 {
		t.Fail()
	}

	if stats.TotalRSS < stats.RSS {
		t.Fail()
	}

	if stats.PgMajFault < 1 {
		t.Fail()
	}

	if stats.MemUsage > stats.MemUsageMax {
		t.Fail()
	}

	if stats.MemSwapUsage < stats.MemUsage {
		t.Fail()
	}

	if stats.MemUsage < 1024 {
		t.Fail()
	}

	t.Logf("%+v\n", stats)
}

func BenchmarkMemoryStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetMemoryStats(Cgroup{ Cgroup: "/system.slice" })
		if err != nil {
			b.Fail()
		}
	}
}
