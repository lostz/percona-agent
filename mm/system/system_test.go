/*
   Copyright (c) 2014, Percona LLC and/or its affiliates. All rights reserved.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package system_test

import (
	"github.com/percona/cloud-protocol/proto"
	"github.com/percona/cloud-tools/mm"
	"github.com/percona/cloud-tools/mm/system"
	"github.com/percona/cloud-tools/pct"
	"github.com/percona/cloud-tools/test"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"path/filepath"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

var sample = test.RootDir + "/mm"

/////////////////////////////////////////////////////////////////////////////
// ProcStat
/////////////////////////////////////////////////////////////////////////////

type ProcStatTestSuite struct {
	logChan chan *proto.LogEntry
	logger  *pct.Logger
}

var _ = Suite(&ProcStatTestSuite{})

func (s *ProcStatTestSuite) SetUpSuite(t *C) {
	s.logChan = make(chan *proto.LogEntry, 10)
	s.logger = pct.NewLogger(s.logChan, "system-monitor-test")
}

// --------------------------------------------------------------------------

func (s *ProcStatTestSuite) TestProcStat001(t *C) {
	files, err := filepath.Glob(sample + "/proc/stat001-*.txt")
	if err != nil {
		t.Fatal(err)
	}

	m := system.NewMonitor(s.logger)

	metrics := []mm.Metric{}

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		got, err := m.ProcStat(content)
		if err != nil {
			t.Fatal(err)
		}
		metrics = append(metrics, got...)
	}

	/*
					Totals		Diff
		stat001-1
			cpu		390817611
			cpu0	97641434
			cpu1	97717127
		stat001-2
			cpu		391386603	568992
			cpu0	97783608	142174	These don't add up because the real input has 4 CPU.
			cpu1	97859411	142284  This does not affect the tests.
		stat001-3
			cpu		391759882	373279
			cpu0	97876875	93267
			cpu1	97952757	93346

			1    2    3      4    5      6   7       8     9     10
			user nice system idle iowait irq softirq steal guest guestlow
	*/
	expect := []mm.Metric{
		// First input, no CPU because that requires previous values, so only cpu-ext:
		{Name: "cpu-ext/intr", Type: "counter", Number: 39222211},    // ok
		{Name: "cpu-ext/ctxt", Type: "counter", Number: 122462971},   // ok
		{Name: "cpu-ext/processes", Type: "counter", Number: 227223}, // ok
		{Name: "cpu-ext/procs_running", Type: "gauge", Number: 1},    // ok
		{Name: "cpu-ext/procs_blocked", Type: "gauge", Number: 0},    // ok
		// Second input, now we have CPU values, plus more cpu-ext values.
		{Name: "cpu/user", Type: "gauge", Number: 0.041477},    // ok 5
		{Name: "cpu/nice", Type: "gauge", Number: 0},           // ok
		{Name: "cpu/system", Type: "gauge", Number: 0.017751},  // ok 7
		{Name: "cpu/idle", Type: "gauge", Number: 99.938488},   // ok
		{Name: "cpu/iowait", Type: "gauge", Number: 0.002285},  // ok 9
		{Name: "cpu/irq", Type: "gauge", Number: 0},            // ok
		{Name: "cpu/softirq", Type: "gauge", Number: 0},        // ok 11
		{Name: "cpu/steal", Type: "gauge", Number: 0},          // ok
		{Name: "cpu/guest", Type: "gauge", Number: 0},          // ok 13
		{Name: "cpu0/user", Type: "gauge", Number: 0.131529},   // ok
		{Name: "cpu0/nice", Type: "gauge", Number: 0},          // ok 15
		{Name: "cpu0/system", Type: "gauge", Number: 0.039388}, // ok
		{Name: "cpu0/idle", Type: "gauge", Number: 99.819939},  // ok 17
		{Name: "cpu0/iowait", Type: "gauge", Number: 0.009144},
		{Name: "cpu0/irq", Type: "gauge", Number: 0}, // 19
		{Name: "cpu0/softirq", Type: "gauge", Number: 0},
		{Name: "cpu0/steal", Type: "gauge", Number: 0}, // 21
		{Name: "cpu0/guest", Type: "gauge", Number: 0},
		{Name: "cpu1/user", Type: "gauge", Number: 0.026707}, // 23
		{Name: "cpu1/nice", Type: "gauge", Number: 0},
		{Name: "cpu1/system", Type: "gauge", Number: 0.023193}, // 25
		{Name: "cpu1/idle", Type: "gauge", Number: 99.950100},
		{Name: "cpu1/iowait", Type: "gauge", Number: 0}, // 27
		{Name: "cpu1/irq", Type: "gauge", Number: 0},
		{Name: "cpu1/softirq", Type: "gauge", Number: 0}, // 29
		{Name: "cpu1/steal", Type: "gauge", Number: 0},
		{Name: "cpu1/guest", Type: "gauge", Number: 0},               // ok 31
		{Name: "cpu-ext/intr", Type: "counter", Number: 39276666},    // ok
		{Name: "cpu-ext/ctxt", Type: "counter", Number: 122631533},   // ok 33
		{Name: "cpu-ext/processes", Type: "counter", Number: 227521}, // ok
		{Name: "cpu-ext/procs_running", Type: "gauge", Number: 2},    // ok 35
		{Name: "cpu-ext/procs_blocked", Type: "gauge", Number: 0},    // ok
		// Third input.
		{Name: "cpu/user", Type: "gauge", Number: 0.038309},   // ok 37
		{Name: "cpu/nice", Type: "gauge", Number: 0},          // ok
		{Name: "cpu/system", Type: "gauge", Number: 0.017681}, // ok 39
		{Name: "cpu/idle", Type: "gauge", Number: 99.941063},
		{Name: "cpu/iowait", Type: "gauge", Number: 0.002947}, // 41
		{Name: "cpu/irq", Type: "gauge", Number: 0},
		{Name: "cpu/softirq", Type: "gauge", Number: 0}, // 43
		{Name: "cpu/steal", Type: "gauge", Number: 0},
		{Name: "cpu/guest", Type: "gauge", Number: 0}, // 45
		{Name: "cpu0/user", Type: "gauge", Number: 0.122230},
		{Name: "cpu0/nice", Type: "gauge", Number: 0}, // 47
		{Name: "cpu0/system", Type: "gauge", Number: 0.041815},
		{Name: "cpu0/idle", Type: "gauge", Number: 99.824161}, // 49
		{Name: "cpu0/iowait", Type: "gauge", Number: 0.011794},
		{Name: "cpu0/irq", Type: "gauge", Number: 0}, // 51
		{Name: "cpu0/softirq", Type: "gauge", Number: 0},
		{Name: "cpu0/steal", Type: "gauge", Number: 0}, // 53
		{Name: "cpu0/guest", Type: "gauge", Number: 0},
		{Name: "cpu1/user", Type: "gauge", Number: 0.021426},   // ok 55
		{Name: "cpu1/nice", Type: "gauge", Number: 0},          // ok
		{Name: "cpu1/system", Type: "gauge", Number: 0.024640}, // -- 57
		{Name: "cpu1/idle", Type: "gauge", Number: 99.953935},
		{Name: "cpu1/iowait", Type: "gauge", Number: 0}, // 59
		{Name: "cpu1/irq", Type: "gauge", Number: 0},
		{Name: "cpu1/softirq", Type: "gauge", Number: 0}, // 61
		{Name: "cpu1/steal", Type: "gauge", Number: 0},
		{Name: "cpu1/guest", Type: "gauge", Number: 0},               // 63
		{Name: "cpu-ext/intr", Type: "counter", Number: 39312673},    // ok
		{Name: "cpu-ext/ctxt", Type: "counter", Number: 122742465},   // 65 ok
		{Name: "cpu-ext/processes", Type: "counter", Number: 227717}, // ok
		{Name: "cpu-ext/procs_running", Type: "gauge", Number: 3},    // 67 ok
		{Name: "cpu-ext/procs_blocked", Type: "gauge", Number: 0},    // ok
	}

	if same, diff := test.IsDeeply(metrics, expect); !same {
		test.Dump(metrics)
		t.Error(diff)
	}
}

