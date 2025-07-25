package domain

import "time"

type ReleaseStats struct {
	ProjectID   ProjectID
	ReleaseID   ReleaseID
	Release     string
	GeneratedAt time.Time

	KnownIssuesTotal uint
	NewIssuesTotal   uint
	RegressionsTotal uint

	ResolvedInVersionTotal uint
	FixedNewInVersionTotal uint
	FixedOldInVersionTotal uint

	AvgFixTime    *time.Duration
	MedianFixTime *time.Duration
	P95FixTime    *time.Duration

	SeverityDistribution map[string]uint
	UsersAffected        uint
}

type ReleaseComparison struct {
	BaseRelease   ReleaseStats
	TargetRelease ReleaseStats
	Delta         map[string]uint
}

type ReleaseAnalyticsDetails struct {
	Release       Release
	Stats         ReleaseStats
	TopIssues     []Issue
	ByPlatform    map[string]uint
	ByBrowser     map[string]uint
	ByOS          map[string]uint
	ByDeviceArch  map[string]uint
	ByRuntimeName map[string]uint
}

type UserSegmentKey string

const (
	SegmentPlatformAndroid UserSegmentKey = "android"
	SegmentPlatformIOS     UserSegmentKey = "ios"
	SegmentPlatformWeb     UserSegmentKey = "web"
	SegmentPlatformWindows UserSegmentKey = "windows"
	SegmentPlatformLinux   UserSegmentKey = "linux"

	SegmentBrowserChrome  UserSegmentKey = "Chrome"
	SegmentBrowserFirefox UserSegmentKey = "Firefox"
	SegmentBrowserSafari  UserSegmentKey = "Safari"
	SegmentBrowserEdge    UserSegmentKey = "Edge"

	SegmentOSWindows UserSegmentKey = "Windows"
	SegmentOSMacOS   UserSegmentKey = "macOS"
	SegmentOSLinux   UserSegmentKey = "Linux"
	SegmentOSAndroid UserSegmentKey = "Android"
	SegmentOSIOS     UserSegmentKey = "iOS"
)

type UserSegmentsAggregation map[UserSegmentKey]uint

type UserSegmentsAnalytics struct {
	Platform    UserSegmentsAggregation
	Browser     UserSegmentsAggregation
	OS          UserSegmentsAggregation
	DeviceArch  UserSegmentsAggregation
	RuntimeName UserSegmentsAggregation
}
