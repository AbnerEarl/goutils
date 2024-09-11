/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/9/9 15:27
 * @desc: about the role of class.
 */

package utils

import (
	"net"
	"regexp"
)

var (
	Ipv4Regex = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	Ipv6Regex = `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:([0-9a-fA-F]{1,4}|:)|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]).){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]).){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
)

func IsIpv4(ip string) bool {
	res, _ := regexp.MatchString(ip, Ipv4Regex)
	return res
}

func IsIpv6(ip string) bool {
	res, _ := regexp.MatchString(ip, Ipv6Regex)
	return res
}

func SubnetMatch(subnet string) (string, string, error) {
	ip, ipv4Net, err := net.ParseCIDR(subnet)
	if err != nil {
		return "", "", err
	}
	return ipv4Net.String(), ip.String(), err
}

func CheckIP(ip string) string {
	result := net.ParseIP(ip)
	if result != nil {
		return result.String()
	}
	return ""
}
