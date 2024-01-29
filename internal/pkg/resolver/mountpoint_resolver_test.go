package resolver

import (
    "testing"
    "os"
)

func TestMain(t *testing.M) {
    code := t.Run()
    os.Exit(code)
}

type ResolveSuccessTestCase struct {
    srcMountPoint string
    dstMountPoint string
    path string
    expected string
}

func providerTestResolveSuccess() map[string]ResolveSuccessTestCase {
    return map[string]ResolveSuccessTestCase{
        "normal case": {
            "/workspace",
            "/home/testuser",
            "/workspace/hoge/fuga",
            "/home/testuser/hoge/fuga",
        },
        "multiple subpath matched case": {
            "/workspace/hoge",
            "/home/testuser",
            "/workspace/hoge/fuga",
            "/home/testuser/fuga",
        },
        "mountPoints are directories case": {
            "/workspace/",
            "/home/testuser/",
            "/workspace/hoge/fuga/piyo",
            "/home/testuser/hoge/fuga/piyo",
        },
        "directory path converts to slash trailed path": {
            "/workspace",
            "/home/testuser",
            "/workspace/hoge/",
            "/home/testuser/hoge",
        },
    }
}

func TestResolveSuccess(t * testing.T) {
    cases := providerTestResolveSuccess()
    for name, c := range cases {
        t.Run(name, func(t *testing.T) {
            c := c
            t.Parallel()
            m := MountPointResolver{c.srcMountPoint, c.dstMountPoint}
            actual, err := m.Resolve(c.path)
            if err != nil {
                t.Error(err)
            }

            if c.expected != actual {
                t.Errorf(
                    "resolved actual path does not match [expected = \"%s\", actual = \"%s\"]",
                    c.expected,
                    actual,
                )
            }
        })
    }
}

type ResolveFailTestCase struct {
    srcMountPoint string
    dstMountPoint string
    path string
    expectedMsg string
}

func providerTestResolveFail() map[string]ResolveFailTestCase {
    return map[string]ResolveFailTestCase{
        "path not containing srcMountPoint": {
            "/workspace",
            "/home/testuser",
            "/home/hoge/fuga",
            "path cannot reach srcMountPoint [path: \"/home/hoge/fuga\", srcMountPoint: \"/workspace\"]",
        },
    }
}

func TestResolveFail(t * testing.T) {
    cases := providerTestResolveFail()
    for name, c := range cases {
        t.Run(name, func(t *testing.T) {
            c := c
            t.Parallel()
            m := MountPointResolver{c.srcMountPoint, c.dstMountPoint}
            result, err := m.Resolve(c.path)
            if err == nil {
                t.Errorf("no error thrown [result = \"%s\"]", result)
                return
            }

            actual := err.Error()
            if c.expectedMsg != actual {
                t.Errorf(
                    "error message does not match [expected = \"%s\", actual = \"%s\"]",
                    c.expectedMsg,
                    actual,
                )
            }
        })
    }
}
