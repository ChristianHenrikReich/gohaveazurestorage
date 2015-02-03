package gohaveazurestoragecommon

import "encoding/xml"

type StorageServiceProperties struct {
	XMLName       xml.Name `xml:"StorageServiceProperties"`
	Logging       LoggingProperties
	HourMetrics   HourMetricsProperties
	MinuteMetrics MinuteMetricsProperties
	Cors          CorsXMLDto `xml:",omitempty"`
}

type LoggingProperties struct {
	Version         string
	Read            bool
	Write           bool
	Delete          bool
	RetentionPolicy RetentionPolicy
}

type HourMetricsProperties struct {
	Version         string
	Enabled         bool
	IncludeAPIs     bool `xml:",omitempty"`
	RetentionPolicy RetentionPolicy
}

type MinuteMetricsProperties struct {
	Version         string
	Enabled         bool
	IncludeAPIs     bool `xml:",omitempty"`
	RetentionPolicy RetentionPolicy
}

type RetentionPolicy struct {
	Enabled bool
	Days    int `xml:",omitempty"`
}

type CorsXMLDto struct {
	//	CorsRule CorsRuleXMLDto `xml:",omitempty"`
}

type CorsRuleXMLDto struct {
	AllowedOrigins  string `xml:",omitempty"`
	AllowedMethods  string `xml:",omitempty"`
	MaxAgeInSeconds string `xml:",omitempty"`
	ExposedHeaders  string `xml:",omitempty"`
	AllowedHeaders  string `xml:",omitempty"`
}
