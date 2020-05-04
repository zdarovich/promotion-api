package log

import (
	"bytes"
	"errors"
	"net"
	"runtime"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func getFields() logrus.Fields {

	f := make(logrus.Fields)
	f["timestamp"] = time.Now().Format(time.RFC3339)
	f["meta"] = struct {
		Version string `json:"version"`
	}{Version: "1.0"}
	f["app_name"] = "APIBuilder"
	f["host"] = host
	f["host_ip"] = ip
	f["process_id"] = pid
	f["logger"] = "internal/log"
	f["thread"] = getGID()

	return f
}

func withContext() *logrus.Entry {
	return logrus.WithFields(getFields())
}

//get goroutine id from rutime.Stack(). Stack formats a stack trace of the calling goroutine into buf and returns the number of bytes written to buf. If all is true, Stack formats stack traces of all other goroutines into buf after the trace for the current goroutine.
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0
	}
	return n
}

// Get preferred outbound ip of this machine
func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
