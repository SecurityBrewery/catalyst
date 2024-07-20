package python

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Python struct {
	Bootstrap string `json:"bootstrap"`
	Script    string `json:"script"`
}

func (a *Python) Run(ctx context.Context, payload string) ([]byte, error) {
	tempDir, err := os.MkdirTemp("", "catalyst_action")
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(tempDir)

	if b, err := pythonSetup(ctx, tempDir); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			b = append(b, ee.Stderr...)
		}

		return nil, fmt.Errorf("failed to setup python, %w: %s", err, string(b))
	}

	if b, err := pythonRunBootstrap(ctx, tempDir, a.Bootstrap); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			b = append(b, ee.Stderr...)
		}

		return nil, fmt.Errorf("failed to run bootstrap, %w: %s", err, string(b))
	}

	b, err := pythonRunScript(ctx, tempDir, a.Script, payload)
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			b = append(b, ee.Stderr...)
		}

		return nil, fmt.Errorf("failed to run script, %w: %s", err, string(b))
	}

	return b, nil
}

func pythonSetup(ctx context.Context, tempDir string) ([]byte, error) {
	pythonPath, err := findExec("python", "python3")
	if err != nil {
		return nil, fmt.Errorf("python or python3 binary not found, %w", err)
	}

	// setup virtual environment
	return exec.CommandContext(ctx, pythonPath, "-m", "venv", tempDir+"/venv").Output()
}

func pythonRunBootstrap(ctx context.Context, tempDir, bootstrap string) ([]byte, error) {
	hasBootstrap := len(strings.TrimSpace(bootstrap)) > 0

	if !hasBootstrap {
		return nil, nil
	}

	bootstrapPath := tempDir + "/requirements.txt"

	if err := os.WriteFile(bootstrapPath, []byte(bootstrap), 0o600); err != nil {
		return nil, err
	}

	// install dependencies
	pipPath := tempDir + "/venv/bin/pip"

	return exec.CommandContext(ctx, pipPath, "install", "-r", bootstrapPath).Output()
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
