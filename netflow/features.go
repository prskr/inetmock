package netflow

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/features"
)

type EBPFFeature uint32

const (
	FeaturePerfEvent EBPFFeature = iota
	FeatureRingBuf
	FeatureBatchAPI
)

func (f EBPFFeature) String() string {
	switch f {
	case FeaturePerfEvent:
		return "PerfEventMap"
	case FeatureRingBuf:
		return "RingBufMap"
	case FeatureBatchAPI:
		return "MapBatchAPI"
	default:
		return fmt.Sprintf("Unknown value %d", f)
	}
}

var (
	featureChecks = map[EBPFFeature]func() error{
		FeaturePerfEvent: func() error {
			return features.HaveMapType(ebpf.PerfEventArray)
		},
		FeatureRingBuf: func() error {
			return features.HaveMapType(ebpf.RingBuf)
		},
		FeatureBatchAPI: func() error {
			m, err := ebpf.NewMap(&ebpf.MapSpec{
				Name:       "batch_api_test",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 1,
			})
			if err != nil {
				return err
			}

			_, err = m.BatchUpdate([]uint32{1}, []uint32{1}, nil)
			_ = m.Close()
			return err
		},
	}

	featureStates = make(map[EBPFFeature]bool)
	featureLock   sync.Mutex
)

func CheckForFeature(feat EBPFFeature) bool {
	featureLock.Lock()
	defer featureLock.Unlock()

	if state, ok := featureStates[feat]; ok {
		return state
	}

	if check, ok := featureChecks[feat]; !ok {
		panic(fmt.Errorf("check for unknown feature: %s", feat))
	} else if err := check(); err != nil {
		if errors.Is(err, ebpf.ErrNotSupported) {
			featureStates[feat] = false
		}
		return false
	} else {
		featureStates[feat] = true
		return true
	}
}
