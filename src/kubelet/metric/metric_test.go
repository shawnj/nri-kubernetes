package metric

import (
	"encoding/json"
	"fmt"
	"testing"

	v1 "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"

	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/definition"
	"github.com/stretchr/testify/assert"
)

var responseOKDataSample = `{ "node": { "nodeName": "fooNode", "startTime": "2018-01-22T06:52:15Z", "cpu": { "time": "2018-01-24T16:40:00Z", "usageNanoCores": 64124211, "usageCoreNanoSeconds": 353998913059080 }, "memory": { "time": "2018-01-24T16:40:00Z", "availableBytes": 502603776, "usageBytes": 687067136, "workingSetBytes": 540618752, "rssBytes": 150396928, "pageFaults": 3067606235, "majorPageFaults": 517653 }, "network": { "time": "2018-01-24T16:40:00Z", "rxBytes": 51419684038, "rxErrors": 0, "txBytes": 25630208577, "txErrors": 0 }, "fs": { "time": "2018-01-24T16:40:00Z", "availableBytes": 92795400192, "capacityBytes": 128701009920, "usedBytes": 30305800192, "inodesFree": 32999604, "inodes": 33554432, "inodesUsed": 554828 }, "runtime": { "imageFs": { "time": "2018-01-24T16:40:00Z", "availableBytes": 92795400192, "capacityBytes": 128701009920, "usedBytes": 20975835934, "inodesFree": 32999604, "inodes": 33554432, "inodesUsed": 554828 } } }, "pods": [ { "podRef": { "name": "newrelic-infra-monitoring-pjp0v", "namespace": "kube-system", "uid": "b5a9c98f-d34f-11e7-95fe-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:37Z", "containers": [ { "name": "kube-state-metrics", "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 }, "rootfs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 35000320, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 24 }, "logs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 20480, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null }, { "name": "newrelic-infra", "startTime": "2017-11-30T09:12:44Z", "cpu": { "time": "2017-11-30T14:48:12Z", "usageNanoCores": 13046199, "usageCoreNanoSeconds": 303855795298 }, "memory": { "time": "2017-11-30T14:48:12Z", "usageBytes": 243638272, "workingSetBytes": 38313984, "rssBytes": 15785984, "pageFaults": 10304448, "majorPageFaults": 217 }, "rootfs": { "time": "2017-11-30T14:48:12Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 1305837568, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 52 }, "logs": { "time": "2017-11-30T14:48:12Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 657747968, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null } ], "network": { "time": "2017-11-30T14:48:12Z", "rxBytes": 15741653, "rxErrors": 0, "txBytes": 19551073, "txErrors": 0 }, "volume": [ { "time": "2017-11-30T09:13:29Z", "availableBytes": 1048637440, "capacityBytes": 1048649728, "usedBytes": 12288, "inodesFree": 256009, "inodes": 256018, "inodesUsed": 9, "name": "default-token-7cg8m" } ] }, { "podRef": { "name": "kube-dns-910330662-pflkj", "namespace": "kube-system", "uid": "a6f2130b-a21e-11e7-8db6-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:36Z", "containers": [ { "name": "dnsmasq", "startTime": "2017-11-30T09:12:43Z", "cpu": { "time": "2017-11-30T14:48:07Z", "usageNanoCores": 208374, "usageCoreNanoSeconds": 3653471654 }, "memory": { "time": "2017-11-30T14:48:07Z", "usageBytes": 19812352, "workingSetBytes": 12828672, "rssBytes": 5201920, "pageFaults": 3376, "majorPageFaults": 139 }, "rootfs": { "time": "2017-11-30T14:48:07Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 42041344, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 20 }, "logs": { "time": "2017-11-30T14:48:07Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 20480, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null } ], "network": { "time": "2017-11-30T14:48:07Z", "rxBytes": 14447980, "rxErrors": 0, "txBytes": 15557657, "txErrors": 0 }, "volume": [ { "time": "2017-11-30T09:13:29Z", "availableBytes": 1048637440, "capacityBytes": 1048649728, "usedBytes": 12288, "inodesFree": 256009, "inodes": 256018, "inodesUsed": 9, "name": "default-token-7cg8m" } ] } ] }`
var responseContainerWithTheSameName = `{ "pods": [ { "podRef": { "name": "newrelic-infra-monitoring-pjp0v", "namespace": "kube-system", "uid": "b5a9c98f-d34f-11e7-95fe-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:37Z", "containers": [ { "name": "kube-state-metrics", "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 }, "rootfs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 35000320, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 24 }, "logs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 20480, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null }, { "name": "newrelic-infra", "startTime": "2017-11-30T09:12:44Z", "cpu": { "time": "2017-11-30T14:48:12Z", "usageNanoCores": 13046199, "usageCoreNanoSeconds": 303855795298 }, "memory": { "time": "2017-11-30T14:48:12Z", "usageBytes": 243638272, "workingSetBytes": 38313984, "rssBytes": 15785984, "pageFaults": 10304448, "majorPageFaults": 217 }, "rootfs": { "time": "2017-11-30T14:48:12Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 1305837568, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 52 }, "logs": { "time": "2017-11-30T14:48:12Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 657747968, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null } ], "network": { "time": "2017-11-30T14:48:12Z", "rxBytes": 15741653, "rxErrors": 0, "txBytes": 19551073, "txErrors": 0 }, "volume": [ { "time": "2017-11-30T09:13:29Z", "availableBytes": 1048637440, "capacityBytes": 1048649728, "usedBytes": 12288, "inodesFree": 256009, "inodes": 256018, "inodesUsed": 9, "name": "default-token-7cg8m" } ] }, { "podRef": { "name": "kube-dns-910330662-pflkj", "namespace": "kube-system", "uid": "a6f2130b-a21e-11e7-8db6-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:36Z", "containers": [ { "name": "kube-state-metrics", "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 }, "rootfs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 35000320, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 24 }, "logs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 20480, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null }, { "name": "dnsmasq", "startTime": "2017-11-30T09:12:43Z", "cpu": { "time": "2017-11-30T14:48:07Z", "usageNanoCores": 208374, "usageCoreNanoSeconds": 3653471654 }, "memory": { "time": "2017-11-30T14:48:07Z", "usageBytes": 19812352, "workingSetBytes": 12828672, "rssBytes": 5201920, "pageFaults": 3376, "majorPageFaults": 139 }, "rootfs": { "time": "2017-11-30T14:48:07Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 42041344, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 20 }, "logs": { "time": "2017-11-30T14:48:07Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 20480, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null } ], "network": { "time": "2017-11-30T14:48:07Z", "rxBytes": 14447980, "rxErrors": 0, "txBytes": 15557657, "txErrors": 0 }, "volume": [ { "time": "2017-11-30T09:13:29Z", "availableBytes": 1048637440, "capacityBytes": 1048649728, "usedBytes": 12288, "inodesFree": 256009, "inodes": 256018, "inodesUsed": 9, "name": "default-token-7cg8m" } ] } ] }`
var responseMissingContainerName = `{ "pods": [ { "podRef": { "name": "newrelic-infra-monitoring-pjp0v", "namespace": "kube-system", "uid": "b5a9c98f-d34f-11e7-95fe-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:37Z", "containers": [ { "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 } } ], "network": { "time": "2017-11-30T14:48:12Z", "rxBytes": 15741653, "txBytes": 52463212, "rxErrors": 0,  "txErrors": 0 } } ] }`
var responseMissingPodName = `{ "pods": [ { "podRef": { "namespace": "kube-system", "uid": "b5a9c98f-d34f-11e7-95fe-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:37Z", "containers": [ { "name": "kube-state-metrics", "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 } } ], "network": { "time": "2017-11-30T14:48:12Z", "rxBytes": 15741653, "txBytes": 52463212, "rxErrors": 0,  "txErrors": 0 } } ] }`
var responseMissingRxBytesForPod = `{ "pods": [ { "podRef": { "name": "newrelic-infra-monitoring-pjp0v", "namespace": "kube-system", "uid": "b5a9c98f-d34f-11e7-95fe-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:37Z", "containers": [ { "name": "kube-state-metrics", "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 } } ], "network": { "time": "2017-11-30T14:48:12Z", "txBytes": 52463212, "rxErrors": 0,  "txErrors": 0 } } ] }`

