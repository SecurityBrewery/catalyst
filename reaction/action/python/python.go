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
	Requirements string `json:"requirements"`
	Script       string `json:"script"`

	env []string
}

func (a *Python) SetEnv(env []string) {
	a.env = env
}

func (a *Python) Run(ctx context.Context, payload string) ([]byte, error) {
	tempDir, err := os.MkdirTemp("", "catalyst_action")
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(tempDir)

	b, err := pythonSetup(ctx, tempDir)
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			b = append(b, ee.Stderr...)
		}

		return nil, fmt.Errorf("failed to setup python, %w: %s", err, string(b))
	}

	b, err = a.pythonInstallRequirements(ctx, tempDir)
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			b = append(b, ee.Stderr...)
		}

		return nil, fmt.Errorf("failed to run install requirements, %w: %s", err, string(b))
	}

	b, err = a.pythonRunScript(ctx, tempDir, payload)
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
	pythonPath, err := findExec("python3", "python")
	if err != nil {
		return nil, fmt.Errorf("python or python3 binary not found, %w", err)
	}

	// setup virtual environment
	return exec.CommandContext(ctx, pythonPath, "-m", "venv", tempDir+"/venv").Output()
}

func (a *Python) pythonInstallRequirements(ctx context.Context, tempDir string) ([]byte, error) {
	hasRequirements := len(strings.TrimSpace(a.Requirements)) > 0

	if !hasRequirements {
		return nil, nil
	}

	requirementsPath := tempDir + "/requirements.txt"

	if err := os.WriteFile(requirementsPath, []byte(a.Requirements), 0o600); err != nil {
		return nil, err
	}

	// install dependencies
	pipPath := tempDir + "/venv/bin/pip"

	return exec.CommandContext(ctx, pipPath, "install", "-r", requirementsPath).Output()
}

func (a *Python) pythonRunScript(ctx context.Context, tempDir, payload string) ([]byte, error) {
	scriptPath := tempDir + "/script.py"

	if err := os.WriteFile(scriptPath, []byte(a.Script), 0o600); err != nil {
		return nil, err
	}

	pythonPath := tempDir + "/venv/bin/python"

	cmd := exec.CommandContext(ctx, pythonPath, scriptPath, payload)

	cmd.Env = a.env

	return cmd.Output()
}

func findExec(name ...string) (string, error) {
	for _, n := range name {
		if p, err := exec.LookPath(n); err == nil {
			return p, nil
		}
	}

	return "", errors.New("no executable found")
}
