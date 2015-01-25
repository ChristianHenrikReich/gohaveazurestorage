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

		t.Fail()
	}
}
