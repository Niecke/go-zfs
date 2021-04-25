package zfs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/go-units"
	"github.com/google/uuid"
)

type command struct {
	Command string
	Stdin   io.Reader
	Stdout  io.Writer
}

func (c *command) Run(arg ...string) ([][]string, error) {

	cmd := exec.Command(c.Command, arg...)

	var stdout, stderr bytes.Buffer

	if c.Stdout == nil {
		cmd.Stdout = &stdout
	} else {
		cmd.Stdout = c.Stdout
	}

	if c.Stdin != nil {
		cmd.Stdin = c.Stdin

	}
	cmd.Stderr = &stderr

	id := uuid.New().String()
	joinedArgs := strings.Join(cmd.Args, " ")

	logger.Log([]string{"ID:" + id, "START", joinedArgs})
	err := cmd.Run()
	logger.Log([]string{"ID:" + id, "FINISH"})

	if err != nil {
		return nil, &Error{
			Err:    err,
			Debug:  strings.Join([]string{cmd.Path, joinedArgs[1:]}, " "),
			Stderr: stderr.String(),
		}
	}

	// assume if you passed in something for stdout, that you know what to do with it
	if c.Stdout != nil {
		return nil, nil
	}

	lines := strings.Split(stdout.String(), "\n")

	//last line is always blank
	lines = lines[0 : len(lines)-1]
	output := make([][]string, len(lines))

	for i, l := range lines {
		output[i] = strings.Fields(l)
	}

	return output, nil
}

func setString(field *string, value string) {
	v := ""
	if value != "-" {
		v = value
	}
	*field = v
}

func setUint(field *uint64, value string) error {
	var v uint64
	if value != "-" {
		var err error
		v, err = strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
	}
	*field = v
	return nil
}

func setStorageValue(field *uint64, value string) error {
	var v int64
	if value == "-" {
		v = 0
	} else {
		var err error
		v, err = units.FromHumanSize(value)
		if err != nil {
			return err
		}
	}
	*field = uint64(v)
	return nil
}

func setBool(field *bool, value string) error {
	if value == "on" || value == "yes" {
		*field = true
	} else if value == "off" || value == "no" || value == "-" {
		*field = false
	} else {
		return errors.New("Value neither on/yes nor off/no")
	}
	return nil
}

