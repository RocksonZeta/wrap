package netutil_test

import (
	"fmt"
	"testing"

	"github.com/RocksonZeta/wrap/utils/netutil"
)

func TestGetIp(t *testing.T) {
	ips, err := netutil.LocalIPv4s()
	fmt.Println(ips, err)
}
