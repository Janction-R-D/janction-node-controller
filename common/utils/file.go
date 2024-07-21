package utils

import (
	"archive/zip"
	"bufio"
	"debug/elf"
	"debug/pe"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func MatchELFFast(ident []byte) bool {
	if len(ident) < 4 {
		return false
	}
	return ident[0] == '\x7f' && ident[1] == 'E' && ident[2] == 'L' && ident[3] != 'F'
}

func MatchPEFast(ident []byte) bool {
	if len(ident) < 2 {
		return false
	}
	return ident[0] == 'M' && ident[1] == 'Z'
}

func IsExecutableFile(path string) bool {
	if IsElfFile(path) {
		return true
	}

	return IsPEFile(path)
}

func IsElfFile(path string) bool {
	f, err := elf.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}

func IsPEFile(path string) bool {
	f, err := pe.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}

// 打包成zip文件
func Zip(src_dir string, zip_file_name string) {
	os.Chdir(src_dir)
	src_dir = "./"
	// 预防：旧文件无法覆盖
	os.RemoveAll(zip_file_name)
	// 创建：zip文件
	zipfile, _ := os.Create(zip_file_name)
	defer zipfile.Close()
	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	// 遍历路径信息
	filepath.Walk(src_dir, func(path string, info os.FileInfo, _ error) error {
		// 如果是源路径，提前进行下一个遍历
		if path == src_dir {
			return nil
		}
		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src_dir+`\`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}
		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
}

const (
	KiB = 1024
	MiB = KiB * 1024
	GiB = MiB * 1024
	TiB = GiB * 1024
)

func FormatFileSize(byteNum int64, decimal int) (to string) {
	size := float64(byteNum) / KiB
	if size < 1024 {
		return fmt.Sprintf("%s KiB", strconv.FormatFloat(size, 'f', decimal, 64))
	}
	size = float64(byteNum) / MiB
	if size < 1024 {
		return fmt.Sprintf("%s MiB", strconv.FormatFloat(size, 'f', decimal, 64))
	}
	size = float64(byteNum) / GiB
	if size < 1024 {
		return fmt.Sprintf("%s GiB", strconv.FormatFloat(size, 'f', decimal, 64))
	}
	size = float64(byteNum) / TiB
	return fmt.Sprintf("%s TiB", strconv.FormatFloat(size, 'f', decimal, 64))
}

func CopyFile(src, dst string) (err error) {
	srcFile, err := os.OpenFile(src, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcBuf := bufio.NewReader(srcFile)
	dstBuf := bufio.NewWriter(dstFile)
	_, err = io.Copy(dstBuf, srcBuf)
	if err != nil {
		return err
	}
	return dstBuf.Flush()
}
