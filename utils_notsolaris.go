// +build !solaris

package zfs

import (
	"strings"
)

// List of ZFS properties to retrieve from zfs list command on a non-Solaris platform
var dsPropList = []string{
	"name",
	"available",
	"compressratio",
	"defer_destroy",
	"mounted",
	"origin",
	"referenced",
	"type",
	"used",
	"usedbychildren",
	"usedbydataset",
	"usedbyrefreservation",
	"usedbysnapshots",
	"userrefs",
	"aclinherit",
	"aclmode",
	"atime",
	"canmount",
	"casesensitivity",
	"checksum",
	"compression",
	"copies",
	"dedup",
	"devices",
	"exec",
	"logbias",
	"mlslabel",
	"mountpoint",
	"nbmand",
	"normalization",
	"primarycache",
	"quota",
	"readonly",
	"recordsize",
	"refquota",
	"refreservation",
	"reservation",
	"secondarycache",
	"setuid",
	"sharenfs",
	"sharesmb",
	"snapdir",
	"utf8only",
	"version",
	"volblocksize",
	"volsize",
	"vscan",
	"xattr",
	"zoned",
}

func indexOfDSProp(element string) int {
	for k, v := range dsPropList {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

var dsPropListOptions = strings.Join(dsPropList, ",")

// List of Zpool properties to retrieve from zpool list command on a non-Solaris platform
var zpoolPropList = []string{
	"name",
	"health",
	"allocated",
	"size",
	"free",
	"capacity",
	"altroot",
	"guid",
	"version",
	"bootfs",
	"delegation",
	"autoreplace",
	"cachefile",
	"failmode",
	"listsnapshots",
	"autoexpand",
	"dedupditto",
	"dedupratio",
	"ashift",
}
var zpoolPropListOptions = strings.Join(zpoolPropList, ",")
var zpoolArgs = []string{"get", zpoolPropListOptions}
