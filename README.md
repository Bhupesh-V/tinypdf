# tinypdf

Reduce PDF file size

## What is this?

`tinypdf` is a utility that acts as a wrapper (or orchestrator) around multiple PDF manipulation tools, using _okayish_ configurations to reduce PDF file size.

## Why?

- I wasn't really comfortable uploading sensitive files on web-server based PDF compression tools. Existing GUI based tools on my system resulted in unsatisfactory results (I was [greedy](https://x.com/bhupeshimself/status/1941499700802355317), since I wanted to reduce a file from 19MB to less than 1.5MB).
- I seem to have a weird liking for doing this [kind](https://bhupesh.me/publishing-my-first-ever-dockerfile-optimization-ugit/) of [stuff](https://bhupesh.me/minimalist-guide-git-clone/).

## Installation

```bash
# TODO
# Pull the Image first: docker pull bhupeshimself/tinypdf
docker run --rm -v $(pwd):/app bhupeshimself/tinypdf
```

If everything goes good, you should see a welcome text like this

```


▗ ▘       ▌▐▘
▜▘▌▛▌▌▌▛▌▛▌▜▘
▐▖▌▌▌▙▌▙▌▙▌▐
     ▄▌▌


Usage:
        tinypdf -i input.pdf [-preset screen|ebook|printer|prepress] [-quality 50]

Example:
        docker run --rm -v $(pwd):/app bhupeshimself/tinypdf -i input.pdf -preset printer -quality 40

Options:
        -preset         One of: screen, ebook, printer, prepress (default: screen)
        -quality        Quality percentage between 10 and 90 (default: 50)

Bugs:
        Please report any issues on the GitHub repository.
                https://github.com/Bhupesh-V/tinypdf/issues

```

## Usage

Follow the instructions shared on `tinypdf --help`

Sample run

```
docker run --rm -v $(pwd):/app tinypdf:latest -i board_certificate.pdf -preset ebook -quality 88
```

Sample output:

```
Original:   976 kB
Compressed: 294 kB
Reduced by: 682 kB (69.9%)
```

## Debug

To report & debug any issues, use the `TINYPDF_DEBUG=1` env variable:

```bash
docker run --rm -v $(pwd):/app -e TINYPDF_DEBUG=1 tinypdf:latest -i board_certificate.pdf -preset ebook -quality 88
```

Sample log

```
2026/02/18 13:49:01 /usr/bin/pdftocairo -pdf board_certificate.pdf /tmp/tinypdf-pdftocairo-2051690074.pdf
2026/02/18 13:49:01 /usr/bin/gs -sDEVICE=pdfwrite -dCompatibilityLevel=1.4 -dBATCH -dNOPAUSE -dQUIET -dDownsampleColorImages=true -dColorImageDownsampleType=/Average -dDownsampleMonoImages=true -dMonoImageDownsampleType=/Subsample -dDownsampleGrayImages=true -dGrayImageDownsampleType=/Average -dColorImageResolution=132 -dMonoImageResolution=264 -dGrayImageResolution=132 -sOutputFile=/tmp/tinypdf-gs-3152810261.pdf /tmp/tinypdf-pdftocairo-2051690074.pdf
2026/02/18 13:49:02 /usr/bin/qpdf --object-streams=generate --stream-data=compress /tmp/tinypdf-gs-3152810261.pdf /tmp/tinypdf-qpdf-2016115314.pdf

Original:   976 kB
Compressed: 294 kB
Reduced by: 682 kB (69.9%)
```

## Acknowledgements

`tinypdf` would not have been possible without these tools 💙

- [qpdf](https://github.com/qpdf/qpdf)
- [ghostscript](https://ghostscript.readthedocs.io/en/latest/index.html)
- [poppler-utils](https://poppler.freedesktop.org/)
- [cpdf](https://community.coherentpdf.com/)

<!-- - Better than sejda[.]com -->

## Alternatives

- `ocrmypdf` on its own does a [good job](https://ocrmypdf.readthedocs.io/en/latest/optimizer.html) for reducing PDF size. However watch out for unintentional image removals from the file.
- [`pdfsizeopt`](https://github.com/pts/pdfsizeopt/) is probably one of the earliest tools to optimize PDF file size.

## Assets

The PDF asset(s) used to test `tinypdf` are available on Internet Archive (except my 10th board certificate).

- [The_Grand_Design](https://archive.org/details/The_Grand_Design)