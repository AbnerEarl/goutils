package cmdc

import (
	"fmt"
	"os"
	"os/exec"
)

func GenHost(index uint64) string {
	base := 100000 + index
	return fmt.Sprintf("s%d.test.com", base)
}

func AddHost(ip, hostname string) bool {
	op := fmt.Sprintf("echo '%s %s' >> /etc/hosts", ip, hostname)
	cmd := exec.Command("sh", "-c", op)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func DelHost(hostname string) bool {
	op := fmt.Sprintf("sed -i '/%s/d' /etc/hosts", hostname)
	cmd := exec.Command("sh", "-c", op)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func RepHost(oldIp, newIp string) bool {
	op := fmt.Sprintf("sed -i 's/%s/%s/' /etc/hosts", oldIp, newIp)
	cmd := exec.Command("sh", "-c", op)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func RcoHost() bool {
	op := fmt.Sprintf("echo '127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4\n::1         localhost localhost.localdomain localhost6 localhost6.localdomain6' > /etc/hosts")
	cmd := exec.Command("sh", "-c", op)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}
