package util

import (
	"net"
	"strconv"
	"strings"
)

var localAddr = "unknown"

func init() {
	ifaces, err := net.Interfaces()
	if err != nil {
		return
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
			return
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
			localAddr = ip.String()
			return
		}
	}
}

// LexicalCompare takes two string, `lhs` and `rhs`, to compare them. If `lhs` < `rhs`, it will return `-idx`,
// else if `lhs` > `rhs`, it will return `idx`. If the two string are equal, it will return 0.
// `idx` here is the first words occurs not equal and where we can decide the order, note the
// `idx` will be 1-indexed, to distinguish from the result when the strings are equal.
func LexicalCompare(lhs, rhs string) int {
	lhsLen := len(lhs)
	rhsLen := len(rhs)

	length := lhsLen
	if lhsLen > rhsLen {
		length = rhsLen
	}

	for idx := 0; idx < length; idx++ {
		if lhs[idx] < rhs[idx] {
			return -idx - 1
		} else if lhs[idx] > rhs[idx] {
			return idx + 1
		}
	}

	if lhsLen == rhsLen {
		return 0
	} else if lhsLen < rhsLen {
		return -lhsLen
	}

	return rhsLen
}

// ParseInt64 will parse a string into a int64, just a strconv.PraseInt wrapper.
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// GetIP will returns the local external ip by iterate all net interfaces.
// If no ip is find, localAddr will be "unknown" by default.
func GetIP() string {
	return localAddr
}

// QuoteJoin will returns a joined string while all elements are quoted.
func QuoteJoin(s []string, sep string) string {
	var quoted []string

	for _, item := range s {
		quoted = append(quoted, strconv.Quote(item))
	}

	return strings.Join(quoted, sep)
}
