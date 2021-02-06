package ctl_test

import (
	"github.com/eloylp/aton/components/ctl"
)

const (
	LeastUtilized       = "LEAST_UTILIZED_DETECTOR"
	OneThirdUtilized    = "ONE_THIRD_UTILIZED_DETECTOR"
	MidUtilized         = "MID_UTILIZED_DETECTOR"
	AverageUtilized     = "AVERAGE_UTILIZED_DETECTOR"
	FullUtilized        = "FULL_UTILIZED_DETECTOR"
	FullCPUUtilized     = "FULL_CPU_UTILIZED_DETECTOR"
	FullMemoryUtilized  = "FULL_MEMORY_UTILIZED_DETECTOR"
	FullNetworkUtilized = "FULL_NETWORK_UTILIZED_DETECTOR"
)

func LeastUtilizedDetector() *ctl.Detector {
	return &ctl.Detector{
		UUID: LeastUtilized,
		Status: &ctl.Status{
			Capturers: nil,
			System: ctl.System{
				CPUCount: 2,
				Network: ctl.Network{
					RxBytesSec: 1000,
					TxBytesSec: 200,
				},
				LoadAverage: ctl.LoadAverage{
					Avg1:  0.1,
					Avg5:  0.1,
					Avg15: 0.1,
				},
				Memory: ctl.Memory{
					UsedMemoryBytes:  32e6,
					TotalMemoryBytes: 256e6,
				},
			},
		},
	}
}

func OneThirdUtilizedDetector() *ctl.Detector {
	return &ctl.Detector{
		UUID: OneThirdUtilized,
		Status: &ctl.Status{
			Capturers: nil,
			System: ctl.System{
				CPUCount: 2,
				Network: ctl.Network{
					RxBytesSec: 30e6,
					TxBytesSec: 10e6,
				},
				LoadAverage: ctl.LoadAverage{
					Avg1:  0.66,
					Avg5:  0.66,
					Avg15: 0.66,
				},
				Memory: ctl.Memory{
					UsedMemoryBytes:  84e6,
					TotalMemoryBytes: 256e6,
				},
			},
		},
	}
}

func MidUtilizedDetector() *ctl.Detector {
	return &ctl.Detector{
		UUID: MidUtilized,
		Status: &ctl.Status{
			Capturers: nil,
			System: ctl.System{
				CPUCount: 2,
				Network: ctl.Network{
					RxBytesSec: 50e6,
					TxBytesSec: 20e6,
				},
				LoadAverage: ctl.LoadAverage{
					Avg1:  1,
					Avg5:  1,
					Avg15: 1,
				},
				Memory: ctl.Memory{
					UsedMemoryBytes:  128e6,
					TotalMemoryBytes: 256e6,
				},
			},
		},
	}
}

func AverageUtilizedDetector() *ctl.Detector {
	return &ctl.Detector{
		UUID: AverageUtilized,
		Status: &ctl.Status{
			Capturers: nil,
			System: ctl.System{
				CPUCount: 2,
				Network: ctl.Network{
					RxBytesSec: 75e6,
					TxBytesSec: 30e6,
				},
				LoadAverage: ctl.LoadAverage{
					Avg1:  1.5,
					Avg5:  1.5,
					Avg15: 1.5,
				},
				Memory: ctl.Memory{
					UsedMemoryBytes:  192e6,
					TotalMemoryBytes: 256e6,
				},
			},
		},
	}
}

func FullUtilizedDetector() *ctl.Detector {
	return &ctl.Detector{
		UUID: FullUtilized,
		Status: &ctl.Status{
			Capturers: nil,
			System: ctl.System{
				CPUCount: 2,
				Network: ctl.Network{
					RxBytesSec: 100e6,
					TxBytesSec: 100e6,
				},
				LoadAverage: ctl.LoadAverage{
					Avg1:  2,
					Avg5:  2,
					Avg15: 2,
				},
				Memory: ctl.Memory{
					UsedMemoryBytes:  256e6,
					TotalMemoryBytes: 256e6,
				},
			},
		},
	}
}

func FullCPUUtilizedDetector() *ctl.Detector {
	return &ctl.Detector{
		UUID: FullCPUUtilized,
		Status: &ctl.Status{
			Capturers: nil,
			System: ctl.System{
				CPUCount: 2,
				Network: ctl.Network{
					RxBytesSec: 75e6,
					TxBytesSec: 30e6,
				},
				LoadAverage: ctl.LoadAverage{
					Avg1:  2,
					Avg5:  2,
					Avg15: 2,
				},
				Memory: ctl.Memory{
					UsedMemoryBytes:  192e6,
					TotalMemoryBytes: 256e6,
				},
			},
		},
	}
}

func FullMemoryUtilizedDetector() *ctl.Detector {
	return &ctl.Detector{
		UUID: FullMemoryUtilized,
		Status: &ctl.Status{
			Capturers: nil,
			System: ctl.System{
				CPUCount: 2,
				Network: ctl.Network{
					RxBytesSec: 75e6,
					TxBytesSec: 30e6,
				},
				LoadAverage: ctl.LoadAverage{
					Avg1:  1.5,
					Avg5:  1.5,
					Avg15: 1.5,
				},
				Memory: ctl.Memory{
					UsedMemoryBytes:  255e6,
					TotalMemoryBytes: 256e6,
				},
			},
		},
	}
}

func FullNetworkUtilizedDetector() *ctl.Detector {
	return &ctl.Detector{
		UUID: FullNetworkUtilized,
		Status: &ctl.Status{
			Capturers: nil,
			System: ctl.System{
				CPUCount: 2,
				Network: ctl.Network{
					RxBytesSec: 99.9e6,
					TxBytesSec: 50e6,
				},
				LoadAverage: ctl.LoadAverage{
					Avg1:  1.5,
					Avg5:  1.5,
					Avg15: 1.5,
				},
				Memory: ctl.Memory{
					UsedMemoryBytes:  192e6,
					TotalMemoryBytes: 256e6,
				},
			},
		},
	}
}
