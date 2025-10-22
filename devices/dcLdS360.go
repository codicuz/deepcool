package devices

type DcLdS360 struct {
	vid uint16 // Vendor ID
	pid uint16 // Product ID
}

func NewDcLdS360() *DcLdS360 {
	return &DcLdS360{
		vid: 0x3633,
		pid: 0x000A,
	}
}

func (dc *DcLdS360) GetVid() uint16 {
	return dc.vid
}

func (dc *DcLdS360) GetPid() uint16 {
	return dc.pid
}
