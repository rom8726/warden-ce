package domain

type SegmentName string

func (s SegmentName) String() string {
	return string(s)
}

const (
	SegmentNamePlatform       SegmentName = "platform"
	SegmentNameBrowserName    SegmentName = "browser_name"
	SegmentNameBrowserVersion SegmentName = "browser_version"
	SegmentNameOSName         SegmentName = "os_name"
	SegmentNameOSVersion      SegmentName = "os_version"
	SegmentNameDeviceArch     SegmentName = "device_arch"
	SegmentNameRuntimeName    SegmentName = "runtime_name"
	SegmentNameRuntimeVersion SegmentName = "runtime_version"
)
