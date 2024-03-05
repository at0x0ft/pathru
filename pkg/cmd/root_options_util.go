package cmd

import "github.com/spf13/pflag"

type OptionData[T int | string | []string] struct {
	value T
	name  string
	fs    *pflag.FlagSet
}

func CreateIntPersistentOptionData(
	fs *pflag.FlagSet,
	defaultValue int,
	name,
	shorthand,
	usage string,
) *OptionData[int] {
	o := &OptionData[int]{name: name, fs: fs}
	o.fs.IntVarP(&o.value, o.name, shorthand, defaultValue, usage)
	return o
}

func CreateIntOptionData(
	fs *pflag.FlagSet,
	defaultValue int,
	name,
	usage string,
) *OptionData[int] {
	o := &OptionData[int]{name: name, fs: fs}
	o.fs.IntVar(&o.value, o.name, defaultValue, usage)
	return o
}

func CreateStringPersistentOptionData(
	fs *pflag.FlagSet,
	defaultValue string,
	name,
	shorthand,
	usage string,
) *OptionData[string] {
	o := &OptionData[string]{name: name, fs: fs}
	o.fs.StringVarP(&o.value, o.name, shorthand, defaultValue, usage)
	return o
}

func CreateStringOptionData(
	fs *pflag.FlagSet,
	defaultValue string,
	name,
	usage string,
) *OptionData[string] {
	o := &OptionData[string]{name: name, fs: fs}
	o.fs.StringVar(&o.value, o.name, defaultValue, usage)
	return o
}

func CreateStringArrayPersistentOptionData(
	fs *pflag.FlagSet,
	defaultValue []string,
	name,
	shorthand,
	usage string,
) *OptionData[[]string] {
	o := &OptionData[[]string]{name: name, fs: fs}
	o.fs.StringArrayVarP(&o.value, o.name, shorthand, defaultValue, usage)
	return o
}

func CreateStringArrayOptionData(
	fs *pflag.FlagSet,
	defaultValue []string,
	name,
	usage string,
) *OptionData[[]string] {
	o := &OptionData[[]string]{name: name, fs: fs}
	o.fs.StringArrayVar(&o.value, o.name, defaultValue, usage)
	return o
}

func (o *OptionData[T]) Value() T {
	return o.value
}

func (o *OptionData[T]) IsSet() bool {
	f := o.fs.Lookup(o.name)
	if f == nil {
		return false
	}
	return f.Changed
}

func stringArrayEquals(ls1 []string, ls2 []string) bool {
	if ls1 == nil && ls2 == nil {
		return true
	} else if ls1 == nil || ls2 == nil {
		return false
	}
	if len(ls1) != len(ls2) {
		return false
	}
	for i, e1 := range ls1 {
		if e1 != ls2[i] {
			return false
		}
	}
	return true
}
