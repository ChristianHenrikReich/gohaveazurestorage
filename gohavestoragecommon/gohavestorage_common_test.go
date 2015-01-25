package gohavestoragecommon

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestStoragePropertiesXMLSerialization(t *testing.T) {
	expectedXML := "<StorageServiceProperties><Logging><Version>1.0</Version><Read>false</Read><Write>false</Write><Delete>false</Delete><RetentionPolicy><Enabled>false</Enabled></RetentionPolicy></Logging><HourMetrics><Version>1.0</Version><Enabled>false</Enabled><RetentionPolicy><Enabled>false</Enabled></RetentionPolicy></HourMetrics><MinuteMetrics><Version>1.0</Version><Enabled>false</Enabled><RetentionPolicy><Enabled>false</Enabled></RetentionPolicy></MinuteMetrics><Cors></Cors></StorageServiceProperties>"

	retentionPolicy := RetentionPolicy{Enabled: false, Days: 0}
	loggingProperties := LoggingProperties{Version: "1.0", Delete: false, Read: false, Write: false, RetentionPolicy: retentionPolicy}
	hourMetricsProperties := HourMetricsProperties{Version: "1.0", Enabled: false, IncludeAPIs: false, RetentionPolicy: retentionPolicy}
	minuteMetricsProperties := MinuteMetricsProperties{Version: "1.0", Enabled: false, IncludeAPIs: false, RetentionPolicy: retentionPolicy}
	storageServiceProperties := StorageServiceProperties{Logging: loggingProperties, HourMetrics: hourMetricsProperties, MinuteMetrics: minuteMetricsProperties}

	output, err := xml.MarshalIndent(storageServiceProperties, "", "")
	if err != nil {
		t.Error("XML Masrshalling error: %+v", err)
	}

	if string(output) != expectedXML {
		fmt.Printf("%s\n\nvs\n\n%s", expectedXML, string(output))
		t.Fail()
	}
}

func TestStorageStatsXMLSerialization(t *testing.T) {
	expectedXML := "<StorageServiceStats><GeoReplication><Status>bootstrap</Status><LastSyncTime></LastSyncTime></GeoReplication></StorageServiceStats>"

	storageServiceStats := &StorageServiceStats{GeoReplication: GeoReplication{Status: "bootstrap", LastSyncTime: ""}}

	output, err := xml.MarshalIndent(storageServiceStats, "", "")
	if err != nil {
		t.Error("XML Masrshalling error: %+v", err)
	}

	if string(output) != expectedXML {
		fmt.Printf("%s\n\nvs\n\n%s", expectedXML, string(output))
		t.Fail()
	}
}

func TestStorageAclXMLSerialization(t *testing.T) {
	expectedXML := "<SignedIdentifiers><SignedIdentifier><Id>b54df8ab0e2d52759110f48c8d0c19e2</Id><AccessPolicy><Start>2014-12-31</Start><Expiry>2114-12-31</Expiry><Permission>raud</Permission></AccessPolicy></SignedIdentifier></SignedIdentifiers>"

	accessPolicy := AccessPolicy{Start: "2014-12-31", Expiry: "2114-12-31", Permission: "raud"}
	signedIdentifier := SignedIdentifier{Id: "b54df8ab0e2d52759110f48c8d0c19e2", AccessPolicy: accessPolicy}
	signedIdentifiers := SignedIdentifiers{[]SignedIdentifier{signedIdentifier}}

	output, err := xml.MarshalIndent(signedIdentifiers, "", "")
	if err != nil {
		t.Error("XML Masrshalling error: %+v", err)
	}

	if string(output) != expectedXML {
		fmt.Printf("%s\n\nvs\n\n%s", expectedXML, string(output))
		t.Fail()
	}
}
