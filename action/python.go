package action

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runPythonAction(name, bootstrap, script, payload string) ([]byte, error) {
	tempDir, err := os.MkdirTemp("", "catalyst_action_"+name)
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(tempDir)

	if err := pythonSetup(tempDir); err != nil {
		return nil, err
	}

	if err := pythonRunBootstrap(tempDir, bootstrap); err != nil {
		return nil, err
	}

	return pythonRunScript(tempDir, script, payload)
}

func pythonSetup(tempDir string) error {
	pythonPath, err := findExec("python", "python3")
	if err != nil {
		return fmt.Errorf("python or python3 binary not found, %w", err)
	}

	// setup virtual environment
	return exec.Command(pythonPath, "-m", "venv", tempDir+"/venv").Run()
}

func pythonRunBootstrap(tempDir, bootstrap string) error {
	hasBootstrap := len(strings.TrimSpace(bootstrap)) > 0

	if !hasBootstrap {
		return nil
	}

	bootstrapPath := tempDir + "/requirements.txt"

	if err := os.WriteFile(bootstrapPath, []byte(bootstrap), 0o600); err != nil {
		return err
	}

	// install dependencies
	pipPath := tempDir + "/venv/bin/pip"

	return exec.Command(pipPath, "install", "-r", bootstrapPath).Run()
}

func pythonRunScript(tempDir, script, payload string) ([]byte, error) {
	scriptPath := tempDir + "/script.py"

	if err := os.WriteFile(scriptPath, []byte(script), 0o600); err != nil {
		return nil, err
	}

	pythonPath := tempDir + "/venv/bin/python"

	return exec.Command(pythonPath, scriptPath, payload).Output()
}

func findExec(name ...string) (string, error) {
	for _, n := range name {
		if p, err := exec.LookPath(n); err == nil {
			return p, nil
		}
	}

	return "", errors.New("no executable found")
}
