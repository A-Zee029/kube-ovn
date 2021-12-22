package health_check

import (
	"fmt"
	"github.com/kubeovn/kube-ovn/pkg/ovs"
	"k8s.io/klog"
	"net"
	"os"
	"os/exec"
	"time"
)

func CmdMain() {

	daemonPid, err := exec.Command("cat", "/var/run/ovn/ovn-nbctl.pid").CombinedOutput()

	if err != nil {
		klog.Errorf("failed to get ovn-nbctl daemon pid, %s", err)
		os.Exit(1)
	}

	if err := os.Setenv("OVN_NB_DAEMON", fmt.Sprintf("/var/run/ovn/ovn-nbctl.%s.ctl", daemonPid)); err != nil {
		klog.Errorf("failed to set env OVN_NB_DAEMON, %v", err)
		os.Exit(1)
	}

	if err := ovs.CheckAlive(); err != nil {
		os.Exit(1)
	}

	conn, err := net.DialTimeout("tcp", "127.0.0.1:10660", 3*time.Second)
	if err != nil {
		klog.Errorf("failed to probe the socket, %s", err)
		os.Exit(1)
	}
	err = conn.Close()
	if err != nil {
		klog.Errorf("Unexpected error closing TCP probe socket: %v (%#v)", err, err)
	}
}
