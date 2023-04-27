package machine

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os/exec"
)

func LinuxMachineOnlyID() string {
	result := ""
	err, productName := ShellScript("cat /sys/class/dmi/id/product_name")
	if err == nil {
		result += "###" + productName
	}
	err, productUuid := ShellScript("cat /sys/class/dmi/id/product_uuid")
	if err == nil {
		result += "###" + productUuid
	}
	err, productSerial := ShellScript("cat /sys/class/dmi/id/board_serial")
	if err == nil {
		result += "###" + productSerial
	}
	err, systemSerial := ShellScript("dmidecode -s system-serial-number")
	if err == nil {
		result += "###" + systemSerial
	}
	return MD5(result)
}

func LinuxMachineMustID() string {
	result := ""
	err, productName := ShellScript("cat /sys/class/dmi/id/product_name")
	if err == nil {
		result += "###" + productName
	}
	err, productUuid := ShellScript("cat /sys/class/dmi/id/product_uuid")
	if err == nil {
		result += "###" + productUuid
	}
	err, productSerial := ShellScript("cat /sys/class/dmi/id/board_serial")
	if err == nil {
		result += "###" + productSerial
	}
	err, systemSerial := ShellScript("dmidecode -s system-serial-number")
	if err == nil {
		result += "###" + systemSerial
	}
	err, netMac := ShellScript("ip a | grep -A 1  \"$(ls /sys/class/net/ | grep -v \"`ls /sys/devices/virtual/net/`\")\" | grep ether | awk '{print $2}' | tr '\\n' ', ' | sed 's/,$//'")
	if err == nil {
		result += "###" + netMac
	}
	err, diskUuid := ShellScript("ls /dev/disk/by-uuid | tr '\\n' ',' | sed 's/,$//'")
	if err == nil {
		result += "###" + diskUuid
	}
	return MD5(result)
}

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func ShellScript(command string) (error, string) {
	result := ""
	cmd := exec.Command("sh", "-c", command)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return err, result
	}
	result = string(stdout.Bytes())
	return nil, result
}
