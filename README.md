# Tool to download files from google drive

This tool was written out of frustration by the lack of resume support when using the browser and the lack of progress  information when downloading with the browser.

## Help

```
./gdrivedl --help
sage of ./gdrivedl:
 -dest string
   	Destination of file
 -url string
   	Url to download from
```

## Example output

```
./gdrivedl -url 'gooldriveid'
2020/03/27 01:59:22 Download id gooledriveid
2020/03/27 01:59:23 Found token XXXX
2020/03/27 01:59:23 Download url https://docs.google.com/uc?confirm=XXXX&export=download&id=googledriveid
2020/03/27 01:59:25 Status: 206 Partial Content
2020/03/27 01:59:25 File size 830184773
   5% |██                                      | (18.1 MB/s) [2s:43s
```