var nodeSampleMissingImageFs = `{ "node": { "nodeName": "fooNode", "startTime": "2018-01-22T06:52:15Z", "cpu": { "time": "2018-01-24T16:40:00Z", "usageNanoCores": 64124211, "usageCoreNanoSeconds": 353998913059080 }, "memory": { "time": "2018-01-24T16:40:00Z", "availableBytes": 502603776, "usageBytes": 687067136, "workingSetBytes": 540618752, "rssBytes": 150396928, "pageFaults": 3067606235, "majorPageFaults": 517653 }, "network": { "time": "2018-01-24T16:40:00Z", "rxBytes": 51419684038, "rxErrors": 0, "txBytes": 25630208577, "txErrors": 0 }, "fs": { "time": "2018-01-24T16:40:00Z", "availableBytes": 92795400192, "capacityBytes": 128701009920, "usedBytes": 30305800192, "inodesFree": 32999604, "inodes": 33554432, "inodesUsed": 554828 }, "runtime": { } } }`

func toSummary(response string) (v1.Summary, error) {
	var summary = new(v1.Summary)
	err := json.Unmarshal([]byte(response), summary)
	if err != nil {
		return v1.Summary{}, fmt.Errorf("Error unmarshaling the response body. Got error: %v", err.Error())
	}
	return *summary, nil
}

