//namespaces limits what we can see from inside the container
//control groups limits the resources

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

	syscall.Sethostname([]byte("container"))

	//change the root directory for this process
	syscall.Chroot("/home/aryanarora/dummy_fs")
	syscall.Chdir("/")

	//proc is a psudo-filesystem used between kernel and user to share information
	//For this process mount proc in the container_filesystem as a proc pseudo-filesystem
	//so that the kernel knows that it has to populate proc in the container_filesystem
	syscall.Mount("proc", "proc", "proc", 0, "")

	cg()
	
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	syscall.Unmount("/proc", 0)

}

func cg() {

	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	err := os.MkdirAll(filepath.Join(pids, "container_from_scratch"), 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	// Setting max process our container can create
	must(ioutil.WriteFile(filepath.Join(pids, "container_from_scratch/pids.max"), []byte("20"), 0700))
	// Removes the new cgroup after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, "container_from_scratch/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "container_from_scratch/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
