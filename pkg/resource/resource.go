// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package resource

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/lock"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

type UsedResource struct {
	l                            lock.RWMutex
	ctx                          context.Context
	client                       client.Client
	logger                       *zap.Logger
	mem                          uint64  // byte
	cpu                          float64 // percent
	systemTotalAppHttpHealthyQPS int64
	systemTotalNetReachQPS       int64
	systemTotalNetDnsQPS         int64
	taskName                     string
}

func InitResource(ctx context.Context, logger *zap.Logger, client client.Client, taskName string) *UsedResource {
	return &UsedResource{
		ctx:      ctx,
		logger:   logger,
		client:   client,
		taskName: taskName,
	}
}

func (r *UsedResource) RunResourceCollector() {
	interval := time.Duration(types.AgentConfig.CollectResourceInSecond) * time.Second

	type task struct {
		name string
		kind string
	}

	runningTask := make(map[task]int64, 0)

	go func() {
		for {
			if r.ctx.Err() != nil {
				return
			}

			// get cpu and mem
			cpuStats, err := cpu.Percent(interval, false)
			if err == nil {
				if r.cpu < cpuStats[0] {
					r.cpu = cpuStats[0]
				}
			}
			m := &runtime.MemStats{}
			runtime.ReadMemStats(m)
			if r.mem < m.Sys {
				r.mem = m.Sys
			}

			// get other running task qps
			// all doing AppHttpHealthy qps
			appHealthyList := new(crd.AppHttpHealthyList)
			if err := r.client.List(r.ctx, appHealthyList); err != nil {
				r.logger.Sugar().Errorf("list AppHealthyList failed,err=%+v", err)
			}

			for _, app := range appHealthyList.Items {
				if len(app.Status.History) > 0 && app.Name != r.taskName {
					if app.Status.History[len(app.Status.History)-1].Status == crd.StatusHistoryRecordStatusOngoing {
						runningTask[task{name: app.Name, kind: app.Kind}] = int64(app.Spec.Request.QPS)
					} else {
						// Avoid counting completed tasks and perform deletion operations
						delete(runningTask, task{name: app.Name, kind: app.Kind})
					}
				}
			}

			// all doing netReach qps
			netReachList := new(crd.NetReachList)
			if err := r.client.List(r.ctx, netReachList); err != nil {
				r.logger.Sugar().Errorf("list NetReachList failed,err=%+v", err)
			}

			for _, app := range netReachList.Items {
				if len(app.Status.History) > 0 && app.Name != r.taskName {
					if app.Status.History[len(app.Status.History)-1].Status == crd.StatusHistoryRecordStatusOngoing {
						taskNum := 0

						if *app.Spec.Target.ClusterIP {
							if *app.Spec.Target.IPv4 {
								taskNum += 1
							}
							if *app.Spec.Target.IPv6 {
								taskNum += 1
							}
						}
						if *app.Spec.Target.LoadBalancer {
							if *app.Spec.Target.IPv4 {
								taskNum += 1
							}
							if *app.Spec.Target.IPv6 {
								taskNum += 1
							}
						}
						if *app.Spec.Target.Endpoint {
							if *app.Spec.Target.IPv4 {
								taskNum += 1
							}
							if *app.Spec.Target.IPv6 {
								taskNum += 1
							}
						}
						if *app.Spec.Target.NodePort {
							if *app.Spec.Target.IPv4 {
								taskNum += 1
							}
							if *app.Spec.Target.IPv6 {
								taskNum += 1
							}
						}
						if *app.Spec.Target.Ingress {
							if *app.Spec.Target.IPv4 {
								taskNum += 1
							}
						}

						// TODO [ii2day] test multus count qps

						runningTask[task{name: app.Name, kind: app.Kind}] = int64(app.Spec.Request.QPS * taskNum)
					} else {
						// Avoid counting completed tasks and perform deletion operations
						delete(runningTask, task{name: app.Name, kind: app.Kind})
					}
				}
			}

			// all doing NetDns qps
			netDnsList := new(crd.NetdnsList)
			if err := r.client.List(r.ctx, netDnsList); err != nil {
				r.logger.Sugar().Errorf("list NetdnsList failed,err=%+v", err)
			}

			for _, app := range netDnsList.Items {
				if len(app.Status.History) > 0 && app.Name != r.taskName {
					if app.Status.History[len(app.Status.History)-1].Status == crd.StatusHistoryRecordStatusOngoing {
						runningTask[task{name: app.Name, kind: app.Kind}] = int64(*app.Spec.Request.QPS)
					} else {
						// Avoid counting completed tasks and perform deletion operations
						delete(runningTask, task{name: app.Name, kind: app.Kind})
					}
				}
			}

			for k, v := range runningTask {
				switch k.kind {
				case types.KindNameAppHttpHealthy:
					r.systemTotalAppHttpHealthyQPS += v
				case types.KindNameNetReach:
					r.systemTotalNetReachQPS += v
				case types.KindNameNetdns:
					r.systemTotalNetDnsQPS += v
				}
			}
		}
	}()
}

func (r *UsedResource) Stats() v1beta1.SystemResource {
	r.l.RLock()
	defer r.l.RUnlock()
	resource := v1beta1.SystemResource{
		MaxCPU:                       fmt.Sprintf("%.3f%%", r.cpu),
		MaxMemory:                    fmt.Sprintf("%.2fMB", float64(r.mem/(1024*1024))),
		SystemTotalAppHttpHealthyQPS: r.systemTotalAppHttpHealthyQPS,
		SystemTotalNetReachQPS:       r.systemTotalNetReachQPS,
		SystemTotalNetDnsQPS:         r.systemTotalNetDnsQPS,
	}
	return resource
}

func (r *UsedResource) CleanStats() {
	r.l.Lock()
	defer r.l.Unlock()
	r.mem = 0
	r.cpu = 0
	r.systemTotalAppHttpHealthyQPS = 0
	r.systemTotalNetDnsQPS = 0
	r.systemTotalNetReachQPS = 0
}
