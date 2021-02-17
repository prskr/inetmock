package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/inetmock/inetmock/internal/format"
	"gitlab.com/inetmock/inetmock/pkg/rpc"
	"go.uber.org/zap"
)

var (
	pcapCmd = &cobra.Command{
		Use:   "pcap",
		Short: "Interact with the PCAP API",
	}

	listAvailableDevicesCmd = &cobra.Command{
		Use:     "list-devices",
		Aliases: []string{"lis-dev", "ls-dev"},
		Short:   "List all devices that might be monitored",
		RunE:    runListAvailableDevices,
	}

	listCurrentlyRecordingsCmd = &cobra.Command{
		Use:     "list-recordings",
		Aliases: []string{"lis-rec", "ls-rec", "ls-recs"},
		Short:   "List currently active recordings",
		RunE:    runListActiveRecordings,
	}

	removeCurrentlyActiveRecording = &cobra.Command{
		Use:     "stop-recording",
		Aliases: []string{"rm-rec", "del-rec", "stop"},
		Short:   "Remove/stop a currently active recording",
		RunE:    runRemoveCurrentlyRunningRecording,
		ValidArgsFunction: func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
			var err error
			pcapClient := rpc.NewPCAPClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
			defer cancel()
			var resp *rpc.ListRecordingsResponse
			if resp, err = pcapClient.ListActiveRecordings(ctx, new(rpc.ListRecordingsRequest)); err == nil {
				return resp.Subscriptions, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveError
		},
	}

	addRecordingCmd = &cobra.Command{
		Use:     "start-recording",
		Aliases: []string{"start"},
		Short:   "[device] [targetPath] - adds a PCAP file subscription to the given path.",
		Long:    `If the path is relative it will be stored in the configured PCAP data directory.`,
		ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
			if len(args) > 2 {
				return nil, cobra.ShellCompDirectiveError
			}

			if len(args) == 2 {
				return nil, cobra.ShellCompDirectiveDefault
			}

			var err error
			pcapClient := rpc.NewPCAPClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
			defer cancel()
			var resp *rpc.ListAvailableDevicesResponse
			if resp, err = pcapClient.ListAvailableDevices(ctx, new(rpc.ListAvailableDevicesRequest)); err == nil {
				var completions []string

				for _, d := range resp.AvailableDevices {
					completions = append(completions, d.Name)
				}

				return completions, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveError
		},
		Args: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 2 {
				return errors.New("expected [device] [targetPath] as parameters")
			}

			if !strings.HasSuffix(strings.ToLower(args[1]), ".pcap") {
				return errors.New("expected .pcap suffix for the file name")
			}
			return nil
		},
		RunE: runAddRecording,
	}
)

func runListAvailableDevices(*cobra.Command, []string) (err error) {
	pcapClient := rpc.NewPCAPClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()

	var resp *rpc.ListAvailableDevicesResponse
	if resp, err = pcapClient.ListAvailableDevices(ctx, new(rpc.ListAvailableDevicesRequest)); err != nil {
		return
	}

	type printableDevice struct {
		Name      string
		Addresses string
	}

	availableDevs := make([]printableDevice, 0)

	for _, dev := range resp.AvailableDevices {
		availableDevs = append(availableDevs, printableDevice{
			Name:      dev.Name,
			Addresses: byteArraysToPrintableIPAddresses(dev.Addresses),
		})
	}

	writer := format.Writer(outputFormat, os.Stdout)
	err = writer.Write(availableDevs)
	return
}

func runListActiveRecordings(*cobra.Command, []string) (err error) {

	type printableSubscription struct {
		Name        string
		Device      string
		ConsumerKey string
	}

	pcapClient := rpc.NewPCAPClient(conn)

	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()

	var resp *rpc.ListRecordingsResponse
	if resp, err = pcapClient.ListActiveRecordings(ctx, new(rpc.ListRecordingsRequest)); err != nil {
		return
	}

	var out []printableSubscription

	for _, subscription := range resp.Subscriptions {
		nameDevSplit := strings.Split(subscription, ":")
		if len(nameDevSplit) != 2 {
			continue
		}
		out = append(out, printableSubscription{
			Name:        nameDevSplit[1],
			Device:      nameDevSplit[0],
			ConsumerKey: subscription,
		})
	}

	writer := format.Writer(outputFormat, os.Stdout)
	err = writer.Write(out)

	return
}

func runAddRecording(_ *cobra.Command, args []string) (err error) {
	pcapClient := rpc.NewPCAPClient(conn)

	if err = isValidRecordDevice(args[0], pcapClient); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()

	var resp *rpc.RegisterPCAPFileRecordResponse
	resp, err = pcapClient.StartPCAPFileRecording(ctx, &rpc.RegisterPCAPFileRecordRequest{
		Device:     args[0],
		TargetPath: args[1],
	})

	if err != nil {
		return
	}

	cliApp.Logger().Info("Added PCAP recording", zap.String("targetPath", resp.ResolvedPath))

	return
}

func runRemoveCurrentlyRunningRecording(*cobra.Command, []string) (err error) {
	return
}

func isValidRecordDevice(device string, pcapClient rpc.PCAPClient) (err error) {
	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()
	var resp *rpc.ListAvailableDevicesResponse
	if resp, err = pcapClient.ListAvailableDevices(ctx, new(rpc.ListAvailableDevicesRequest)); err != nil {
		return
	}

	for _, dev := range resp.AvailableDevices {
		if strings.ToLower(dev.Name) == strings.ToLower(device) {
			return nil
		}
	}

	return fmt.Errorf("device %s not found in available devices", device)
}

func byteArraysToPrintableIPAddresses(arrs [][]byte) string {
	ipsArr := make([]string, 0)
	for _, b := range arrs {
		ip := net.IP(b)
		ipsArr = append(ipsArr, ip.String())
	}

	return strings.Join(ipsArr, ", ")
}
