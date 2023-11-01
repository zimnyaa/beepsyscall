package main
import (
	"reflect"
	"syscall"
	"unsafe"
	"fmt"
)

type ProxyArgs struct {
	Addr    uintptr
	ArgsLen uintptr
	Args1   uintptr
	Args2   uintptr
	Args3   uintptr
	Args4   uintptr
	Args5   uintptr
	Args6   uintptr
	Args7   uintptr
	Args8   uintptr
	Args9   uintptr
	Args10  uintptr
}

func ProxyCall(Addr uintptr, Args ...uintptr) {
	Args0 := proxySetargs(Addr, Args...)
	ProxyCallWithStruct(Args0)
}

func ProxyCallWithStruct(Args0 uintptr) {
	tag := []byte{0x60, 0x70, 0x80, 0x90, 0x90, 0x90, 0x90}
	addr := reflect.ValueOf(proxyTag).Pointer()
	BaseAddr := addr
	addr = findTag(tag, BaseAddr)
	modkernel32 := syscall.MustLoadDLL("kernel32.dll")
	modkernelbase := syscall.MustLoadDLL("kernelbase.dll")
    procCreateThread := modkernel32.MustFindProc("CreateThread")
    procQueueUserApc := modkernel32.MustFindProc("QueueUserAPC")
    procResumeThread := modkernel32.MustFindProc("ResumeThread")

    fmt.Printf("starting search from %x\n", uintptr(modkernelbase.MustFindProc("Beep").Addr()))
    sleepgadget_addr := uintptr(modkernelbase.MustFindProc("Beep").Addr())
    for *(*uint64)(unsafe.Pointer(sleepgadget_addr)) != 0x8BD38B00000001BB {
    	sleepgadget_addr = sleepgadget_addr + 1
    }
    fmt.Printf("sleepgadget_addr: 0x%x\n", sleepgadget_addr)
    r1, _, _ := procCreateThread.Call(0, 0, sleepgadget_addr, Args0, 0x00000004, 0)

    procQueueUserApc.Call(addr, r1, Args0)
    procQueueUserApc.Call(modkernel32.MustFindProc("ExitThread").Addr(), r1, 0)
    procResumeThread.Call(r1)
	
	syscall.WaitForSingleObject(0xffffffffffffffff, 0x100)
}

func proxySetargs(Addr uintptr, Args ...uintptr) uintptr {
	newArgs := ProxyArgs{}
	newArgs.Addr = Addr
	if Args == nil {
		newArgs.ArgsLen = 0
		return uintptr(unsafe.Pointer(&newArgs))
	}
	if len(Args) > 10 {
		panic("Too much args")
	}
	len0 := len(Args)
	newArgs.ArgsLen = uintptr(len0)

	//使用反射遍历赋值
	pArgs := &newArgs
	value := reflect.ValueOf(pArgs).Elem()
	for i := 0; i < len0; i++ {
		if value.Field(i).CanSet() {
			ptr := unsafe.Pointer(value.Field(i + 2).UnsafeAddr())
			*(*uintptr)(ptr) = Args[i]
		}
	}
	return uintptr(unsafe.Pointer(&newArgs))
}

func findTag(egg []byte, startAddress uintptr) uintptr {
	var currentOffset = uintptr(0)
	currentAddress := startAddress
	for {
		currentOffset++
		currentAddress = startAddress + currentOffset
		if memcmp(unsafe.Pointer(&egg[0]), unsafe.Pointer(currentAddress), 7) == 0 {
			return currentAddress + 7
		}
	}
}

func memcmp(dest, src unsafe.Pointer, len uintptr) int {
	cnt := len >> 3
	var i uintptr = 0
	for i = 0; i < cnt; i++ {
		var pdest *uint64 = (*uint64)(unsafe.Pointer(uintptr(dest) + uintptr(8*i)))
		var psrc *uint64 = (*uint64)(unsafe.Pointer(uintptr(src) + uintptr(8*i)))
		switch {
		case *pdest < *psrc:
			return -1
		case *pdest > *psrc:
			return 1
		default:
		}
	}

	left := len & 7
	for i = 0; i < left; i++ {
		var pdest *uint8 = (*uint8)(unsafe.Pointer(uintptr(dest) + uintptr(8*cnt+i)))
		var psrc *uint8 = (*uint8)(unsafe.Pointer(uintptr(src) + uintptr(8*cnt+i)))
		switch {
		case *pdest < *psrc:
			return -1
		case *pdest > *psrc:
			return 1
		default:
		}
	}
	return 0
}

func proxyTag() {
	proxyC()
}

func proxyC()
