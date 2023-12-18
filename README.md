# csguard

**csguard** is a command-line interface (CLI) application written in Go for calculating and validating checksums for files and folders. The current version supports the MD5 hashing algorithm.

## Features

- Calculate MD5 checksum for files and folders
- Validate checksums against provided values
- Easy-to-use command-line interface

## Usage

### Calculate Checksum

To calculate the MD5 checksum for a file, use the following command:

```bash
csguard calculate --input-file=path/to/file.txt --output=result.txt
```
OR
```bash
csguard calculate --input-folder=path/to/folder --output=result.txt
```
If `--output` is not provided, only console output will show.

Avaliable format are `.json`, `.txt`, `.yaml`.

### Validate Checksum

To validate the MD5 checksum for a file, use the following command:

```bash
csguard validate --input-file=path/to/file.txt --checksum=checksum-value-here --output=result.txt
```
OR
To validate the MD5 checksum for multiple files, use the following command:

```bash
csguard validate --checksum-file=path/to/file.txt --output=result.txt
```

Available format for `--checksum-file` is `.json`, `.txt`, `.yaml`.

`json` format
```json
[
	{
		"file_name": "folder1\\f610zrj8bn281.png",
		"value": "654b0435d0b202ac1654b79d088d4be5"
	},
	{
		"file_name": "folder1\\sub1\\2702140556378.jpg",
		"value": "75e61e44d231ae781e335e2703e94914"
	},
	{
		"file_name": "folder1\\1702140556378.jpg",
		"value": "75e61e44d231ae781e335e2703e94914"
	}
]
```
Warning: filepath are marshalled to `\\` but validate clean up those so no need to worries about `\\`

`txt` format
```text
f610zrj8bn281.png 654b0435d0b202ac1654b79d088d4be5
```

`yaml` format
```yaml
- file_name: f610zrj8bn281.png
  value: 654b0435d0b202ac1654b79d088d4be5
```

Available format for `--output` are `.json`, `.txt`, `.yaml`.

## License

The MIT License (MIT)

Copyright (c) 2015 Chris Kibble

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)