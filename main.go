package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"
	"simjoha/stego-demo/stego"
)

const BUFSIZE = 4096

var flg_ifile string // the file to encode
var flg_cfile string // the file that should carry the encoded data
var flg_ofile string // the encoded output file
var flg_read bool    // read stego payload
var flg_write bool   // write stego payload
var stderr *log.Logger

func init() {
	flag.BoolVar(&flg_read, "r", false, "Read stego payload.")
	flag.BoolVar(&flg_write, "w", false, "Write stego payload.")
	flag.StringVar(&flg_ifile, "i", "", "The file you wish to hide.")
	flag.StringVar(&flg_ofile, "o", "", "Output image.")
	flag.StringVar(&flg_cfile, "c", "", "The carrier image file.")
	flag.Parse()
	stderr = log.New(os.Stderr, "", 0)
	if flg_read == flg_write {
		log.Fatalln("You must specify exactly one of -r or -w.")
	}
}

func main() {
	if flg_write {
		StegoWrite(flg_ifile, flg_cfile, flg_ofile)
	} else {
		StegoRead(flg_ifile, flg_ofile)
	}
}

func StegoWrite(ifile_path, cfile_path, ofile_path string) {
	var ifile, ofile, cfile *os.File
	var err error
	var img image.Image
	var ftype string
	b := make([]byte, BUFSIZE)
	if ifile_path == "" {
		log.Fatalln("Missing input file.")
	}
	if cfile_path == "" {
		stderr.Fatalln("Missing carrier file.")
	}
	if ofile_path == "" {
		ofile_path = cfile_path
	}
	if ifile, err = os.Open(ifile_path); err != nil {
		log.Fatalln(err)
	}
	defer ifile.Close()
	if cfile, err = os.Open(cfile_path); err != nil {
		log.Fatalln(err)
	}
	defer cfile.Close()
	if img, ftype, err = image.Decode(cfile); err != nil {
		log.Fatalln(err)
	}
	steg := stego.NewImage(img)
	info, _ := ifile.Stat()
	size := info.Size()
	if size > int64(steg.Entropy()) {
		log.Fatalln("Not enough entropy in carrier image.")
	}
	w := steg.Writer()
	for size > 0 {
		n, _ := ifile.Read(b)
		w.Write(b[:n])
		size -= int64(n)
	}
	if ofile, err = os.OpenFile(ofile_path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0755); err != nil {
		log.Fatalln(err)
	}
	defer ofile.Close()
	switch ftype {
	case "png":
		png.Encode(ofile, steg.Export())
	default:
		stderr.Fatalln("Image type not supported.")
	}
}

func StegoRead(ifile_path, ofile_path string) {
	var ifile, ofile *os.File
	var err error
	var img image.Image
	b := make([]byte, BUFSIZE)
	if ifile_path == "" {
		log.Fatalln("Missing input file.")
	}
	if ifile, err = os.Open(ifile_path); err != nil {
		log.Fatalln(err)
	}
	defer ifile.Close()
	if ofile_path == "" {
		ofile_path = ifile_path + ".out"
	}
	if ofile, err = os.OpenFile(ofile_path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0755); err != nil {
		log.Fatalln(err)
	}
	defer ofile.Close()
	if img, _, err = image.Decode(ifile); err != nil {
		log.Fatalln(err)
	}
	steg := stego.NewImage(img)
	r := steg.Reader()
	for n := 0; err == nil; {
		n, err = r.Read(b)
		ofile.Write(b[:n])
	}
}
