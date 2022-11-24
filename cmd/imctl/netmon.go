package main

import (
	"context"
	"net"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/format"
	"inetmock.icb4dc0.de/inetmock/internal/netutils"
	rpcv1 "inetmock.icb4dc0.de/inetmock/pkg/rpc/v1"
)

var (
	netMonCmd = &cobra.Command{
		Use:   "netflow",
		Short: "Interact with the network monitoring API",
	}
	listAvailableNetInterfacesCmd = &cobra.Command{
		Use:          "list-devices",
		Aliases:      []string{"lis-dev", "ls-dev"},
		Short:        "List all devices that might be monitored",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return runListAvailableNetInterfaces()
		},
	}
	listRunningNetworkMonitors = &cobra.Command{
		Use:          "list-monitors",
		Aliases:      []string{"ls-mons", "ls-monitors"},
		Short:        "List all running monitors",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return runListRunningMonitors()
		},
	}
	startInterfaceMonitorCmd = &cobra.Command{
		Use:          "start",
		Aliases:      []string{"create"},
		Short:        "Start a new network interface monitor",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runStartInterfaceMonitor(args[0], startInterfaceMonitorRequest)
		},
	}
	stopInterfaceMonitorCmd = &cobra.Command{
		Use:          "stop",
		Aliases:      []string{"rm"},
		Short:        "Stop monitoring a network interface",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runStopInterfaceMonitor(args[0])
		},
	}
	startInterfaceMonitorRequest interfaceMonitorRequest
)

type interfaceMonitorRequest struct {
	WhitelistPorts         []int
	SourceIPWhitelist      []net.IP
	DestinationIPWhitelist []net.IP
	InterfaceToMonitor     string
	RemoveMemLock          bool
}

//nolint:lll
func init() {
	startInterfaceMonitorCmd.Flags().IntSliceVar(&startInterfaceMonitorRequest.WhitelistPorts, "whitelist-ports", nil, "List of destination ports that are ignored in the monitor")
	startInterfaceMonitorCmd.Flags().IPSliceVar(&startInterfaceMonitorRequest.SourceIPWhitelist, "whitelist-source-ip", nil, "Source IP addresses that should be ignored")
	startInterfaceMonitorCmd.Flags().IPSliceVar(&startInterfaceMonitorRequest.DestinationIPWhitelist, "whitelist-destination-ip", nil, "Destination IP addresses that should be ignored")
	startInterfaceMonitorCmd.Flags().BoolVar(&startInterfaceMonitorRequest.RemoveMemLock, "remove-memlock", false, "Remove rlimit memlock for Linux kernels < 5.11")
	netMonCmd.AddCommand(listAvailableNetInterfacesCmd, listRunningNetworkMonitors, startInterfaceMonitorCmd, stopInterfaceMonitorCmd)
}

func runListAvailableNetInterfaces() (err error) {
	netMonClient := rpcv1.NewNetFlowControlServiceClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), cfg.GRPCTimeout)
	defer cancel()

	var resp *rpcv1.ListAvailableNetworkInterfacesResponse
	if resp, err = netMonClient.ListAvailableNetworkInterfaces(ctx, new(rpcv1.ListAvailableNetworkInterfacesRequest)); err != nil {
		return err
	}
	type printableDevice struct {
		Name      string
		Addresses string
	}

	availableDevs := make([]printableDevice, 0, len(resp.AvailableDevices))

	for idx := range resp.AvailableDevices {
		availableDevs = append(availableDevs, printableDevice{
			Name:      resp.AvailableDevices[idx].Name,
			Addresses: byteArraysToPrintableIPAddresses(resp.AvailableDevices[idx].Addresses),
		})
	}

	writer := format.Writer(cfg.Format, os.Stdout)
	return writer.Write(availableDevs)
}

func runListRunningMonitors() (err error) {
	netMonClient := rpcv1.NewNetFlowControlServiceClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), cfg.GRPCTimeout)
	defer cancel()

	var resp *rpcv1.ListControlledInterfacesResponse
	if resp, err = netMonClient.ListControlledInterfaces(ctx, new(rpcv1.ListControlledInterfacesRequest)); err != nil {
		return err
	}

	type printableRunningMonitor struct {
		InterfaceName string
	}

	runningMonitors := make([]printableRunningMonitor, 0, len(resp.ControlledInterfaces))

	for idx := range resp.ControlledInterfaces {
		runningMonitors = append(runningMonitors, printableRunningMonitor{InterfaceName: resp.ControlledInterfaces[idx]})
	}

	writer := format.Writer(cfg.Format, os.Stdout)
	return writer.Write(runningMonitors)
}

func runStartInterfaceMonitor(interfaceName string, req interfaceMonitorRequest) (err error) {
	netMonClient := rpcv1.NewNetFlowControlServiceClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), cfg.GRPCTimeout)
	defer cancel()

	rpcReq := &rpcv1.StartPacketFlowControlRequest{
		InterfaceName:          interfaceName,
		SourceIpWhitelist:      netutils.IPAddressesToBytes(req.SourceIPWhitelist),
		DestinationIpWhitelist: netutils.IPAddressesToBytes(req.DestinationIPWhitelist),
		WhitelistPorts:         make([]uint32, 0, len(req.WhitelistPorts)),
		RemoveRlimitMemlock:    req.RemoveMemLock,
	}

	for idx := range req.WhitelistPorts {
		rpcReq.WhitelistPorts = append(rpcReq.WhitelistPorts, uint32(req.WhitelistPorts[idx]))
	}

	if _, err = netMonClient.StartPacketFlowControl(ctx, rpcReq); err != nil {
		return err
	}

	return nil
}

func runStopInterfaceMonitor(interfaceName string) error {
	netMonClient := rpcv1.NewNetFlowControlServiceClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), cfg.GRPCTimeout)
	defer cancel()

	resp, err := netMonClient.StopPacketFlowControl(ctx, &rpcv1.StopPacketFlowControlRequest{
		InterfaceName: interfaceName,
	})
	if err != nil {
		return err
	}

	cliApp.Logger().Info("Stop monitor result", zap.Bool("monitor_was_running", resp.InterfaceWasControlled))

	return nil
}
