package ctl_test

import (
	"github.com/eloylp/aton/components/ctl"
)

const (
	LeastUtilized       = "LEAST_UTILIZED_NODE"
	OneThirdUtilized    = "ONE_THIRD_UTILIZED_NODE"
	MidUtilized         = "MID_UTILIZED_NODE"
	AverageUtilized     = "AVERAGE_UTILIZED_NODE"
	FullUtilized        = "FULL_UTILIZED_NODE"
	FullCPUUtilized     = "FULL_CPU_UTILIZED_NODE"
	FullMemoryUtilized  = "FULL_MEMORY_UTILIZED_NODE"
	FullNetworkUtilized = "FULL_NETWORK_UTILIZED_NODE"
)

func LeastUtilizedNode() *ctl.Node {
	return &ctl.Node{
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

func OneThirdUtilizedNode() *ctl.Node {
	return &ctl.Node{
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

func MidUtilizedNode() *ctl.Node {
	return &ctl.Node{
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

func AverageUtilizedNode() *ctl.Node {
	return &ctl.Node{
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

func FullUtilizedNode() *ctl.Node {
	return &ctl.Node{
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

func FullCPUUtilizedNode() *ctl.Node {
	return &ctl.Node{
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

func FullMemoryUtilizedNode() *ctl.Node {
	return &ctl.Node{
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

func FullNetworkUtilizedNode() *ctl.Node {
	return &ctl.Node{
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