func TestGroupStatsSummary_CorrectValue(t *testing.T) {
	expectedRawData := definition.RawGroups{
		"node": {
			"fooNode": definition.RawMetrics{
				"nodeName": "fooNode",
				// CPU
				"usageNanoCores":       uint64(64124211),
				"usageCoreNanoSeconds": uint64(353998913059080),
				// Memory
				"memoryUsageBytes":      uint64(687067136),
				"memoryAvailableBytes":  uint64(502603776),
				"memoryWorkingSetBytes": uint64(540618752),
				"memoryRssBytes":        uint64(150396928),
				"memoryPageFaults":      uint64(3067606235),
				"memoryMajorPageFaults": uint64(517653),
				// Network
				"rxBytes": uint64(51419684038),
				"txBytes": uint64(25630208577),
				"errors":  uint64(0),
				// Fs
				"fsAvailableBytes": uint64(92795400192),
				"fsCapacityBytes":  uint64(128701009920),
				"fsUsedBytes":      uint64(30305800192),
				"fsInodesFree":     uint64(32999604),
				"fsInodes":         uint64(33554432),
				"fsInodesUsed":     uint64(554828),
				// Runtime
				"runtimeAvailableBytes": uint64(92795400192),
				"runtimeCapacityBytes":  uint64(128701009920),
				"runtimeUsedBytes":      uint64(20975835934),
				"runtimeInodesFree":     uint64(32999604),
				"runtimeInodes":         uint64(33554432),
				"runtimeInodesUsed":     uint64(554828),
			},
		},
		"pod": {
			"kube-system_newrelic-infra-monitoring-pjp0v": definition.RawMetrics{
				"podName":   "newrelic-infra-monitoring-pjp0v",
				"namespace": "kube-system",
				"rxBytes":   uint64(15741653),
				"errors":    uint64(0),
				"txBytes":   uint64(19551073),
			},
			"kube-system_kube-dns-910330662-pflkj": definition.RawMetrics{
				"podName":   "kube-dns-910330662-pflkj",
				"namespace": "kube-system",
				"rxBytes":   uint64(14447980),
				"errors":    uint64(0),
				"txBytes":   uint64(15557657),
			},
		},
		"container": {
			"kube-system_newrelic-infra-monitoring-pjp0v_kube-state-metrics": definition.RawMetrics{
				"containerName":  "kube-state-metrics",
				"usageBytes":     uint64(22552576),
				"usageNanoCores": uint64(184087),
				"podName":        "newrelic-infra-monitoring-pjp0v",
				"namespace":      "kube-system",
			},
			"kube-system_newrelic-infra-monitoring-pjp0v_newrelic-infra": definition.RawMetrics{
				"containerName":  "newrelic-infra",
				"usageBytes":     uint64(243638272),
				"usageNanoCores": uint64(13046199),
				"podName":        "newrelic-infra-monitoring-pjp0v",
				"namespace":      "kube-system",
			},
			"kube-system_kube-dns-910330662-pflkj_dnsmasq": definition.RawMetrics{
				"containerName":  "dnsmasq",
				"usageBytes":     uint64(19812352),
				"usageNanoCores": uint64(208374),
				"podName":        "kube-dns-910330662-pflkj",
				"namespace":      "kube-system",
			},
		},
	}
	summary, err := toSummary(responseOKDataSample)
	assert.NoError(t, err)

	rawData, errs := GroupStatsSummary(summary)
	assert.Empty(t, errs)
	assert.Equal(t, expectedRawData, rawData)
}