func (ds *Dataset) parseLine(line []string) error {
	var err error

	if len(line) != len(dsPropList) {
		return errors.New("Output does not match what is expected on this platform")
	}
	setString(&ds.Name, line[indexOfDSProp("name")])
	if err = setStorageValue(&ds.Available, line[indexOfDSProp("available")]); err != nil {
		return err
	}
	if ds.CompressRatio, err = strconv.ParseFloat(line[indexOfDSProp("compressratio")][:len(line[indexOfDSProp("compressratio")])-1], 64); err != nil {
		return err
	}
	if err = setBool(&ds.DeferDestroy, line[indexOfDSProp("defer_destroy")]); err != nil {
		return err
	}
	if err = setBool(&ds.Mounted, line[indexOfDSProp("mounted")]); err != nil {
		return err
	}
	setString(&ds.Origin, line[indexOfDSProp("origin")])
	if err = setStorageValue(&ds.Referenced, line[indexOfDSProp("referenced")]); err != nil {
		return err
	}
	setString(&ds.Type, line[indexOfDSProp("type")])
	if err = setStorageValue(&ds.Used, line[indexOfDSProp("used")]); err != nil {
		return err
	}
	if err = setStorageValue(&ds.UsedByChildren, line[indexOfDSProp("usedbychildren")]); err != nil {
		return err
	}
	if err = setStorageValue(&ds.UsedByDataset, line[indexOfDSProp("usedbydataset")]); err != nil {
		return err
	}
	if err = setStorageValue(&ds.UsedByRefReservation, line[indexOfDSProp("usedbyrefreservation")]); err != nil {
		return err
	}
	if err = setStorageValue(&ds.UsedBysnapshots, line[indexOfDSProp("usedbysnapshots")]); err != nil {
		return err
	}
	if err = setUint(&ds.UserRefs, line[indexOfDSProp("userrefs")]); err != nil {
		return err
	}
	setString(&ds.Aclinherit, line[indexOfDSProp("aclinherit")])
	setString(&ds.AclMode, line[indexOfDSProp("aclmode")])
	if err = setBool(&ds.Atime, line[indexOfDSProp("atime")]); err != nil {
		return err
	}
	setString(&ds.CanMount, line[indexOfDSProp("canmount")])
	setString(&ds.CaseSensitivity, line[indexOfDSProp("casesensitivity")])
	setString(&ds.Checksum, line[indexOfDSProp("checksum")])
	setString(&ds.Compression, line[indexOfDSProp("compression")])
	setString(&ds.Copies, line[indexOfDSProp("copies")])
	setString(&ds.Dedup, line[indexOfDSProp("dedup")])
	if err = setBool(&ds.Devices, line[indexOfDSProp("devices")]); err != nil {
		return err
	}
	if err = setBool(&ds.Exec, line[indexOfDSProp("exec")]); err != nil {
		return err
	}
	setString(&ds.Logbias, line[indexOfDSProp("logbias")])
	setString(&ds.Mlslabel, line[indexOfDSProp("mlslabel")])
	setString(&ds.Mountpoint, line[indexOfDSProp("mountpoint")])
	if err = setBool(&ds.Nbmand, line[indexOfDSProp("nbmand")]); err != nil {
		return err
	}
	setString(&ds.Normalization, line[indexOfDSProp("normalization")])
	setString(&ds.Primarycache, line[indexOfDSProp("primarycache")])
	setString(&ds.Quota, line[indexOfDSProp("quota")])
	if err = setBool(&ds.Readonly, line[indexOfDSProp("readonly")]); err != nil {
		return err
	}
	setString(&ds.Recordsize, line[indexOfDSProp("recordsize")])
	setString(&ds.RefQuota, line[indexOfDSProp("refquota")])
	setString(&ds.RefReservation, line[indexOfDSProp("refreservation")])
	setString(&ds.Reservation, line[indexOfDSProp("reservation")])
	setString(&ds.SecondaryCache, line[indexOfDSProp("secondarycache")])
	if err = setBool(&ds.Setuid, line[indexOfDSProp("setuid")]); err != nil {
		return err
	}
	setString(&ds.Sharenfs, line[indexOfDSProp("sharenfs")])
	setString(&ds.Sharesmb, line[indexOfDSProp("sharesmb")])
	setString(&ds.Snapdir, line[indexOfDSProp("snapdir")])
	if err = setBool(&ds.UTF8only, line[indexOfDSProp("utf8only")]); err != nil {
		return err
	}
	setString(&ds.Version, line[indexOfDSProp("version")])
	setString(&ds.VolBlockSize, line[indexOfDSProp("volblocksize")])
	setString(&ds.VolSize, line[indexOfDSProp("volsize")])
	if err = setBool(&ds.Vscan, line[indexOfDSProp("vscan")]); err != nil {
		return err
	}
	if err = setBool(&ds.Xattr, line[indexOfDSProp("xattr")]); err != nil {
		return err
	}
	if err = setBool(&ds.Zoned, line[indexOfDSProp("zoned")]); err != nil {
		return err
	}
	return nil
}

/*
 * from zfs diff`s escape function:
 *
 * Prints a file name out a character at a time.  If the character is
 * not in the range of what we consider "printable" ASCII, display it
 * as an escaped 3-digit octal value.  ASCII values less than a space
 * are all control characters and we declare the upper end as the
 * DELete character.  This also is the last 7-bit ASCII character.
 * We choose to treat all 8-bit ASCII as not printable for this
 * application.
 */
func unescapeFilepath(path string) (string, error) {
	buf := make([]byte, 0, len(path))
	llen := len(path)
	for i := 0; i < llen; {
		if path[i] == '\\' {
			if llen < i+4 {
				return "", fmt.Errorf("Invalid octal code: too short")
			}
			octalCode := path[(i + 1):(i + 4)]
			val, err := strconv.ParseUint(octalCode, 8, 8)
			if err != nil {
				return "", fmt.Errorf("Invalid octal code: %v", err)
			}
			buf = append(buf, byte(val))
			i += 4
		} else {
			buf = append(buf, path[i])
			i++
		}
	}
	return string(buf), nil
}

var changeTypeMap = map[string]ChangeType{
	"-": Removed,
	"+": Created,
	"M": Modified,
	"R": Renamed,
}
var inodeTypeMap = map[string]InodeType{
	"B": BlockDevice,
	"C": CharacterDevice,
	"/": Directory,
	">": Door,
	"|": NamedPipe,
	"@": SymbolicLink,
	"P": EventPort,
	"=": Socket,
	"F": File,
}

// matches (+1) or (-1)
var referenceCountRegex = regexp.MustCompile("\\(([+-]\\d+?)\\)")

func parseReferenceCount(field string) (int, error) {
	matches := referenceCountRegex.FindStringSubmatch(field)
	if matches == nil {
		return 0, fmt.Errorf("Regexp does not match")
	}
	return strconv.Atoi(matches[1])
}

