// ESXi Overview in golang

package main

import (
    "context"
    "fmt"
    "log"

    "github.com/vmware/govmomi/examples"
    "github.com/vmware/govmomi/view"
    "github.com/vmware/govmomi/vim25/mo"
)

func main() {
    ctx := context.Background()

    // Connect and login to ESX or vCenter
    c, err := examples.NewClient(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer c.Logout(ctx)

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
        "guest.hostName",
        "config.guestFullName",
        "summary.quickStats.overallCpuUsage",
        "summary.quickStats.uptimeSeconds",
        "summary.customValue",
        "summary.quickStats",
        "config.hardware.memoryMB",
        "config.hardware.numCPU",
        "config.hardware.device",
        "config.annotation",
        "summary.runtime.powerState",
        "runtime.host",
        "network",
        "guest.toolsRunningStatus",
        "guest.toolsVersionStatus",
        "guest.net",
        "guest.guestState",
        "config.guestId",
        "config.version",
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

    // fmt.Println(vms)
    for _, vm := range vms {
        fmt.Printf("%s\t%s\t%s\t%d\t%s\n", vm.Name, vm.Config.GuestFullName, vm.Config.GuestId, vm.Config.Hardware.MemoryMB, vm.Config.Annotation)
    }
}
