package ctriface

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	ctrdlog "github.com/containerd/containerd/log"
	"github.com/containerd/containerd/namespaces"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestPauseSnapResume(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: ctrdlog.RFC3339NanoFixed,
		FullTimestamp:   true,
	})
	//log.SetReportCaller(true) // FIXME: make sure it's false unless debugging

	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)

	testTimeout := 120 * time.Second
	ctx, cancel := context.WithTimeout(namespaces.WithNamespace(context.Background(), namespaceName), testTimeout)
	defer cancel()

	orch := NewOrchestrator("devmapper", 1, true)

	vmID := "4"

	message, _, err := orch.StartVM(ctx, vmID, "ustiugov/helloworld:runner_workload")
	require.NoError(t, err, "Failed to start VM, "+message)

	message, err = orch.PauseVM(ctx, vmID)
	require.NoError(t, err, "Failed to pause VM, "+message)

	message, err = orch.CreateSnapshot(ctx, vmID, "/tmp/snapshot_file", "/tmp/mem_file")
	require.NoError(t, err, "Failed to create snapshot of VM, "+message)

	message, err = orch.ResumeVM(ctx, vmID)
	require.NoError(t, err, "Failed to resume VM, "+message)

	message, err = orch.StopSingleVM(ctx, vmID)
	require.NoError(t, err, "Failed to stop VM, "+message)

	orch.Cleanup()
}

func TestStartStopSerial(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: ctrdlog.RFC3339NanoFixed,
		FullTimestamp:   true,
	})
	//log.SetReportCaller(true) // FIXME: make sure it's false unless debugging

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)

	testTimeout := 120 * time.Second
	ctx, cancel := context.WithTimeout(namespaces.WithNamespace(context.Background(), namespaceName), testTimeout)
	defer cancel()

	orch := NewOrchestrator("devmapper", 1, true)

	vmID := "5"

	message, _, err := orch.StartVM(ctx, vmID, "ustiugov/helloworld:runner_workload")
	require.NoError(t, err, "Failed to start VM, "+message)

	message, err = orch.StopSingleVM(ctx, vmID)
	require.NoError(t, err, "Failed to stop VM, "+message)

	orch.Cleanup()
}

func TestPauseResumeSerial(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: ctrdlog.RFC3339NanoFixed,
		FullTimestamp:   true,
	})
	//log.SetReportCaller(true) // FIXME: make sure it's false unless debugging

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)

	testTimeout := 120 * time.Second
	ctx, cancel := context.WithTimeout(namespaces.WithNamespace(context.Background(), namespaceName), testTimeout)
	defer cancel()

	orch := NewOrchestrator("devmapper", 1, true)

	vmID := "6"

	message, _, err := orch.StartVM(ctx, vmID, "ustiugov/helloworld:runner_workload")
	require.NoError(t, err, "Failed to start VM, "+message)

	message, err = orch.PauseVM(ctx, vmID)
	require.NoError(t, err, "Failed to pause VM, "+message)

	message, err = orch.ResumeVM(ctx, vmID)
	require.NoError(t, err, "Failed to resume VM, "+message)

	message, err = orch.StopSingleVM(ctx, vmID)
	require.NoError(t, err, "Failed to stop VM, "+message)

	orch.Cleanup()
}

func TestStartStopParallel(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: ctrdlog.RFC3339NanoFixed,
		FullTimestamp:   true,
	})
	//log.SetReportCaller(true) // FIXME: make sure it's false unless debugging

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)

	testTimeout := 120 * time.Second
	ctx, cancel := context.WithTimeout(namespaces.WithNamespace(context.Background(), namespaceName), testTimeout)
	defer cancel()

	vmNum := 10
	orch := NewOrchestrator("devmapper", vmNum, true)

	{
		var vmGroup sync.WaitGroup
		for i := 0; i < vmNum; i++ {
			vmGroup.Add(1)
			go func(i int) {
				defer vmGroup.Done()
				vmID := fmt.Sprintf("test_%d", i)
				message, _, err := orch.StartVM(ctx, vmID, "ustiugov/helloworld:runner_workload")
				require.NoError(t, err, "Failed to start VM, "+message)
			}(i)
		}
		vmGroup.Wait()
	}

	{
		var vmGroup sync.WaitGroup
		for i := 0; i < vmNum; i++ {
			vmGroup.Add(1)
			go func(i int) {
				defer vmGroup.Done()
				vmID := fmt.Sprintf("test_%d", i)
				message, err := orch.StopSingleVM(ctx, vmID)
				require.NoError(t, err, "Failed to stop VM, "+message)
			}(i)
		}
		vmGroup.Wait()
	}

	orch.Cleanup()
}

func TestPauseResumeParallel(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: ctrdlog.RFC3339NanoFixed,
		FullTimestamp:   true,
	})
	//log.SetReportCaller(true) // FIXME: make sure it's false unless debugging

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)

	testTimeout := 120 * time.Second
	ctx, cancel := context.WithTimeout(namespaces.WithNamespace(context.Background(), namespaceName), testTimeout)
	defer cancel()

	vmNum := 10
	orch := NewOrchestrator("devmapper", vmNum, true)

	{
		var vmGroup sync.WaitGroup
		for i := 0; i < vmNum; i++ {
			vmGroup.Add(1)
			go func(i int) {
				defer vmGroup.Done()
				vmID := fmt.Sprintf("test_%d", i)
				message, _, err := orch.StartVM(ctx, vmID, "ustiugov/helloworld:runner_workload")
				require.NoError(t, err, "Failed to start VM, "+message)
			}(i)
		}
		vmGroup.Wait()
	}

	{
		var vmGroup sync.WaitGroup
		for i := 0; i < vmNum; i++ {
			vmGroup.Add(1)
			go func(i int) {
				defer vmGroup.Done()
				vmID := fmt.Sprintf("test_%d", i)
				message, err := orch.PauseVM(ctx, vmID)
				require.NoError(t, err, "Failed to pause VM, "+message)
			}(i)
		}
		vmGroup.Wait()
	}

	{
		var vmGroup sync.WaitGroup
		for i := 0; i < vmNum; i++ {
			vmGroup.Add(1)
			go func(i int) {
				defer vmGroup.Done()
				vmID := fmt.Sprintf("test_%d", i)
				message, err := orch.ResumeVM(ctx, vmID)
				require.NoError(t, err, "Failed to resume VM, "+message)
			}(i)
		}
		vmGroup.Wait()
	}

	{
		var vmGroup sync.WaitGroup
		for i := 0; i < vmNum; i++ {
			vmGroup.Add(1)
			go func(i int) {
				defer vmGroup.Done()
				vmID := fmt.Sprintf("test_%d", i)
				message, err := orch.StopSingleVM(ctx, vmID)
				require.NoError(t, err, "Failed to stop VM, "+message)
			}(i)
		}
		vmGroup.Wait()
	}

	orch.Cleanup()
}