func parseInodeChange(line []string) (*InodeChange, error) {
	llen := len(line)
	if llen < 1 {
		return nil, fmt.Errorf("Empty line passed")
	}

	changeType := changeTypeMap[line[0]]
	if changeType == 0 {
		return nil, fmt.Errorf("Unknown change type '%s'", line[0])
	}

	switch changeType {
	case Renamed:
		if llen != 4 {
			return nil, fmt.Errorf("Mismatching number of fields: expect 4, got: %d", llen)
		}
	case Modified:
		if llen != 4 && llen != 3 {
			return nil, fmt.Errorf("Mismatching number of fields: expect 3..4, got: %d", llen)
		}
	default:
		if llen != 3 {
			return nil, fmt.Errorf("Mismatching number of fields: expect 3, got: %d", llen)
		}
	}

	inodeType := inodeTypeMap[line[1]]
	if inodeType == 0 {
		return nil, fmt.Errorf("Unknown inode type '%s'", line[1])
	}

	path, err := unescapeFilepath(line[2])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse filename: %v", err)
	}

	var newPath string
	var referenceCount int
	switch changeType {
	case Renamed:
		newPath, err = unescapeFilepath(line[3])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse filename: %v", err)
		}
	case Modified:
		if llen == 4 {
			referenceCount, err = parseReferenceCount(line[3])
			if err != nil {
				return nil, fmt.Errorf("Failed to parse reference count: %v", err)
			}
		}
	default:
		newPath = ""
	}

	return &InodeChange{
		Change:               changeType,
		Type:                 inodeType,
		Path:                 path,
		NewPath:              newPath,
		ReferenceCountChange: referenceCount,
	}, nil
}

// example input
//M       /       /testpool/bar/
//+       F       /testpool/bar/hello.txt
//M       /       /testpool/bar/hello.txt (+1)
//M       /       /testpool/bar/hello-hardlink
func parseInodeChanges(lines [][]string) ([]*InodeChange, error) {
	changes := make([]*InodeChange, len(lines))

	for i, line := range lines {
		c, err := parseInodeChange(line)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse line %d of zfs diff: %v, got: '%s'", i, err, line)
		}
		changes[i] = c
	}
	return changes, nil
}

func listByType(t, filter string) ([]*Dataset, error) {
	args := []string{"list", "-rH", "-t", t, "-o", dsPropListOptions}

	if filter != "" {
		args = append(args, filter)
	}
	out, err := zfs(args...)
	if err != nil {
		return nil, err
	}

	var datasets []*Dataset

	name := ""
	var ds *Dataset
	for _, line := range out {
		if name != line[0] {
			name = line[0]
			ds = &Dataset{Name: name}
			datasets = append(datasets, ds)
		}
		if err := ds.parseLine(line); err != nil {
			return nil, err
		}
	}

	return datasets, nil
}

func propsSlice(properties map[string]string) []string {
	args := make([]string, 0, len(properties)*3)
	for k, v := range properties {
		args = append(args, "-o")
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}
	return args
}

func (z *Zpool) parseLine(line []string) error {
	prop := line[1]
	val := line[2]

	var err error

	switch prop {
	case "name":
		setString(&z.Name, val)
	case "health":
		setString(&z.Health, val)
	case "allocated":
		err = setStorageValue(&z.Allocated, val)
	case "size":
		err = setStorageValue(&z.Size, val)
	case "free":
		err = setStorageValue(&z.Free, val)
	case "dedupratio":
		// Trim trailing "x" before parsing float64
		z.DedupRatio, err = strconv.ParseFloat(val[:len(val)-1], 64)
	case "capacity":
		err = setUint(&z.Capacity, val[:len(val)-1])
	case "altroot":
		setString(&z.Altroot, val)
	case "guid":
		setString(&z.Guid, val)
	case "version":
		err = setUint(&z.Version, val)
	case "bootfs":
		setString(&z.BootFS, val)
	case "delegation":
		err = setBool(&z.Delegation, val)
	case "autoreplace":
		err = setBool(&z.Autoreplace, val)
	case "cachefile":
		setString(&z.Cachefile, val)
	case "failmode":
		setString(&z.Failmode, val)
	case "listsnapshots":
		err = setBool(&z.ListSnapshots, val)
	case "autoexpand":
		err = setBool(&z.Autoexpand, val)
	case "dedupditto":
		setString(&z.Dedupditto, val)
	case "ashift":
		setString(&z.Ashift, val)
	}
	return err
}