func TestGroupStatsSummary_MissingNodeData_ContainerWithTheSameName(t *testing.T) {
	expectedRawData := definition.RawGroups{
		"pod": {
			"kube-system_newrelic-infra-monitoring-pjp0v": definition.RawMetrics{
				"podName":   "newrelic-infra-monitoring-pjp0v",
				"namespace": "kube-system",
				"rxBytes":   uint64(15741653),
				"errors":    uint64(0),
				"txBytes":   uint64(19551073),
			},
			"kube-system_kube-dns-910330662-pflkj": definition.RawMetrics{
				"podName":   "kube-dns-910330662-pflkj",
				"namespace": "kube-system",
				"rxBytes":   uint64(14447980),
				"errors":    uint64(0),
				"txBytes":   uint64(15557657),
			},
		},
		"container": {
			"kube-system_newrelic-infra-monitoring-pjp0v_kube-state-metrics": definition.RawMetrics{
				"containerName":  "kube-state-metrics",
				"usageBytes":     uint64(22552576),
				"usageNanoCores": uint64(184087),
				"podName":        "newrelic-infra-monitoring-pjp0v",
				"namespace":      "kube-system",
			},
			"kube-system_kube-dns-910330662-pflkj_kube-state-metrics": definition.RawMetrics{
				"containerName":  "kube-state-metrics",
				"usageBytes":     uint64(22552576),
				"usageNanoCores": uint64(184087),
				"podName":        "kube-dns-910330662-pflkj",
				"namespace":      "kube-system",
			},
			"kube-system_newrelic-infra-monitoring-pjp0v_newrelic-infra": definition.RawMetrics{
				"containerName":  "newrelic-infra",
				"usageBytes":     uint64(243638272),
				"usageNanoCores": uint64(13046199),
				"podName":        "newrelic-infra-monitoring-pjp0v",
				"namespace":      "kube-system",
			},
			"kube-system_kube-dns-910330662-pflkj_dnsmasq": definition.RawMetrics{
				"containerName":  "dnsmasq",
				"usageBytes":     uint64(19812352),
				"usageNanoCores": uint64(208374),
				"podName":        "kube-dns-910330662-pflkj",
				"namespace":      "kube-system",
			},
		},
		"node": {},
	}
	summary, err := toSummary(responseContainerWithTheSameName)
	assert.NoError(t, err)

	rawData, errs := GroupStatsSummary(summary)
	assert.EqualError(t, errs[0], "empty node identifier, possible data error in /stats/summary response")
	assert.Equal(t, expectedRawData, rawData)
}

func TestGroupStatsSummary_IncompleteStatsSummaryMessage_MissingNodeData_MissingContainerName(t *testing.T) {
	expectedRawData := definition.RawGroups{
		"pod": {
			"kube-system_newrelic-infra-monitoring-pjp0v": definition.RawMetrics{
				"podName":   "newrelic-infra-monitoring-pjp0v",
				"namespace": "kube-system",
				"rxBytes":   uint64(15741653),
				"txBytes":   uint64(52463212),
				"errors":    uint64(0),
			},
		},
		"container": {},
		"node":      {},
	}

	summary, err := toSummary(responseMissingContainerName)
	assert.NoError(t, err)

	rawData, errs := GroupStatsSummary(summary)
	assert.Len(t, errs, 2, "Not expected length of errors")
	assert.Equal(t, expectedRawData, rawData)
}

