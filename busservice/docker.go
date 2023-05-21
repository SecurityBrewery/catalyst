package busservice

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/SecurityBrewery/catalyst/database"
)

func createContainer(ctx context.Context, image, script, data, network string) (string, string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", "", err
	}

	logs, err := pullImage(ctx, cli, image)
	if err != nil {
		return "", logs, err
	}

	config := &container.Config{
		Image:        image,
		Cmd:          []string{"/script", data},
		WorkingDir:   "/home",
		AttachStderr: true,
		AttachStdout: true,
	}
	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode(network),
	}
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, "")
	if err != nil {
		return "", logs, err
	}

	if err := copyFile(ctx, cli, "/script", script, resp.ID); err != nil {
		return "", logs, err
	}

	return resp.ID, logs, nil
}

func pullImage(ctx context.Context, cli *client.Client, image string) (string, error) {
	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, reader)

	return buf.String(), err
}

func copyFile(ctx context.Context, cli *client.Client, path string, contentString string, id string) error {
	tarBuf := &bytes.Buffer{}
	tw := tar.NewWriter(tarBuf)
	header := &tar.Header{Name: path, Mode: 0o755, Size: int64(len(contentString))}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := tw.Write([]byte(contentString)); err != nil {
		return err
	}

	if err := tw.Close(); err != nil {
		return err
	}

	return cli.CopyToContainer(ctx, id, "/", tarBuf, types.CopyToContainerOptions{})
}

func runDocker(ctx context.Context, jobID, containerID string, db *database.Database) (stdout []byte, stderr []byte, err error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, nil, err
	}

	defer func(cli *client.Client, ctx context.Context, containerID string, options types.ContainerRemoveOptions) {
		err := cli.ContainerRemove(ctx, containerID, options)
		if err != nil {
			log.Println(err)
		}
	}(cli, ctx, containerID, types.ContainerRemoveOptions{Force: true})

	if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return nil, nil, err
	}

	stderrBuf, err := streamStdErr(ctx, cli, jobID, containerID, db)
	if err != nil {
		return nil, nil, err
	}

	if err := waitForContainer(ctx, cli, containerID, stderrBuf); err != nil {
		return nil, nil, err
	}

	output, err := getStdOut(ctx, cli, containerID)
	if err != nil {
		log.Println(err)
	}

	return output.Bytes(), stderrBuf.Bytes(), nil
}

func streamStdErr(ctx context.Context, cli *client.Client, jobID, containerID string, db *database.Database) (*bytes.Buffer, error) {
	stderrBuf := &bytes.Buffer{}
	containerLogs, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStderr: true, Follow: true})
	if err != nil {
		return nil, err
	}
	go func() {
		err := scanLines(ctx, jobID, containerLogs, stderrBuf, db)
		if err != nil {
			log.Println(err)

			return
		}
		if err := containerLogs.Close(); err != nil {
			log.Println(err)

			return
		}
	}()

	return stderrBuf, nil
}

func scanLines(ctx context.Context, jobID string, input io.ReadCloser, output io.Writer, db *database.Database) error {
	r, w := io.Pipe()
	go func() {
		_, err := stdcopy.StdCopy(w, w, input)
		if err != nil {
			log.Println(err)

			return
		}
		if err := w.Close(); err != nil {
			log.Println(err)

			return
		}
	}()
	s := bufio.NewScanner(r)
	for s.Scan() {
		b := s.Bytes()
		_, _ = output.Write(b)
		_, _ = output.Write([]byte("\n"))

		if err := db.JobLogAppend(ctx, jobID, string(b)+"\n"); err != nil {
			log.Println(err)

			continue
		}
	}

	return s.Err()
}

func waitForContainer(ctx context.Context, cli *client.Client, containerID string, stderrBuf *bytes.Buffer) error {
	statusCh, errCh := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case exitStatus := <-statusCh:
		if exitStatus.StatusCode != 0 {
			return fmt.Errorf("container returned status code %d: stderr: %s", exitStatus.StatusCode, stderrBuf.String())
		}
	}

	return nil
}

func getStdOut(ctx context.Context, cli *client.Client, containerID string) (*bytes.Buffer, error) {
	output := &bytes.Buffer{}
	containerLogs, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStdout: true, Follow: true})
	if err != nil {
		return nil, err
	}
	defer containerLogs.Close()

	_, err = stdcopy.StdCopy(output, output, containerLogs)
	if err != nil {
		return nil, err
	}

	return output, nil
}
