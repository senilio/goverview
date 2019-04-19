// ESXi Overview in golang

package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	//"sort"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

// VM struct
type VM struct {
	Name        string
	MemoryMB    int32
	NumCPU      int32
	RuntimeHost string
	PowerState  interface{}
	Annotation  string
}

// Host struct
type Host struct {
	id string
}

func main() {
	ctx := context.Background()

	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse URL from environment
	u := os.Getenv("GOVERVIEW_URL")

	url, err := url.Parse(u)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// Connecting to vCenter
	c, err := govmomi.NewClient(ctx, url, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// Create view of VirtualMachine objects
	mv := view.NewManager(c.Client)
	v, err := mv.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		log.Fatal(err)
	}
	defer v.Destroy(ctx)

	// Create view of HostSystem objects
	mh := view.NewManager(c.Client)
	h, err := mh.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		log.Fatal(err)
	}
	defer h.Destroy(ctx)

	// Retrieve properties for all machines
	var vms []mo.VirtualMachine
	var vmprops = []string{
		"name",
		"snapshot",
		"layoutEx.file",
		"summary.quickStats.overallCpuUsage",
		"summary.quickStats.uptimeSeconds",
		"summary.customValue",
		"summary.quickStats",
		"summary.runtime.powerState",
		"runtime.host",
		"network",
		"guest.toolsRunningStatus",
		"guest.toolsVersionStatus",
		"guest.hostName",
		"guest.net",
		"guest.guestState",
		"runtime",
		"config",
	}

	err = v.Retrieve(ctx, []string{"VirtualMachine"}, vmprops, &vms)

	if err != nil {
		log.Fatal(err)
	}

	// Retrieve properties for hosts
	var hosts []mo.HostSystem
	var hostprops = []string{
		"summary.hardware.cpuModel",
		"config.product.fullName",
		"hardware.cpuInfo.numCpuCores",
		"hardware.cpuInfo.hz",
		"hardware.cpuInfo.numCpuThreads",
		"hardware.cpuInfo.numCpuPackages",
		"hardware.memorySize",
		"hardware.systemInfo.model",
		"hardware.systemInfo.vendor",
		"hardware.biosInfo.biosVersion",
		"name",
		"summary.hardware.cpuModel",
		"summary.hardware.numCpuPkgs",
		"summary.hardware.numNics",
		"summary.overallStatus",
	}

	err = h.Retrieve(ctx, []string{"HostSystem"}, hostprops, &hosts)

	if err != nil {
		log.Fatal(err)
	}

	/////////////////////////

	Hypervisors := map[string]*Host{}

	// Used to translate hypervisor ID to real hostname
	HypervisorTranslate := map[string]string{}

	for _, v := range hosts {
		HypervisorTranslate[v.Self.Value] = v.Name
		entry := new(Host)
		entry.id = v.Self.Value
		Hypervisors[v.Name] = entry
	}
	for k, v := range Hypervisors {
		fmt.Printf("k: %+v v: %+v\n", k, v)
	}

	VirtualMachines := map[string]*VM{}
	for _, vm := range vms {
		entry := new(VM)
		entry.Name = vm.Name
		entry.MemoryMB = vm.Config.Hardware.MemoryMB
		entry.NumCPU = vm.Config.Hardware.NumCPU
		entry.RuntimeHost = HypervisorTranslate[vm.Runtime.Host.Value]
		entry.PowerState = vm.Runtime.PowerState
		entry.Annotation = vm.Config.Annotation
		VirtualMachines[vm.Config.Uuid] = entry
	}

	fmt.Println("Before sort:")
	for k, v := range VirtualMachines {
		fmt.Printf("k: %+v v: %+v\n", k, v)
	}

}
