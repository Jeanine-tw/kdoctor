// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package apphttphealth_test

import (
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/kdoctor-io/kdoctor/test/e2e/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spidernet-io/e2eframework/tools"
	"net"
	"time"
)

var _ = Describe("testing appHttpHealth test ", Label("appHttpHealth"), func() {
	var termMin = int64(1)
	var requestTimeout = 3000
	var replicas = int32(1)
	It("deployment threshold", Label("appHealth threshold"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"

		appHttpHealth := new(v1beta1.AppHttpHealthy)

		// agent
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		agentSpec.DeploymentReplicas = &replicas
		agentSpec.Kind = types.KindDeployment

		appHttpHealth.Spec.AgentSpec = agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d", testSvcIP, httpPort)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d", testSvcIP, httpPort)
		}
		appHttpHealth.Spec.Target = target

		// request
		request := new(v1beta1.NetHttpRequest)
		request.PerRequestTimeoutInMS = requestTimeout
		request.DurationInSecond = 10
		appHttpHealth.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		appHttpHealth.Spec.Schedule = Schedule

		maxQPS := 300
		appHttpHealth.Spec.Request.QPS = maxQPS

		for {
			appHttpHealthName := "apphttphealth-get" + tools.RandomName()
			appHttpHealth.Name = appHttpHealthName
			appHttpHealth.ResourceVersion = ""
			e = frame.CreateResource(appHttpHealth)
			Expect(e).NotTo(HaveOccurred(), "create appHttpHealth resource")

			e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
			Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

			time.Sleep(time.Second * 10)
			r, e := common.GetPluginReportResult(frame, appHttpHealthName, int(replicas))
			Expect(e).NotTo(HaveOccurred(), "get report failed")

			Reports := *r.Spec.Report

			if !Reports[0].HttpAppHealthyTask.Succeed {
				e = frame.DeleteResource(appHttpHealth)
				Expect(e).NotTo(HaveOccurred(), "delete resource failed")
				break
			} else {
				maxQPS = appHttpHealth.Spec.Request.QPS
				appHttpHealth.Spec.Request.QPS += 100
			}
			e = frame.DeleteResource(appHttpHealth)
			Expect(e).NotTo(HaveOccurred(), "delete resource failed")
		}

		GinkgoWriter.Printf("max QPS in AppHttpHealthy is %d \n", maxQPS)
	})

})