func TestGroupStatsSummary_IncompleteStatsSummaryMessage_MissingNodeData_MissingPodName(t *testing.T) {
	expectedRawData := definition.RawGroups{
		"pod":       {},
		"container": {},
		"node":      {},
	}

	summary, err := toSummary(responseMissingPodName)
	assert.NoError(t, err)

	rawData, errs := GroupStatsSummary(summary)
	assert.Len(t, errs, 2, "Not expected length of errors")
	assert.Len(t, rawData, 3, "Not expected length of rawData for pods and containers")
	assert.Equal(t, expectedRawData, rawData)
	assert.Empty(t, rawData["pod"])
	assert.Empty(t, rawData["container"])
	assert.Empty(t, rawData["node"])
}

func TestGroupStatsSummary_IncompleteStatsSummaryMessage_MissingNodeData_NoRxBytesForPod_ReportedAsZero(t *testing.T) {
	expectedRawData := definition.RawGroups{
		"pod": {
			"kube-system_newrelic-infra-monitoring-pjp0v": definition.RawMetrics{
				"podName":   "newrelic-infra-monitoring-pjp0v",
				"namespace": "kube-system",
				"errors":    uint64(0),
				"txBytes":   uint64(52463212),
			},
		},
		"container": {
			"kube-system_newrelic-infra-monitoring-pjp0v_kube-state-metrics": definition.RawMetrics{
				"containerName":  "kube-state-metrics",
				"usageBytes":     uint64(22552576),
				"usageNanoCores": uint64(184087),
				"podName":        "newrelic-infra-monitoring-pjp0v",
				"namespace":      "kube-system",
			},
		},
		"node": {},
	}

	summary, err := toSummary(responseMissingRxBytesForPod)
	assert.NoError(t, err)

	rawData, errs := GroupStatsSummary(summary)
	assert.Len(t, errs, 1, "Not expected length of errors")
	assert.Equal(t, expectedRawData, rawData)
}

func TestGroupStatsSummary_EmptyStatsSummaryMessage(t *testing.T) {
	var summary = new(v1.Summary)

	rawData, errs := GroupStatsSummary(*summary)

	assert.Len(t, errs, 2, "Not expected length of errors")
	assert.Len(t, rawData, 3, "Not expected length of rawData for pods and containers")
	assert.Empty(t, rawData["pod"])
	assert.Empty(t, rawData["container"])
	assert.Empty(t, rawData["node"])
}

func TestAddUint64RawMetric(t *testing.T) {
	r := definition.RawMetrics{
		"nodeName": "fooNode",
	}

	expected := definition.RawMetrics{
		"nodeName": "fooNode",
		"foo":      uint64(353998913059080),
	}

	summary, err := toSummary(`{ "node": { "cpu": { "usageCoreNanoSeconds": 353998913059080 } } }`)
	assert.NoError(t, err)

	AddUint64RawMetric(r, "foo", summary.Node.CPU.UsageCoreNanoSeconds)
	assert.Equal(t, expected, r)
}

func TestFetchNodeStats_MissingImageFs(t *testing.T) {
	expectedRawData := definition.RawMetrics{
		"nodeName": "fooNode",
		// CPU
		"usageNanoCores":       uint64(64124211),
		"usageCoreNanoSeconds": uint64(353998913059080),
		// Memory
		"memoryUsageBytes":      uint64(687067136),
		"memoryAvailableBytes":  uint64(502603776),
		"memoryWorkingSetBytes": uint64(540618752),
		"memoryRssBytes":        uint64(150396928),
		"memoryPageFaults":      uint64(3067606235),
		"memoryMajorPageFaults": uint64(517653),
		// Network
		"rxBytes": uint64(51419684038),
		"txBytes": uint64(25630208577),
		"errors":  uint64(0),
		// Fs
		"fsAvailableBytes": uint64(92795400192),
		"fsCapacityBytes":  uint64(128701009920),
		"fsUsedBytes":      uint64(30305800192),
		"fsInodesFree":     uint64(32999604),
		"fsInodes":         uint64(33554432),
		"fsInodesUsed":     uint64(554828),
	}
	s, err := toSummary(nodeSampleMissingImageFs)
	assert.NoError(t, err)

	rawData, ID, errs := fetchNodeStats(s.Node)
	assert.Empty(t, errs)
	assert.Equal(t, "fooNode", ID)
	assert.Equal(t, expectedRawData, rawData)
}
