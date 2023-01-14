package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/udhos/equalfile"
)

func TestCopy(t *testing.T) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		require.NoError(t, err)
	}

	from := "./testdata/input.txt"
	to := f.Name()
	cmp := equalfile.New(nil, equalfile.Options{})

	defer os.Remove(to)

	t.Run("empty from", func(t *testing.T) {
		err := Copy("", to, 0, 0)
		require.Error(t, err)
	})

	t.Run("empty to", func(t *testing.T) {
		err = Copy(from, "", 0, 0)
		require.Error(t, err)
	})

	t.Run("wrong offset", func(t *testing.T) {
		err = Copy(from, to, -1, 0)
		require.Error(t, err)
	})

	t.Run("wrong limit", func(t *testing.T) {
		err = Copy(from, to, 0, -1)
		require.Error(t, err)
	})

	t.Run("copy file from to", func(t *testing.T) {
		err := Copy(from, to, 0, 0)
		require.NoError(t, err)
		equal, err := cmp.CompareFile("./testdata/out_offset0_limit0.txt", to)
		require.NoError(t, err)
		require.True(t, equal)
	})

	t.Run("copy no exist file", func(t *testing.T) {
		fakeFrom := "./testdata/input_fake.txt"
		err = Copy(fakeFrom, to, 0, 0)
		require.Error(t, err)
	})

	t.Run("copy set limit", func(t *testing.T) {
		err := Copy(from, to, 0, 10)
		require.NoError(t, err)
		equal, err := cmp.CompareFile("./testdata/out_offset0_limit10.txt", to)
		require.NoError(t, err)
		require.True(t, equal)
	})

	t.Run("copy set limit over file", func(t *testing.T) {
		err := Copy(from, to, 0, 10000)
		require.NoError(t, err)
		equal, err := cmp.CompareFile("./testdata/out_offset0_limit10000.txt", to)
		require.NoError(t, err)
		require.True(t, equal)
	})

	t.Run("copy set offset", func(t *testing.T) {
		err := Copy(from, to, 100, 1000)
		require.NoError(t, err)
		equal, err := cmp.CompareFile("./testdata/out_offset100_limit1000.txt", to)
		require.NoError(t, err)
		require.True(t, equal)
	})

	t.Run("copy set offset over file", func(t *testing.T) {
		err := Copy(from, to, 100000, 0)
		require.Error(t, err)
	})

	t.Run("copy not regular files", func(t *testing.T) {
		notRegularFile := "/dev/urandom"
		err = Copy(notRegularFile, to, 0, 0)
		require.Error(t, err)

		notRegularFile = "/dev/zero"
		err = Copy(notRegularFile, to, 0, 0)
		require.Error(t, err)
	})
}
