// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
)

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KdoctorReport
// +k8s:openapi-gen=true
type KdoctorReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec KdoctorReportSpec `json:"spec,omitempty"`
}

// KdoctorReportList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KdoctorReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []KdoctorReport `json:"items"`
}

// KdoctorReportSpec defines the desired state of KdoctorReport
type KdoctorReportSpec struct {
	TaskName            string    `json:"TaskName"`
	TaskType            string    `json:"TaskType"`
	ToTalRoundNumber    int64     `json:"RoundNumber"`
	FinishedRoundNumber int64     `json:"FinishedRoundNumber"`
	FailedRoundNumber   []int64   `json:"FailedRoundNumber"`
	Status              string    `json:"Status"`
	ReportRoundNumber   int64     `json:"ReportRoundNumber"`
	Report              *[]Report `json:"Report,omitempty"`
}

type Report struct {
	TaskName       string      `json:"TaskName"`
	TaskType       string      `json:"TaskType"`
	RoundNumber    int64       `json:"RoundNumber"`
	RoundResult    string      `json:"RoundResult"`
	NodeName       string      `json:"NodeName"`
	PodName        string      `json:"PodName"`
	FailedReason   *string     `json:"FailedReason,omitempty"`
	StartTimeStamp metav1.Time `json:"StartTimeStamp"`
	EndTimeStamp   metav1.Time `json:"EndTimeStamp"`
	RoundDuration  string      `json:"RoundDuration"`
	ReportType     string      `json:"ReportType"`

	NetReachTaskSpec *v1beta1.NetReachSpec `json:"NetReachTaskSpec,omitempty"`
	NetReachTask     *NetReachTask         `json:"NetReachTask,omitempty"`

	HttpAppHealthyTaskSpec *v1beta1.AppHttpHealthySpec `json:"HttpAppHealthyTaskSpec,omitempty"`
	HttpAppHealthyTask     *AppHttpHealthyTask         `json:"HttpAppHealthyTask,omitempty"`

	NetDNSTaskSpec *v1beta1.NetdnsSpec `json:"netDNSTaskSpec,omitempty"`
	NetDNSTask     *NetDNSTask         `json:"netDNSTask,omitempty"`
}
