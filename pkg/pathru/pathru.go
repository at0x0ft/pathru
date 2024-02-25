package pathru

import "fmt" // 4debug
import (
	"github.com/at0x0ft/pathru/pkg/compose"
	"github.com/at0x0ft/pathru/pkg/mount"
	"github.com/at0x0ft/pathru/pkg/resolver"
	docker_compose "github.com/docker/compose/v2/cmd/compose"
	"os"
	"os/exec"
	"path/filepath"
)

// TODO: delete constants
const (
	SrcMountPointEnv     = "CONTAINER_WORKSPACE_FOLDER"
	HostMountPointEnv    = "LOCAL_WORKSPACE_FOLDER"
	DevContainerDirname  = ".devcontainer"
	BaseShellServiceName = "base_shell"
)

func Process(opts *docker_compose.ProjectOptions, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("[Error] Not enough arguments are given!")
	}

	// TODO: runtime service name validation
	// 1. make sure 'CONTAINER_WORKSPACE_FOLDER' is set.
	// note: should validate in later?
	runtimeServiceName := args[0]
	// executeCommand := args[1]
	// args = args[2:]
	args = args[1:]

	var err error
	var mounts map[string]mount.BindMount
	mounts, err = (&compose.ComposeParser{}).Parse([]string{})
	if err != nil {
		return err
	}

	args, err = resolveArgs(args, runtimeServiceName, mounts)
	if err != nil {
		return err
	}

	output, exitCode := execDockerCompose(args)
	if exitCode != 0 {
		fmt.Println("[Error] docker-compose run failed.")
		fmt.Println(output)
		os.Exit(exitCode)
	}
	fmt.Printf("%s\n", output)
	return nil
}

func execDockerCompose(args []string) (string, int) {
	args = append([]string{"run", "--rm"}, args...)
	cmd := exec.Command("docker-compose", args...)
	// TODO: split stdout & stderr
	// out, err := cmd.CombinedOutput()
	// result := string(out)
	// if err != nil {
	// 	result += err.Error()
	// }
	// return result, cmd.ProcessState.ExitCode()
	fmt.Println(cmd) // 4debug
	return "", 0     // 4debug
}

func pathExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	return err == nil
}

func resolveArgs(args []string, runtimeServiceName string, mounts map[string]mount.BindMount) ([]string, error) {
	r := resolver.PathResolver{Mounts: mounts}
	res := make([]string, len(args))
	for _, arg := range args {
		if pathExists(arg) {
			p, err := r.Resolve(arg, BaseShellServiceName, runtimeServiceName)
			if err != nil {
				return nil, err
			}
			res = append(res, p)
		} else {
			res = append(res, arg)
		}
	}
	return res, nil
}
