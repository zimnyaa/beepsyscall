# a fork of `timwhitez/proxycall` that uses Beep and QueueUserAPC
```
this is a fork of timwhitez's implementation of proxycall. it is an unstable,
untested implementation, which is basically a practical joke that illustrates how even
kernelbase.dll!Beep can be used for malware.
the implementation is based off of https://0xdarkvortex.dev/hiding-in-plainsight/

instead of TpAllocWork, this fork creates a suspended thread that will enter an 
alertable state by using a gadget in the Beep function in kernelbase.dll that calls 

SleepEx(..., TRUE).

it then schedules two APCs to that thread: 
- a call to the trampoline that extracts syscall arguments from the passed struct;
- and a call to ExitThread.

the only changes made here are in proxysyscall.go!ProxyCallWithStruct and 
unpacking the struct from the first argument instead of the second in
asm_x64.s (mov rbx,rdx -> mov rbx,rcx)

full credit goes to the original authors.
```
```
main.go contains a usage example.
```