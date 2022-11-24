package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"

	rpcv1 "inetmock.icb4dc0.de/inetmock/pkg/rpc/v1"
)

const (
	outFileMode              fs.FileMode = 0o600
	defaultProfilingDuration             = 30 * time.Second
)

var (
	pprofCmd = &cobra.Command{
		Use:   "pprof",
		Short: "Fetch pprof dumps via gRPC socket/API",
	}

	dumpProfileCmd = &cobra.Command{
		Use:          "dump",
		Short:        "Dump a pprof profile:",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runDumpProfile(args[0])
		},
	}

	dumpCPUProfile = &cobra.Command{
		Use:          "cpu-profile",
		Short:        "Dump a CPU profile for a given amount of time",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return runDumpCPUProfile()
		},
	}

	dumpTraceCmd = &cobra.Command{
		Use:          "trace",
		Short:        "Dump a trace",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return runDumpTrace()
		},
	}

	pprofDumpOutFile         string
	pprofDebugFlag           int32
	pprofRequestGCBeforeDump bool
	pprofProfileDuration     time.Duration
)

//nolint:lll // setup of commands and flags might result in rather long lines
func init() {
	knownProfiles := pprof.Profiles()
	profileNames := make([]string, 0, len(knownProfiles))
	for idx := range knownProfiles {
		profileNames = append(profileNames, knownProfiles[idx].Name())
	}

	dumpProfileCmd.ValidArgs = profileNames
	dumpProfileCmd.Short = fmt.Sprintf("%s [%s]", dumpProfileCmd.Short, strings.Join(profileNames, ", "))
	dumpProfileCmd.Flags().Int32Var(&pprofDebugFlag, "debug", 0, "pprof debug flag - set e.g. to 1 to get legacy text format")
	dumpProfileCmd.Flags().BoolVar(&pprofRequestGCBeforeDump, "request-gc", false, "Request GC before dumping the profile - only applied for heap dumps")

	dumpCPUProfile.Flags().DurationVar(&pprofProfileDuration, "duration", defaultProfilingDuration, "Duration how long to profile CPU usage")
	dumpTraceCmd.Flags().DurationVar(&pprofProfileDuration, "duration", defaultProfilingDuration, "Duration how long to profile CPU usage")

	var defaultOutFile string
	if wd, err := os.Getwd(); err == nil {
		defaultOutFile = filepath.Join(wd, "profile.dump")
	}

	pprofCmd.PersistentFlags().StringVar(&pprofDumpOutFile, "out-file", defaultOutFile, "path where fetched pprof dump will be stored")
	pprofCmd.AddCommand(dumpProfileCmd, dumpCPUProfile, dumpTraceCmd)
}

func runDumpProfile(profile string) (err error) {
	profileClient := rpcv1.NewProfilingServiceClient(conn)

	req := &rpcv1.ProfileDumpRequest{ProfileName: profile, Debug: pprofDebugFlag, GcBeforeDump: pprofRequestGCBeforeDump}
	resp, err := profileClient.ProfileDump(cliApp.Context(), req)
	if err != nil {
		return err
	}

	return os.WriteFile(pprofDumpOutFile, resp.ProfileData, outFileMode)
}

func runDumpCPUProfile() error {
	profileClient := rpcv1.NewProfilingServiceClient(conn)
	req := &rpcv1.CPUProfileRequest{ProfileDuration: durationpb.New(pprofProfileDuration)}
	resp, err := profileClient.CPUProfile(cliApp.Context(), req)
	if err != nil {
		return err
	}

	return os.WriteFile(pprofDumpOutFile, resp.ProfileData, outFileMode)
}

func runDumpTrace() error {
	profileClient := rpcv1.NewProfilingServiceClient(conn)
	req := &rpcv1.TraceRequest{TraceDuration: durationpb.New(pprofProfileDuration)}
	resp, err := profileClient.Trace(cliApp.Context(), req)
	if err != nil {
		return err
	}

	return os.WriteFile(pprofDumpOutFile, resp.ProfileData, outFileMode)
}
