package connection

import (
	"os"
	"testing"
	"time"
)

func ImageMerge_Integration_test(t *testing.T) {
	go func() {
		time.Sleep(2 * time.Second)
		t.Error("frame not received within time limit")
		os.Exit(1)
	}()

	var _ = SimpleServer("8097", RedFrame())
	var _ = SimpleServer("8098", BlueFrame())

}
