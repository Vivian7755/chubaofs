// Copyright 2018 The Container File System Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package exporter

import (
	"fmt"
	"testing"
)

func TestNewMetricLabels(t *testing.T) {
	N := 100
	exitCh := make(chan int, N)
	for i := 0; i < N; i++ {
		go func(i int) {
			name := fmt.Sprintf("name_%d_gauge", i%6)
			label := fmt.Sprintf("label-%d:name", i)
			m := RegistMetricWithLabels(name, Gauge, map[string]string{"volname": label, "cluster": name})
			if m != nil {
				m.Set(float64(i))
				t.Logf("metric: %v, %v", name, m.Key)
			}
			name2 := fmt.Sprintf("name_%d_counter", i%6)
			c := RegistMetricWithLabels(name2, Counter, map[string]string{"volname": label, "cluster": name})
			if c != nil {
				c.Set(float64(i))
				t.Logf("metric: %v, %v", name2, c.Key)
			}
			exitCh <- i

		}(i)
	}

	x := 0
	select {
	case <-exitCh:
		x += 1
		if x == N {
			return
		}
	}
}

func TestNewMetric(t *testing.T) {
	N := 100
	exitCh := make(chan int, N)
	for i := 0; i < N; i++ {
		go func() {
			m := RegistCounter(fmt.Sprintf("name_%d_metric_count", i%17))
			if m != nil {
				m.Add(float64(i))
				t.Logf("metric: %v", m.Desc())
			}
			g := RegistGauge(fmt.Sprintf("name_%d_metric", i%17))
			if g != nil {
				g.Set(float64(i))
				t.Logf("metric: %v", g.Desc())
			}
			exitCh <- i

		}()
	}

	x := 0
	select {
	case <-exitCh:
		x += 1
		if x == N {
			return
		}
	}
}

func TestRegistGauge(t *testing.T) {
	N := 100
	exitCh := make(chan int, 100)
	for i := 0; i < N; i++ {
		go func() {
			m := RegistGauge(fmt.Sprintf("name_%d", i%7))
			if m != nil {
				t.Logf("metric: %v", m.Desc().String())
			}
			exitCh <- i

		}()
	}

	x := 0
	select {
	case <-exitCh:
		x += 1
		if x == N {
			return
		}
	}
}

func TestRegistTp(t *testing.T) {
	N := 100
	exitCh := make(chan int, 100)
	for i := 0; i < N; i++ {
		go func() {
			m := RegistTp(fmt.Sprintf("name_%d", i%7))
			if m != nil {
				t.Logf("metric: %v", m.Name)
			}

			defer m.CalcTp()

			exitCh <- i
		}()
	}

	x := 0
	select {
	case <-exitCh:
		x += 1
		if x == N {
			return
		}
	}

}
