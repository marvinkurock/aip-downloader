#!/bin/bash

# install with pip install download-tiles
# might need pip install setuptools

download-tiles icao.mbtiles --tiles-url=https://ais.dfs.de/static-maps/icao500/tiles/{z}/{x}/{y}.png --bbox=5.8663,47.2701,15.0419,55.0584 --zoom-levels=7-12 --verbose
