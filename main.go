package main

import (
	"os"

	"github.com/mkurock/aipDownloader/aip"
)

func main()  {
  airports := aip.ParseConfigJson()
  os.RemoveAll("charts")
  os.MkdirAll("charts/byop", 0755)
  aip.DownloadCharts(airports)
  aip.SetManifestJson()
  aip.MakeZip("aip_bundle.zip")
}