/////////////////////////////////////////////////////////////////////////////
// ProcMeminfo
/////////////////////////////////////////////////////////////////////////////

type ProcMeminfoTestSuite struct {
	logChan chan *proto.LogEntry
	logger  *pct.Logger
}

var _ = Suite(&ProcMeminfoTestSuite{})

func (s *ProcMeminfoTestSuite) SetUpSuite(t *C) {
	s.logChan = make(chan *proto.LogEntry, 10)
	s.logger = pct.NewLogger(s.logChan, "system-monitor-test")
}

// --------------------------------------------------------------------------

func (s *ProcMeminfoTestSuite) TestProcMeminfo001(t *C) {
	m := system.NewMonitor(s.logger)
	content, err := ioutil.ReadFile(sample + "/proc/meminfo001.txt")
	if err != nil {
		t.Fatal(err)
	}
	got, err := m.ProcMeminfo(content)
	if err != nil {
		t.Fatal(err)
	}
	// Remember: the order of this array must match order in which each
	// stat appears in the input file:
	expect := []mm.Metric{
		{Name: "memory/MemTotal", Type: "gauge", Number: 8046892},  // ok
		{Name: "memory/MemFree", Type: "gauge", Number: 5273644},   // ok
		{Name: "memory/Buffers", Type: "gauge", Number: 300684},    // ok
		{Name: "memory/Cached", Type: "gauge", Number: 946852},     // ok
		{Name: "memory/SwapCached", Type: "gauge", Number: 0},      // ok
		{Name: "memory/Active", Type: "gauge", Number: 1936436},    // ok
		{Name: "memory/Inactive", Type: "gauge", Number: 598916},   // ok
		{Name: "memory/SwapTotal", Type: "gauge", Number: 8253436}, // ok
		{Name: "memory/SwapFree", Type: "gauge", Number: 8253436},  // ok
		{Name: "memory/Dirty", Type: "gauge", Number: 0},           // ok
	}
	if same, diff := test.IsDeeply(got, expect); !same {
		test.Dump(got)
		t.Error(diff)
	}
}

