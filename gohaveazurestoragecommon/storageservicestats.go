package gohaveazurestoragecommon

import "encoding/xml"

type StorageServiceStats struct {
	XMLName        xml.Name `xml:"StorageServiceStats"`
	GeoReplication GeoReplication
}

type GeoReplication struct {
	Status       string
	LastSyncTime string
}
