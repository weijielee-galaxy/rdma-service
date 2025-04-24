package acs

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func getPciDevices() ([]string, error) {
	cmd := exec.Command("lspci", "-d", "*:*:*")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")
	var devices []string
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			devices = append(devices, fields[0]) // Get the first column which is the BDF
		}
	}

	return devices, nil
}

// supportsACS checks if a PCI device supports ACS
func supportsACS(bdf string) bool {
	// Use setpci to check for ACS support
	cmd := exec.Command("setpci", "-v", "-s", bdf, "ECAP_ACS+0x6.w")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return err == nil
}

func getCurrentACS(bdf string) error {
	// Verify if ACS has been disabled
	cmd := exec.Command("setpci", "-v", "-s", bdf, "ECAP_ACS+0x6.w")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	newVal := strings.TrimSpace(out.String())
	fields := strings.Fields(newVal)

	if fields[len(fields)-1] != "0000" {
		return fmt.Errorf("failed to disable ACS on %s, value:%s", bdf, newVal)
	}
	log.Printf("check ACS on bdf: %s, msg: %s", bdf, newVal)
	return nil
}

func disableACS(bdf string) error {
	// check disabled
	if getCurrentACS(bdf) == nil {
		return nil
	}

	// Disable ACS by setting it to 0000
	cmd := exec.Command("setpci", "-v", "-s", bdf, "ECAP_ACS+0x6.w=0000")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Verify if ACS has been disabled
	cmd = exec.Command("setpci", "-v", "-s", bdf, "ECAP_ACS+0x6.w")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	newVal := strings.TrimSpace(out.String())
	fields := strings.Fields(newVal)

	if fields[len(fields)-1] != "0000" {
		return fmt.Errorf("failed to disable ACS on %s, value:%s", bdf, newVal)
	}
	log.Printf("disable ACS on bdf: %s, msg: %s", bdf, newVal)
	return nil
}

func CheckACS() {
	log.Printf("check ACS is disabled")

	pciDevices, err := getPciDevices()
	if err != nil {
		log.Fatalf("Failed to get PCI devices: %v", err)
	}

	for _, bdf := range pciDevices {
		if !supportsACS(bdf) {
			continue
		}

		if err := disableACS(bdf); err != nil {
			log.Printf("Error disabling ACS on %s: %v", bdf, err)
		}
	}
}
