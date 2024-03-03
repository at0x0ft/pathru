package mount

import "testing"

type ConvertTargetToSourceSuccessTestCase struct {
	source   string
	target   string
	path     string
	expected string
}

func providerTestConvertTargetToSourceSuccess() map[string]ConvertTargetToSourceSuccessTestCase {
	return map[string]ConvertTargetToSourceSuccessTestCase{
		"normal case": {
			"/home/testuser",
			"/workspace",
			"/workspace/hoge/fuga",
			"/home/testuser/hoge/fuga",
		},
		"multiple subpath matched case": {
			"/home/testuser",
			"/workspace/hoge",
			"/workspace/hoge/fuga",
			"/home/testuser/fuga",
		},
		"mountPoints are directories case": {
			"/home/testuser/",
			"/workspace/",
			"/workspace/hoge/fuga/piyo",
			"/home/testuser/hoge/fuga/piyo",
		},
		"directory path converts to slash trailed path": {
			"/home/testuser",
			"/workspace",
			"/workspace/hoge/",
			"/home/testuser/hoge",
		},
	}
}

func TestConvertTargetToSourceSuccess(t *testing.T) {
	cases := providerTestConvertTargetToSourceSuccess()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			bm := BindMount{c.source, c.target}
			actual, err := bm.ConvertTargetToSource(c.path)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}

			if c.expected != actual {
				t.Errorf(
					"resolved actual path does not match [expected = \"%s\", actual = \"%s\"]",
					c.expected,
					actual,
				)
				t.FailNow()
			}
		})
	}
}

type ConvertTargetToSourceFailTestCase struct {
	source      string
	target      string
	path        string
	expectedMsg string
}

func providerTestConvertTargetToSourceFail() map[string]ConvertTargetToSourceFailTestCase {
	return map[string]ConvertTargetToSourceFailTestCase{
		"path not containing mount.source": {
			"/home/testuser",
			"/workspace",
			"/home/hoge/fuga",
			"given path cannot reach to its mount base path [given: \"/home/hoge/fuga\", base: \"/workspace\"]",
		},
	}
}

func TestResolveTargetToSourceFail(t *testing.T) {
	cases := providerTestConvertTargetToSourceFail()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			bm := BindMount{c.source, c.target}
			result, err := bm.ConvertTargetToSource(c.path)
			if err == nil {
				t.Errorf("no error thrown [result = \"%s\"]", result)
				t.FailNow()
			}

			actual := err.Error()
			if c.expectedMsg != actual {
				t.Errorf(
					"error message does not match [expected = \"%s\", actual = \"%s\"]",
					c.expectedMsg,
					actual,
				)
				t.FailNow()
			}
		})
	}
}
