# tinypdf

Reduce PDF Size

## Configurations Possible

### Modes

Default values set when only the mode is specified.

1. `screen`
   1. `dColorImageResolution`: 72
   2. `dColorImageDownsampleType`: `/Average`
   3. `dMonoImageResolution`: 300
   4. `dMonoImageDownsampleType`: `/Subsample`
   5. `dGrayImageResolution`: 72
   6. `dGrayImageDownsampleType`: `/Average`
   7. `dColorConversionStrategy`: `RGB`
2. `default`
   1. `dColorImageResolution`: 72
   2. `dColorImageDownsampleType`: `/Subsample`
   3. `dMonoImageResolution`: 300
   4. `dMonoImageDownsampleType`: `/Subsample`
   5. `dGrayImageResolution`: 72
   6. `dGrayImageDownsampleType`: `/Subsample`
   7. `dColorConversionStrategy`: `LeaveColorUnchanged`
3. `ebook`
   1. `dColorImageResolution`: 150
   2. `dColorImageDownsampleType`: `/Average`
   3. `dMonoImageResolution`: 300
   4. `dMonoImageDownsampleType`: `/Subsample`
   5. `dGrayImageResolution`: 150
   6. `dGrayImageDownsampleType`: `/Bicubic`
   7. `dColorConversionStrategy`: `RGB`
4. `printer`
   1. `dColorImageResolution`: 300
   2. `dColorImageDownsampleType`: `/Average`
   3. `dMonoImageResolution`: 1200
   4. `dMonoImageDownsampleType`: `/Subsample`
   5. `dGrayImageResolution`: 300
   6. `dGrayImageDownsampleType`: `/Bicubic`
   7. `dColorConversionStrategy`: `UseDeviceIndependentColor`
5. `prepress`
   1. `dColorImageResolution`: 300
   2. `dColorImageDownsampleType`: `/Bicubic`
   3. `dMonoImageResolution`: 1200
   4. `dMonoImageDownsampleType`: `/Subsample`
   5. `dGrayImageResolution`: 300
   6. `dGrayImageDownsampleType`: `/Bicubic`
   7. `dColorConversionStrategy`: `LeaveColorUnchanged`

### Configurations to tinker

- `EmbedAllFonts`
- `dColorImageDownsampleThreshold`
- `dUCRandBGInfo`
- `dCompressEntireFile` & `dCompressPages` (pretty useless though)

## Packages

- https://github.com/MrSaints/go-ghostscript
- [qpdf](https://github.com/qpdf/qpdf)
  ```
  qpdf --object-streams=generate --stream-data=compress small.pdf small-v2.pdf
  ```
- [pdfsizeopt](https://github.com/pts/pdfsizeopt)
  ```
  ./pdfsizeopt --optimize --quiet ~/Documents/personal/tinypdf/tinypdf-small.pdf  test.pdf
  ```

## Notes

### `dColorImageDownsampleType`, `dMonoImageDownsampleType`, `dGrayImageDownsampleType`

```
0 = None. No downsampling occurs.

1 = Average. Downsampling averages the pixel color in a sample group and replaces the entire group with an average of the pixel color at the specified resolution.

2 = Bicubic. Downsampling occurs using a weighted average of the pixel color in a sample group and replaces the entire group with an average of the pixel color at the specified resolution. Bicubic downsampling is more precise, resulting in smoother gradations than Average and Subsample downsampling, but it is also the slowest process.

3 = Subsample. Subsampling chooses a pixel in the center of a sample group and replaces the entire group with that pixel at the specified resolution. Subsampling is the fastest sampling process, but can result in a less smooth and continuous image.
```

Size impact: ✅ Yes — /Subsample gives smallest output, /Bicubic best visuals.
Tweakable?: ✅ Use /Subsample for aggressive compression.

- https://mostafaznv.gitbook.io/pdf-optimizer/advanced-usage/configuring-ghostscript/colorimagedownsampletype


### `dColorConversionStrategy`

- `LeaveColorUnchanged`: Keep original color space
- `RGB`: Convert everything to RGB
- `CMYK`: Convert everything to CMYK
  
Size impact: ✅ Yes — converting to RGB can reduce file size.
Tweakable?: ✅ Set to RGB if you're not targeting print.

### `dColorImageDownsampleThreshold`

Only downsample an image if its actual resolution is more than X times the target resolution.

```
-dColorImageResolution=150
-dColorImageDownsampleThreshold=1.5

```

An image in the PDF is:
- 300 DPI → gets downsampled to 150 DPI (since 300 > 150 × 1.5)
- 200 DPI → untouched (200 < 225)
- 150 DPI → untouched

### `dUCRandBGInfo`

> Transfer functions (dot gain compensation, gamma correction) can be applied to the PDF file, altering the colour values, but allowing the information to be dropped, if it is present in the input file. Set -dTransferFunctionInfo to /Apply instead of /Preserve.
> Undercolour removal and black generation functions are used when converting RGB to CMYK, and PDF files can carry around rules on how to do this. Since printers will always have their own defaults, it is safe to drop this too by setting `UCRandBGInfo` to /Remove.

[Source](https://ghostscript.com/blog/optimizing-pdfs.html#:~:text=Assume%20we've%20set%20the,results%20but%20will%20take%20longer)


## Debugging

### pdfimages

```
page   num  type   width height color comp bpc  enc interp  object ID x-ppi y-ppi size ratio
--------------------------------------------------------------------------------------------
   5     0 image    1231  1032  rgb     3   8  image  no        20  0   410   410  269K 7.2%
   5     1 smask    1231  1032  gray    1   8  image  no        20  0   410   410 1253B 0.1%
   6     2 image    2717  1658  rgb     3   8  image  no        26  0   814   814 81.8K 0.6%
   6     3 smask    2717  1658  gray    1   8  image  no        26  0   814   814 48.4K 1.1%
   7     4 image    2471  1821  rgb     3   8  image  no        31  0   823   823 85.9K 0.7%
   7     5 smask    2471  1821  gray    1   8  image  no        31  0   823   823 4396B 0.1%
```


### Ghostscript options for individual modes

```
gs \
  -q \
  -dNODISPLAY \
  -c ".distillersettings /screen get {exch ==only ( ) print ===} forall quit" \
| sort
```

src: https://superuser.com/a/440573


## Comparisons

- Better than sejda[.]com
- Better than https://pdfcpu[.]io/core/optimize