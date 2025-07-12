# tinypdf

Reduce PDF file size

## What is this?

At its core, `tinypdf` is a wrapper based utility that uses multiple PDF tools behind the scenes to reduce file sizes, with a touch of configuration magic to make it seamless.

## Why?

- I wasn't really comfortable uploading sensitive files on web-server based PDF compression tools. Existing GUI based tools on my system resulted in unsatisfactory results (I was [greedy](https://x.com/bhupeshimself/status/1941499700802355317), since I wanted to reduce a file from 19MB to less than 1.5MB).
- I like doing this [kind](https://bhupesh.me/publishing-my-first-ever-dockerfile-optimization-ugit/) of [stuff](https://bhupesh.me/minimalist-guide-git-clone/) (reducing asset sizes).

## Installation

```bash
# TODO
docker pull bhupeshimself/tinypdf
docker run --rm -it -v $(pwd):/app bhupeshimself/tinypdf
```

If everything goes good, you should see a welcome text like this

```

▗ ▘       ▌▐▘
▜▘▌▛▌▌▌▛▌▛▌▜▘
▐▖▌▌▌▙▌▙▌▙▌▐
     ▄▌▌


Usage: tinypdf -i input.pdf [-preset screen|ebook|printer|prepress] [-quality 50]

Example: tinypdf -i input.pdf -preset ebook -quality 60

Options:
  -preset         One of: screen, ebook, printer, prepress (default: screen)
  -quality        Quality percentage between 10 and 90 (default: 50)

```

## Usage

TODO

## Acknowledgements

`tinypdf` would not have been possible without these tools 💙

- [qpdf](https://github.com/qpdf/qpdf)
- [ghostscript](https://ghostscript.readthedocs.io/en/latest/index.html)
- [poppler-utils](https://poppler.freedesktop.org/)
- [cpdf](https://community.coherentpdf.com/)

<!-- ## Comparisons -->

TODO

<!-- - Better than sejda[.]com -->

## Alternatives

- `ocrmypdf` on its own does a [good job](https://ocrmypdf.readthedocs.io/en/latest/optimizer.html) for reducing PDF size. However watch out for unintentional image removals from the file.
- [`pdfsizeopt`](https://github.com/pts/pdfsizeopt/) is probably one of the earliest tools to optimize PDF file size.

## Assets

The PDF assets used to test `tinypdf` are available on Internet Archive.

- [The_Grand_Design](https://archive.org/details/The_Grand_Design)