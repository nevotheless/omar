package bootc

import "os/exec"

func Switch(image string) error {
	cmd := exec.Command("bootc", "switch", image)
	return cmd.Run()
}

func Upgrade() error {
	cmd := exec.Command("bootc", "upgrade")
	return cmd.Run()
}

func Rollback() error {
	cmd := exec.Command("bootc", "rollback")
	return cmd.Run()
}

func Install(disk string) error {
	cmd := exec.Command("bootc", "install", "--target-disk", disk)
	return cmd.Run()
}

func Status() (string, error) {
	out, err := exec.Command("bootc", "status").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
