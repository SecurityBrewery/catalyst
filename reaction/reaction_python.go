package reaction

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runPythonReaction(ctx context.Context, name, bootstrap, script, payload string) ([]byte, error) {
	tempDir, err := os.MkdirTemp("", "catalyst_action_"+name)
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(tempDir)

	if err := pythonSetup(ctx, tempDir); err != nil {
		return nil, err
	}

	if err := pythonRunBootstrap(ctx, tempDir, bootstrap); err != nil {
		return nil, err
	}

	return pythonRunScript(ctx, tempDir, script, payload)
}

func pythonSetup(ctx context.Context, tempDir string) error {
	pythonPath, err := findExec("python", "python3")
	if err != nil {
		return fmt.Errorf("python or python3 binary not found, %w", err)
	}

	// setup virtual environment
	return exec.CommandContext(ctx, pythonPath, "-m", "venv", tempDir+"/venv").Run()
}

func pythonRunBootstrap(ctx context.Context, tempDir, bootstrap string) error {
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

	return exec.CommandContext(ctx, pipPath, "install", "-r", bootstrapPath).Run()
}

func pythonRunScript(ctx context.Context, tempDir, script, payload string) ([]byte, error) {
	scriptPath := tempDir + "/script.py"

	if err := os.WriteFile(scriptPath, []byte(script), 0o600); err != nil {
		return nil, err
	}

	pythonPath := tempDir + "/venv/bin/python"

	return exec.CommandContext(ctx, pythonPath, scriptPath, payload).Output()
}

func findExec(name ...string) (string, error) {
	for _, n := range name {
		if p, err := exec.LookPath(n); err == nil {
			return p, nil
		}
	}

	return "", errors.New("no executable found")
}
