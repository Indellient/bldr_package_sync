package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func packageUpload(target BldrApi, fileName string, channel string) {
	env := []string{"HAB_BLDR_URL=" + target.Url, "HAB_AUTH_TOKEN=" + target.AuthToken}

	cmd := fmt.Sprintf("pkg upload --channel \"%s\" %s", channel, fileName)

	log.Debug("Running `hab " + cmd + "`")

	runHabCommandEnv(cmd, env)
}

func packagePromote(target BldrApi, pkgName string, channel string) {
	env := []string{"HAB_BLDR_URL=" + target.Url, "HAB_AUTH_TOKEN=" + target.AuthToken}

	cmd := fmt.Sprintf("pkg promote \"%s\" %s", pkgName, channel)

	log.Debug("Running `hab " + cmd + "`")

	runHabCommandEnv(cmd, env)
}

func importPublicKey(target BldrApi, dir string, fileName string) {
	env := []string{"HAB_BLDR_URL=" + target.Url, "HAB_AUTH_TOKEN=" + target.AuthToken}

	cmd := fmt.Sprintf("origin key upload --pubfile \"%s\"", fileName)

	log.Debug("Running `hab " + cmd + "`")

	runHabCommandEnv(cmd, env)
}

// Run a habitat command given a hab environment variables and a directory to be executed from.
func runHabCommand(command string) {
	command = "hab " + command
	cmd := exec.Command("/bin/bash", "-c", command)
	path := fmt.Sprintf("PATH=%s", os.Getenv("PATH"))
	cmd.Env = append(cmd.Env, path)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			log.Info(scanner.Text())
		}
	}()

	scannerErr := bufio.NewScanner(stderr)
	go func() {
		for scannerErr.Scan() {
			log.Error(scannerErr.Text())
		}
	}()

	cmd.Wait()
}

// Run a habitat command given a hab environment variables and a directory to be executed from.
func runHabCommandEnv(command string, habEnv []string) {
	command = "hab " + command
	cmd := exec.Command("/bin/bash", "-c", command)
	path := fmt.Sprintf("PATH=%s", os.Getenv("PATH"))
	cmd.Env = append(cmd.Env, path)

	habEnv = append(habEnv, config.Env...)
	cmd.Env = append(cmd.Env, habEnv...)

	log.Debug(fmt.Sprintf("Running with Hab Environment Variables %s", habEnv))

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			log.Info(scanner.Text())
		}
	}()

	scannerErr := bufio.NewScanner(stderr)
	go func() {
		for scannerErr.Scan() {
			log.Error(scannerErr.Text())
		}
	}()

	cmd.Wait()
}

// Run a habitat command given a hab environment variables and a directory to be executed from.
func runHabCommandFromDirectory(command string, habEnv []string, dir string) {
	command = "hab " + command
	cmd := exec.Command("/bin/bash", "-c", command)
	path := fmt.Sprintf("PATH=%s", os.Getenv("PATH"))
	cmd.Env = append(cmd.Env, path)
	cmd.Dir = dir
	cmd.Env = append(cmd.Env, habEnv...)
	stdout, _ := cmd.StdoutPipe()
	// stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			log.Info(scanner.Text())
		}
	}()

	cmd.Wait()
}
