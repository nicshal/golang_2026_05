package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
)

const blockSize = 8192

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	// не даем копировать файл сам в себя
	if fromPath == toPath {
		return fmt.Errorf("error copying file onto itself")
	}

	// Открытие исходного файла
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
			slog.Error("error closing source file: %v", "error", err)
		}
	}(srcFile)

	// Получение размера исходного файла
	srcStat, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting source file size: %w", err)
	}
	srcSize := srcStat.Size()

	if srcSize == 0 && !srcStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	// Проверка offset и limit
	if offset > srcSize {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 {
		limit = srcSize
	}
	if limit > srcSize {
		limit = srcSize
	}

	// Открытие целевого файла для записи
	dstFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error creating target file: %w", err)
	}
	defer func(dstFile *os.File) {
		err := dstFile.Close()
		if err != nil {
			slog.Error("error closing target file: %v", "error", err)
		}
	}(dstFile)

	// Перемещение курсора исходного файла на заданный offset
	if _, err := srcFile.Seek(offset, 0); err != nil {
		return fmt.Errorf("error moving source file cursor: %w", err)
	}

	// Копирование данных
	buf := make([]byte, blockSize)
	totalCopied := int64(0)
	ratio := 0.0
	for totalCopied < limit {
		n, err := srcFile.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("error reading source file: %w", err)
		}

		if n == 0 {
			break
		}

		if n > int(limit-totalCopied) {
			n = int(limit - totalCopied)
		}

		_, err = dstFile.Write(buf[:n])
		if err != nil {
			return fmt.Errorf("error writing to target file: %w", err)
		}

		totalCopied += int64(n)
		ratioNew := float64(totalCopied) * 100 / float64(limit)
		if int(ratioNew) > int(ratio) {
			fmt.Printf("Progress %s -> %s: %.2f%%\n", fromPath, toPath, ratioNew)
		}
		ratio = ratioNew
	}

	return nil
}
