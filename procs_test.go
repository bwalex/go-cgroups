package cgroups

import (
	"testing"
)

func TestProcs(t *testing.T) {
	pids, err := GetProcs(Cgroup{ Cgroup: "/system.slice" })

	if err != nil {
		t.Fail()
	}

	if len(pids) < 1 {
		t.Fail()
	}

	t.Logf("%+v\n", pids)
}


func BenchmarkProcs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetProcs(Cgroup{ Cgroup: "/system.slice" })
		if err != nil {
			b.Fail()
		}
	}
}
