
package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"image/gif"
	"io"
	"os"
	"flag"
)

func main() {

	format := flag.String("format", "", "image format, choose of jpeg, gif, png")
	flag.Parse()
	in := os.Stdin
	img, kind, err := image.Decode(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decode: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)

	switch *format {
	case "jpeg":
		if err := toJPEG(img, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
			os.Exit(1)
		}
	case "png":
		if err := toPNG(img, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "png: %v\n", err)
			os.Exit(1)
		}
	case "gif":
		if err := toGIF(img, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "gif: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unkown format %s\n", *format)
		os.Exit(1)
	}
}

func toJPEG(img image.Image, out io.Writer) error {
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}
func toGIF(img image.Image, out io.Writer) error {
	return gif.Encode(out, img, &gif.Options{})
}
func toPNG(img image.Image, out io.Writer) error {
	return png.Encode(out, img)
}
//!-main

/*
//!+with
$ go build gopl.io/ch3/mandelbrot
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
Input format = png
//!-with

//!+without
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
jpeg: image: unknown format
//!-without
*/
