# go-image-scale-script

This project is a Go-based command-line application for processing images. It resizes images to predefined dimensions and optionally adds a watermark.

## Features
- Loads configuration from a `.env` file.
- Resizes images to specified dimensions.
- Supports watermarking for large, medium, and small sizes.
- Ensures processed files belong to a specific user.

## Installation

### Install Dependencies
First, install Go dependencies:

```sh
 go mod tidy
```

### Configure `.env` File
Create a `.env` file in the project directory and add the following variables:

```ini
OUTPUT_BASE_DIR=/path/to/output
OWNER_USER=username
WATERMARK_FILE=/path/to/watermark.png
DIMENSION_S=200
DIMENSION_M=500
DIMENSION_L=1000
DIMENSION_XL=2000
```

## Usage

### Basic Usage
Run the following command to process an image:

```sh
go run main.go -a /path/to/image.jpg
```

This processes the image in all sizes.

### Options

| Parameter | Description |
|-----------|-------------|
| `-env <path>` | Specifies a custom `.env` file. Default: `.env` |
| `-w` | Adds a watermark. |
| `-a` | Processes all sizes. |
| `-s` | Processes only the small size. |
| `-m` | Processes only the medium size. |
| `-l` | Processes only the large size. |
| `-xl` | Processes only the extra-large size. |

### Example Commands

#### Add a watermark to medium and large sizes only
```sh
go run main.go -env ./.env -m -l -w /path/to/image.jpg
```

#### Process only the small size
```sh
go run main.go ./.env -s /path/to/image.jpg
```

## Build and Run
To build the project as a standalone binary:

```sh
go build -o image-processor
```

Then run:

```sh
./image-processor -env ./.env -a /path/to/image.jpg
```

