//namespaces limits what we can see from inside the container
//control groups limits the resources
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bad command")
	}

}
func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	// "/proc/self/exe" runs itself(this program)
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// namespacing(limit what a process can see  ) hostname, pid and mount
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS}

	cmd.Run()
}
func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	syscall.Sethostname([]byte("container"))

	//change the root directory for this process
	syscall.Chroot("/home/aryanarora/dummy_fs")
	syscall.Chdir("/")

	//proc is a psudo-filesystem used between kernel and user to share information
	//For this process mount proc in the container_filesystem as a proc pseudo-filesystem
	//so that the kernel knows that it has to populate proc in the container_filesystem
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd.Run()

	syscall.Unmount("/proc", 0)

}
