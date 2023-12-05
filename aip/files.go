package aip

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sunshineplan/imgconv"
)

func MakeZip(name string) {
  zipFile, err := os.Create(name)
  if err != nil {
    panic(err)
  }
  defer zipFile.Close()

  w := zip.NewWriter(zipFile)
  defer w.Close()

  walker := func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        f, err := w.Create(path)
        if err != nil {
            return err
        }

        _, err = io.Copy(f, file)
        if err != nil {
            return err
        }

        return nil
    }
  err = filepath.Walk("charts", walker)
  if err != nil {
    panic(err)
  }
}

func saveChartToDisk(chart Chart, b64 string) {
  dec, err := base64.StdEncoding.DecodeString(b64)
  if err != nil {
    panic(err)
  }
  filePath := fmt.Sprintf("charts/byop/%s", chart.Name)
  os.WriteFile(filePath, dec, 0644)
}

func saveAsPdf(chart Chart, b64 string) {
  dec, err := base64.StdEncoding.DecodeString(b64)
  if err != nil {
    panic(err)
  }
  breader := bytes.NewReader(dec)
  img, err := imgconv.Decode(breader)
  if err != nil {
    panic(err)
  }
  fHandle, err := os.Create(fmt.Sprintf("charts/byop/%s", chart.Name))
  
  err = imgconv.Write(fHandle, img, &imgconv.FormatOption{Format: imgconv.PDF})
  if err != nil {
    panic(err)
  }
}
