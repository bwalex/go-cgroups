package cgroups

import (
	"testing"
)

func TestProcNetInterfaces(t *testing.T) {
	devs, err := GetNetInterfaces(Cgroup{ Cgroup: "/system.slice" })

	if err != nil {
		t.Fail()
	}

	if len(devs) < 1 {
		t.Fail()
	}

	t.Logf("%+v\n", devs)
}

func TestNetStatSingle(t *testing.T) {
	stats, err := GetNetStats(Cgroup{ Cgroup: "/system.slice" }, "lo")

	if err != nil {
		t.Fail()
	}

	if stats.RxPackets >= stats.RxBytes {
		t.Fail()
	}

	if stats.TxPackets >= stats.TxBytes {
		t.Fail()
	}

	if stats.RxPackets < 1 {
		t.Fail()
	}

	if stats.TxPackets < 1 {
		t.Fail()
	}

	t.Logf("%+v\n", stats)
}

func TestNetStatTotals(t *testing.T) {
	stats, err := GetNetStats(Cgroup{ Cgroup: "/system.slice" }, "")

	if err != nil {
		t.Fail()
	}

	if stats.RxPackets >= stats.RxBytes {
		t.Fail()
	}

	if stats.TxPackets >= stats.TxBytes {
		t.Fail()
	}

	if stats.RxPackets < 1 {
		t.Fail()
	}

	if stats.TxPackets < 1 {
		t.Fail()
	}

	t.Logf("%+v\n", stats)
}

func BenchmarkNetStatTotals(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetNetStats(Cgroup{ Cgroup: "/system.slice" }, "")
		if err != nil {
			b.Fail()
		}
	}
}
