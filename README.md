# tinypdf

Reduce/Compress PDF file size

## What is this?

At its core, `tinypdf` is a utility that uses multiple PDF tools behind the scenes to reduce file size, with a touch of configuration magic to make it seamless.

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

## Packages

- https://github.com/MrSaints/go-ghostscript
- [pdfsizeopt](https://github.com/pts/pdfsizeopt)
  ```
  ./pdfsizeopt --optimize --quiet ~/Documents/personal/tinypdf/tinypdf-small.pdf  test.pdf
  ```

## Comparisons

TODO

<!-- - Better than sejda[.]com -->