/////////////////////////////////////////////////////////////////////////////
// ProcVmstat
/////////////////////////////////////////////////////////////////////////////

type ProcVmstatTestSuite struct {
	logChan chan *proto.LogEntry
	logger  *pct.Logger
}

var _ = Suite(&ProcVmstatTestSuite{})

func (s *ProcVmstatTestSuite) SetUpSuite(t *C) {
	s.logChan = make(chan *proto.LogEntry, 10)
	s.logger = pct.NewLogger(s.logChan, "system-monitor-test")
}

// --------------------------------------------------------------------------

func (s *ProcVmstatTestSuite) TestProcVmstat001(t *C) {
	m := system.NewMonitor(s.logger)
	content, err := ioutil.ReadFile(sample + "/proc/vmstat001.txt")
	if err != nil {
		t.Fatal(err)
	}
	got, err := m.ProcVmstat(content)
	if err != nil {
		t.Fatal(err)
	}
	// Remember: the order of this array must match order in which each
	// stat appears in the input file:
	expect := []mm.Metric{
		{Name: "vmstat/numa_hit", Type: "counter", Number: 42594095},    // ok
		{Name: "vmstat/numa_miss", Type: "counter", Number: 0},          // ok
		{Name: "vmstat/numa_foreign", Type: "counter", Number: 0},       // ok
		{Name: "vmstat/numa_interleave", Type: "counter", Number: 7297}, // ok
		{Name: "vmstat/numa_local", Type: "counter", Number: 42594095},  // ok
		{Name: "vmstat/numa_other", Type: "counter", Number: 0},         // ok
		{Name: "vmstat/pgpgin", Type: "counter", Number: 646645},        // ok
		{Name: "vmstat/pgpgout", Type: "counter", Number: 5401659},      // ok
		{Name: "vmstat/pswpin", Type: "counter", Number: 0},             // ok
		{Name: "vmstat/pswpout", Type: "counter", Number: 0},            // ok
	}
	if same, diff := test.IsDeeply(got, expect); !same {
		test.Dump(got)
		t.Error(diff)
	}
}
