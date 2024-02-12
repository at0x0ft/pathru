/*
Copyright © 2023 at0x0ft <26642966+at0x0ft@users.noreply.github.com>

*/
package cmd

import "fmt"	// 4debug
import (
	"os"
	"os/exec"
	"strings"
	"path/filepath"
	"github.com/spf13/cobra"
	// "github.com/at0x0ft/pathru/internal/pkg/schema"
	"github.com/at0x0ft/pathru/pkg/mount"
	"github.com/at0x0ft/pathru/pkg/parser"
	"github.com/at0x0ft/pathru/pkg/resolver"
)

const (
	SrcMountPointEnv = "CONTAINER_WORKSPACE_FOLDER"
	HostMountPointEnv = "LOCAL_WORKSPACE_FOLDER"
	DevContainerDirname = ".devcontainer"
	BaseShellServiceName = "base_shell"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pathru",
	Short: "Command pass-through helper with path conversion",
	Long: `pathru is a CLI command for help executing command in external container.
Usage: pathru <runtime service name> <execute command> -- [command arguments & options]`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := cmd.PersistentFlags().GetString("file")
		if err != nil {
			return err
		}
		files := strings.Split(file, ",")
		return execBody(files, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.museum.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("file", "f", "", "compose file path(s)")
}

// command body
func execBody(files []string, args []string) error {
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

	mergedComposeFiles := mergeComposeFiles(files)
	p := parser.ComposeYamlParser{mergedComposeFiles}

	var err error
	var mounts map[string]mount.BindMount
	mounts, err = p.Parse()
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

func mergeComposeFiles(files []string) string {
	return ""
}

func tryResolvingPath(arg string) (bool, string) {
	// [Warning] naive implementation
	absPath, err := filepath.Abs(arg)
	if err != nil {
		return false, ""
	}
	_, err = os.Stat(absPath)
	if err != nil {
		return false, ""
	}
	return true, absPath
}

func getDockerComposeFileAbsPathList() ([]string, error) {
	// devcontainerDirPath := filepath.Join(os.Getenv(SrcMountPointEnv), DevContainerDirname)
	// devcontainer, err := schema.LoadDevcontainer(devcontainerDirPath)
	// if err != nil {
		// return nil, err
		return nil, nil
	// }

	// var absPathList []string
	// for _, dockerComposeFileRelpath := range devcontainer.DockerComposeFile {
	//	 absPathList = append(absPathList, filepath.Join(devcontainerDirPath, dockerComposeFileRelpath))
	// }
	// return absPathList, nil
}

func serviceExists(serviceName string, dockerComposeFileList []string) error {
	// dockerCompose, err := schema.LoadMultipleDockerComposes(dockerComposeFileList)
	// if err != nil {
	//	 return err
	// }
	// // fmt.Println(dockerCompose)   // 4debug

	// for definedService, _ := range dockerCompose.Services {
	//	 if serviceName == definedService {
			return nil
	//	 }
	// }
	// return fmt.Errorf(
	//	 "[Error] service = '%s' not exists in docker-compose.yml files (%v) .",
	//	 serviceName,
	//	 strings.Join(dockerComposeFileList, ", "),
	// )
}

func getRuntimeMountPoints(serviceName string) (string, string, error) {
	dockerComposeFileList, err := getDockerComposeFileAbsPathList()
	if err != nil {
		return "", "", err
	}

	if err := serviceExists(serviceName, dockerComposeFileList); err != nil {
		return "", "", err
	}

	// os.Exit(1)  // 4debug
	// TODO: implement later

	// 3. get runtime container mount point path
	// 4. (host) source-path, (container) destination-path, error
	return "", "", nil
}

func convertPath(baseAbsPath string, runtimeServiceName string) (string, error) {
	runtimeSrcMountPoint, runtimeDstMountPoint, err := getRuntimeMountPoints(runtimeServiceName)
	if err != nil {
		return "", err
	}

	hostAbsPath, err := filepath.Rel(os.Getenv(SrcMountPointEnv), baseAbsPath)
	if err != nil {
		return "", err
	}
	runtimeDstRelPath, err := filepath.Rel(runtimeSrcMountPoint, hostAbsPath)
	if err != nil {
		return "", err
	}
	runtimeDstAbsPath := filepath.Join(runtimeDstMountPoint, runtimeDstRelPath)
	return runtimeDstAbsPath, nil
}

func convertPathIfFileExists(runtimeServiceName string, executeCommand string, args []string) ([]string, error) {
	result := []string {runtimeServiceName, executeCommand}
	var err error
	for _, arg := range args {
		isFilePath, absPath := tryResolvingPath(arg)
		if !isFilePath {
			continue
		}

		arg, err = convertPath(absPath, runtimeServiceName)
		if err != nil {
			return nil, err
		}
		result = append(result, arg)
	}
	return result, nil
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
	fmt.Println(cmd)	// 4debug
	return "", 0	// 4debug
}

func pathExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	return err == nil
}

func resolveArgs(args []string, runtimeServiceName string, mounts map[string]mount.BindMount) ([]string, error) {
	r := resolver.PathResolver{mounts}
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
