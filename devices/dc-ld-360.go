package devices

type DcLd360 struct {
	vid         uint16 // DeepCool Vendor ID
	pid         uint16 // LD-Series Product ID
	cpuTdpWatts int
}

func NewDcLd360() DcLd360 {
	return DcLd360{
		vid:         0x3633,
		pid:         0x000A,
		cpuTdpWatts: 170,
	}
}

func (dc *DcLd360) GetVid() uint16 {
	return dc.vid
}

func (dc *DcLd360) GetPid() uint16 {
	return dc.pid
}

func (dc *DcLd360) GetCpuTdpWatts() int {
	return dc.cpuTdpWatts
}
