package kernelmodule

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func IsContainerGPU() bool {
	cmd := exec.Command("sh", "-c", "lspci | grep -i NVIDIA")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	if strings.Contains(string(output), "NVIDIA") {
		return true
	} else {
		return false
	}
}

func isKernelModuleLoaded(moduleName string) (bool, error) {
	file, err := os.Open("/proc/modules")
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return false, err
		}
		if err == io.EOF {
			break
		}
		if strings.HasPrefix(line, moduleName) {
			return true, nil
		}
	}
	return false, nil
}

func CheckKernelMod() {
	log.Printf("check IB kernel module")

	var modInfo string

	var IBModule = []string{
		"rdma_ucm",
		"rdma_cm",
		"ib_ipoib",
		"mlx5_core",
		"mlx5_ib",
		"ib_uverbs",
		"ib_umad",
		"ib_cm",
		"ib_core",
		"mlxfw",
	}

	if IsContainerGPU() {
		IBModule = append(IBModule, "nvidia_peermem")
	}

	for i := 0; i < len(IBModule); i++ {
		loaded, err := isKernelModuleLoaded(IBModule[i])
		if err != nil {
			fmt.Println("Error :", err)
		} else {
			if !loaded {
				modInfo += IBModule[i] + ","
				cmd := exec.Command("sh", "-c", "modprobe "+IBModule[i])
				_, err := cmd.Output()
				if err != nil {
					log.Printf("fail fix load run modprobe %s", IBModule[i])
					return
				}

				log.Printf("success fix load run modprobe %s", IBModule[i])
			}
		}
	}
	if len(modInfo) == 0 {
		modInfo = "success load all modules"
	} else {
		modInfo = "fail load " + modInfo
	}

	log.Printf("%s", modInfo)
}
