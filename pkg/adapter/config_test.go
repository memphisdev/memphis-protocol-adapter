package adapter_test

import (
	"testing"

	adapter "github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"
)

type TestConf struct {
	FIRSTSTRING string `mapstructure:"FIRSTSTRING"`
	SECONDINT   int    `mapstructure:"SECONDINT"`
}

func TestPrepare(t *testing.T) {

	confFolderPath := "./_conf_test/"

	cfact := adapter.ConfigFactory(confFolderPath)

	expected := TestConf{"1.0.3", 300}

	var tf TestConf

	err := cfact("example", &tf)

	if err != nil {
		t.Errorf("unmarshal error %v", err)
	}

	if tf.FIRSTSTRING != expected.FIRSTSTRING {
		t.Errorf("Expected %s Actual %s", expected.FIRSTSTRING, tf.FIRSTSTRING)
	}

	if tf.SECONDINT != expected.SECONDINT {
		t.Errorf("Expected %d Actual %d", expected.SECONDINT, tf.SECONDINT)
	}
}
