package tools

import (
	"log"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type DiskStatus struct {
	All  int64
	Used int64
	Free int64
}

func (ds DiskStatus) DiskInfo(diskPath string) (disk DiskStatus) {
	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		log.Panic(err)
	}
	defer syscall.FreeLibrary(kernel32)
	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")

	if err != nil {
		log.Panic(err)
	}

	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	_, _, _ = syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(diskPath))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)

	ds.All = lpTotalNumberOfBytes
	ds.Free = lpTotalNumberOfFreeBytes
	ds.Used = ds.All - ds.Free
	return ds
	// log.Printf("Available  %dmb", lpFreeBytesAvailable/1024/1024.0)
	// log.Printf("Total      %dmb", lpTotalNumberOfBytes/1024/1024.0)
	// log.Printf("Free       %dmb", lpTotalNumberOfFreeBytes/1024/1024.0)
}

func GetFree() int64 {
	var ds DiskStatus
	path, _ := os.Getwd()
	disPath := path[:strings.Index(path, "\\")]
	ds = ds.DiskInfo(disPath)
	return ds.Free
}
