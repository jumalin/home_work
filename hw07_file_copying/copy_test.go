package main

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

/*
* offset больше, чем размер файла - невалидная ситуация;
* limit больше, чем размер файла - валидная ситуация, копируется исходный файл до его EOF;
* программа может НЕ обрабатывать файлы, у которых неизвестна длина (например, /dev/urandom);
 */

var readPath = "testdata/out_offset0_limit10.txt"

func getWritePath() string {
	return fmt.Sprintf("/tmp/result%s.txt", uuid.New())
}

func TestCopy(t *testing.T) {
	// Place your code here.
	err := Copy(readPath, getWritePath(), 0, 0)
	require.NoError(t, err)
}

func TestFileLessThanOffset(t *testing.T) {
	err := Copy(readPath, getWritePath(), 10, 0)
	require.ErrorContains(t, err, "Offset exceedes size of the file")
}

func TestFileLessThanLimit(t *testing.T) {
	err := Copy(readPath, getWritePath(), 0, 10)
	require.ErrorContains(t, err, "Offset exceedes size of the file")
}

func TestFileWithRandomSize(t *testing.T) {
	err := Copy(readPath, getWritePath(), 0, 0)
	require.ErrorContains(t, err, "Offset exceedes size of the file")
}
