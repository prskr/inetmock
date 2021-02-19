package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/format"
	"gitlab.com/inetmock/inetmock/pkg/rpc"
)

const (
	expectedAddRecordsArgsLength = 2
	defaultReadTimeout           = 30 * time.Second
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
	promiscuousMode bool
	readTimeout     time.Duration
	addRecordingCmd = &cobra.Command{
		Use:     "start-recording",
		Aliases: []string{"start"},
		Short:   "[device] [targetPath] - adds a PCAP file subscription to the given path.",
		Long:    `If the path is relative it will be stored in the configured PCAP data directory.`,
		ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
			if len(args) > expectedAddRecordsArgsLength {
				return nil, cobra.ShellCompDirectiveError
			}

			if len(args) == expectedAddRecordsArgsLength {
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
			if len(args) != expectedAddRecordsArgsLength {
				return errors.New("expected [device] [targetPath] as parameters")
			}

			if !strings.HasSuffix(strings.ToLower(args[1]), ".pcap") {
				return errors.New("expected .pcap suffix for the file name")
			}
			return nil
		},
		RunE: runAddRecording,
	}
	removeCurrentlyActiveRecording = &cobra.Command{
		Use:     "stop-recording",
		Aliases: []string{"rm-rec", "del-rec", "stop"},
		Short:   "Remove/stop a currently active recording",
		RunE:    runRemoveCurrentlyRunningRecording,
		Args:    cobra.ExactArgs(1),
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
)

func init() {
	addRecordingCmd.Flags().BoolVar(
		&promiscuousMode,
		"promiscuous",
		false,
		"Start the recording in promiscuous mode which means it also captures traffic not only meant for the given interface",
	)

	addRecordingCmd.Flags().DurationVar(
		&readTimeout,
		"read-timeout",
		defaultReadTimeout,
		"configure the read time for the recording - supported values are Go time.Duration strings",
	)

	pcapCmd.AddCommand(
		listAvailableDevicesCmd,
		listCurrentlyRecordingsCmd,
		addRecordingCmd,
		removeCurrentlyActiveRecording,
	)
}

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

func runListActiveRecordings(*cobra.Command, []string) error {
	type printableSubscription struct {
		Name        string
		Device      string
		ConsumerKey string
	}

	pcapClient := rpc.NewPCAPClient(conn)

	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()

	var err error
	var resp *rpc.ListRecordingsResponse
	if resp, err = pcapClient.ListActiveRecordings(ctx, new(rpc.ListRecordingsRequest)); err != nil {
		return err
	}

	var out []printableSubscription
	for _, subscription := range resp.Subscriptions {
		splitIdx := strings.Index(subscription, ":")
		if splitIdx < 0 {
			continue
		}

		out = append(out, printableSubscription{
			Name:        subscription[splitIdx:],
			Device:      subscription[:splitIdx],
			ConsumerKey: subscription,
		})
	}

	writer := format.Writer(outputFormat, os.Stdout)

	return writer.Write(out)
}

func runAddRecording(_ *cobra.Command, args []string) (err error) {
	pcapClient := rpc.NewPCAPClient(conn)

	if err = isValidRecordDevice(args[0], pcapClient); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()

	var resp *rpc.StartPCAPFileRecordResponse
	resp, err = pcapClient.StartPCAPFileRecording(ctx, &rpc.StartPCAPFileRecordRequest{
		Device:     args[0],
		TargetPath: args[1],
	})

	if err != nil {
		return
	}

	cliApp.Logger().Info("Added PCAP recording", zap.String("targetPath", resp.ResolvedPath))

	return
}

func runRemoveCurrentlyRunningRecording(_ *cobra.Command, args []string) error {
	pcapClient := rpc.NewPCAPClient(conn)

	listRecsCtx, listRecsCancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer listRecsCancel()

	var err error
	var listRecsResp *rpc.ListRecordingsResponse
	if listRecsResp, err = pcapClient.ListActiveRecordings(listRecsCtx, new(rpc.ListRecordingsRequest)); err != nil {
		return err
	}

	var knownSubscription = false
	for i := range listRecsResp.Subscriptions {
		knownSubscription = knownSubscription || listRecsResp.Subscriptions[i] == args[0]
		if knownSubscription {
			break
		}
	}

	if !knownSubscription {
		return fmt.Errorf("the given subscription is not known: %s", args[0])
	}

	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()

	var stopRecResp *rpc.StopPCAPFileRecordResponse
	stopRecResp, err = pcapClient.StopPCAPFileRecord(ctx, &rpc.StopPCAPFileRecordRequest{
		ConsumerKey: args[0],
	})

	if err != nil {
		return err
	}

	if !stopRecResp.Removed {
		return fmt.Errorf("apparently no recording got removed for the given key %s", args[0])
	}
	return nil
}

func isValidRecordDevice(device string, pcapClient rpc.PCAPClient) (err error) {
	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()
	var resp *rpc.ListAvailableDevicesResponse
	if resp, err = pcapClient.ListAvailableDevices(ctx, new(rpc.ListAvailableDevicesRequest)); err != nil {
		return
	}

	for _, dev := range resp.AvailableDevices {
		if strings.EqualFold(dev.Name, device) {
			return nil
		}
	}

	return fmt.Errorf("device %s not found in available devices", device)
}

func byteArraysToPrintableIPAddresses(addresses [][]byte) string {
	ipsArr := make([]string, 0)
	for _, b := range addresses {
		ip := net.IP(b)
		ipsArr = append(ipsArr, ip.String())
	}

	return strings.Join(ipsArr, ", ")
}
