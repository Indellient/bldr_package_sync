package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

// ${_HAB} pkg upload --url "${local_depot_url}" --channel "${channel}" "${hartfile_path}" && break

func packageUpload(target BldrApi, fileName string, channel string) {
	// env := []string{"HAB_BLDR_URL=" + target.Url, "HAB_AUTH_TOKEN=" + target.AuthToken, "HAB_CACHE_KEY_PATH=" + dir, "SSL_CERT_FILE=/usr/local/etc/openssl/cert.pem"}
	env := []string{"HAB_BLDR_URL=" + target.Url, "HAB_AUTH_TOKEN=" + target.AuthToken, "SSL_CERT_FILE=/usr/local/etc/openssl/cert.pem"}

	cmd := fmt.Sprintf("pkg upload --channel \"%s\" %s", channel, fileName)

	log.Debug("Running `hab " + cmd + "`")

	runHabCommandEnv(cmd, env)
}

func importPublicKey(target BldrApi, dir string, fileName string) {
	// env := []string{"HAB_BLDR_URL=" + target.Url, "HAB_AUTH_TOKEN=" + target.AuthToken, "HAB_CACHE_KEY_PATH=" + dir, "SSL_CERT_FILE=/usr/local/etc/openssl/cert.pem"}
	env := []string{"HAB_BLDR_URL=" + target.Url, "HAB_AUTH_TOKEN=" + target.AuthToken, "SSL_CERT_FILE=/usr/local/etc/openssl/cert.pem"}

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
	cmd.Env = append(cmd.Env, habEnv...)

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
