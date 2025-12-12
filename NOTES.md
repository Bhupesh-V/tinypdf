<!-- omit from toc -->
# Notes

- [Checklist](#checklist)
  - [General](#general)
  - [Fonts](#fonts)
  - [Images](#images)
- [Resources](#resources)
- [Tools](#tools)
- [Author Notes](#author-notes)

## Checklist

The following checklist contains items that help in reducing PDF file sizes with an indication of whether its taken care of in `tinypdf`

Reference: https://github.com/pts/pdfsizeopt/tree/master/pts_pdfsizeopt2009_talk

### General

- [ ] Remove extra whitespace and comments.
- [x] Compress streams with high-eﬀort ZIP.
- [ ] Use cross-reference streams (with the y predictor).
- [x] Use object streams.
- [ ] Remove page thumbnails.
- [ ] Flatten the page structure.
- [ ] Inline indirect references.
- [ ] Remove keys ignored by PDF specification.
- [ ] Compact the cross reference tables.
- [ ] Find duplicate or equivalent objects, and keep only one copy.

### Fonts

- [ ] Convert Type 1 (PostScript Type 1) fonts to [CFF](https://adobe-type-tools.github.io/font-tech-notes/pdfs/5176.CFF.pdf) (Type 1C, Type 2)
- [ ] Unify [subsets](https://fonts.google.com/knowledge/glossary/subsetting) of the same font.
- [ ] Embed vector fonts.


### Images

- [ ] Convert some inline images to objects to help deduplication.
- [ ] Use grayscale or a palette instead of `RGB` or `CMYK`.
- [ ] Use the smallest bit depth.
- [ ] Get rid of image duplicates based on pixel colors.
- [x] Reduce image resolution (300 DPI or 600 DPI): no need for a higher resolution than the printer’s for the scaled image.
- [x] Choose the JPEG quality.

<!-- ## Questions

- **What is a stream in PDF?**
  - A stream in a PDF is a sequence of bytes that can store large data, like images, fonts, or page content. Streams are usually compressed and are paired with a dictionary that describes their metadata.
  ```
  << /Length 123 /Filter /FlateDecode >>
  stream
  ...binary data...
  endstream
  ``` -->

## Resources

- [The Structure of a PDF File](https://medium.com/@jberkenbilt/the-structure-of-a-pdf-file-6f08114a58f6)
- [PDF Specificaiton 2.0](https://pdfa.org/resource/iso-32000-2/) - Its behind a paywall (although free) and requires personal details to get access (weird?).
- [PDF Specification 1.7 (free)](https://opensource.adobe.com/dc-acrobat-sdk-docs/pdfstandards/PDF32000_2008.pdf)

## Tools

The PDF optimise ecosystem seems to be weirdly big, I cane able to discover following tools:

- qpdf
- cpdf
- mutool
- ghostscript
- pdfsizeopt
- pdftocairo
- pdfcpu
- ocrmypdf

## Author Notes

- Make a PDF structure readable with `qpdf --qdf --object-streams=disable original.pdf  output.pdf`.
- Check images inside a PDF, using
  - `cpdf`
    ```bash
    cpdf -list-images in.pdf
    ```
  - `pdfimages`
    ```bash
    pdfimages -list tinypdf-original.pdf
    ```
- gs fails to run properly with `ebook` preset. Logs
  ```
  Failed to initialise downsample filter, downsampling aborted
  ```
- Keep an eye out for MacOS support on building executables from docker image using [dockerc](https://github.com/NilsIrl/dockerc).