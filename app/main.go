// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A simple example exposing fictional RPC latencies with different types of
// random distributions (uniform, normal, and exponential) as Prometheus
// metrics.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	httpDurationsHistogram *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer, factor float64) *metrics {
	m := &metrics{
		httpDurationsHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:                         "http_request_durations",
			Help:                         "HTTP latency distributions.",
			Buckets:                      prometheus.DefBuckets,
			NativeHistogramZeroThreshold: 0.05,
			NativeHistogramBucketFactor:  factor,
		}, []string{"service", "code"}),
	}
	reg.MustRegister(m.httpDurationsHistogram)
	return m
}

func main() {
	var (
		addr        = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
		factor      = flag.Float64("native-factor", 1.1, "The factor of histogram bucket to setting.")
		metricCount = flag.Int("metrics-count", 10, "The count of metrics.")
	)
	flag.Parse()

	reg := prometheus.NewRegistry()
	m := NewMetrics(reg, *factor)

	go func() {
		for {
			for i := 1; i <= *metricCount; i++ {
				v := rand.Float64() * 12
				observer := m.httpDurationsHistogram.WithLabelValues("service#"+strconv.Itoa(i), "200")

				if v > 10 {
					observer.(prometheus.ExemplarObserver).ObserveWithExemplar(
						v,
						prometheus.Labels{"traceID": fmt.Sprint(rand.Intn(100000))},
					)
				} else {
					observer.Observe(v)
				}
			}

			time.Sleep(10 * time.Millisecond)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
			Registry:          reg,
		},
	))

	log.Fatal(http.ListenAndServe(*addr, nil))
}
