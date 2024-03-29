package selfupdate

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"
)

var (
	ErrCannotDecompressFile        = errors.New("failed to decompress")
	ErrExecutableNotFoundInArchive = errors.New("executable not found")
)

var (
	fileTypes = []struct {
		ext        string
		decompress func(src io.Reader, cmd, os, arch string) (io.Reader, error)
	}{
		{".zip", unzip},
		{".tar.gz", untar},
		{".tgz", untar},
		{".gzip", gunzip},
		{".gz", gunzip},
		{".tar.xz", untarxz},
		{".xz", unxz},
		{".bz2", unbz2},
	}
)

// DecompressCommand decompresses the given source. Archive and compression format is
// automatically detected from 'url' parameter, which represents the URL of asset,
// or simply a filename (with an extension).
// This returns a reader for the decompressed command given by 'cmd'. '.zip',
// '.tar.gz', '.tar.xz', '.tgz', '.gz', '.bz2' and '.xz' are supported.
//
// These wrapped errors can be returned:
//   - ErrCannotDecompressFile
//   - ErrExecutableNotFoundInArchive
func DecompressCommand(src io.ReadCloser, url, cmd, os, arch string) (io.Reader, error) {
	for _, fileType := range fileTypes {
		if strings.HasSuffix(url, fileType.ext) {
			return fileType.decompress(src, cmd, os, arch)
		}
	}
	return src, nil
}

func unzip(src io.Reader, cmd, os, arch string) (io.Reader, error) {
	// Zip format requires its file size for Decompressing.
	// So we need to read the HTTP response into a buffer at first.
	buf, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("%w zip file: %v", ErrCannotDecompressFile, err)
	}

	r := bytes.NewReader(buf)
	z, err := zip.NewReader(r, r.Size())
	if err != nil {
		return nil, fmt.Errorf("%w zip file: %s", ErrCannotDecompressFile, err)
	}

	for _, file := range z.File {
		_, name := filepath.Split(file.Name)
		if !file.FileInfo().IsDir() && matchExecutableName(cmd, os, arch, name) {
			return file.Open()
		}
	}

	return nil, fmt.Errorf("%w in zip file: %q", ErrExecutableNotFoundInArchive, cmd)
}

func untar(src io.Reader, cmd, os, arch string) (io.Reader, error) {
	gz, err := gzip.NewReader(src)
	if err != nil {
		return nil, fmt.Errorf("%w tar.gz file: %s", ErrCannotDecompressFile, err)
	}

	return unarchiveTar(gz, cmd, os, arch)
}

func gunzip(src io.Reader, cmd, os, arch string) (io.Reader, error) {
	r, err := gzip.NewReader(src)
	if err != nil {
		return nil, fmt.Errorf("%w gzip file: %s", ErrCannotDecompressFile, err)
	}

	name := r.Header.Name
	if !matchExecutableName(cmd, os, arch, name) {
		return nil, fmt.Errorf("%w: expected %q but found %q", ErrExecutableNotFoundInArchive, cmd, name)
	}

	return r, nil
}

func untarxz(src io.Reader, cmd, os, arch string) (io.Reader, error) {
	xzip, err := xz.NewReader(src)
	if err != nil {
		return nil, fmt.Errorf("%w tar.xz file: %s", ErrCannotDecompressFile, err)
	}

	return unarchiveTar(xzip, cmd, os, arch)
}

func unxz(src io.Reader, cmd, os, arch string) (io.Reader, error) {
	xzip, err := xz.NewReader(src)
	if err != nil {
		return nil, fmt.Errorf("%w xzip file: %s", ErrCannotDecompressFile, err)
	}
	return xzip, nil
}

func unbz2(src io.Reader, cmd, os, arch string) (io.Reader, error) {
	bz2 := bzip2.NewReader(src)
	return bz2, nil
}

func matchExecutableName(cmd, os, arch, target string) bool {
	if cmd == target || cmd+".exe" == target {
		return true
	}

	// When the contained executable name is full name (e.g. foo_darwin_amd64),
	// it is also regarded as a target executable file.
	for _, delimiter := range []rune{'_', '-'} {
		c := fmt.Sprintf("%s%c%s%c%s", cmd, delimiter, os, delimiter, arch)
		if os == "windows" {
			c += ".exe"
		}
		if c == target {
			return true
		}
	}

	return false
}

func unarchiveTar(src io.Reader, cmd, os, arch string) (io.Reader, error) {
	t := tar.NewReader(src)
	for {
		h, err := t.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%w tar file: %s", ErrCannotDecompressFile, err)
		}
		_, name := filepath.Split(h.Name)
		if matchExecutableName(cmd, os, arch, name) {
			return t, nil
		}
	}
	return nil, fmt.Errorf("%w in tar: %q", ErrExecutableNotFoundInArchive, cmd)
}
