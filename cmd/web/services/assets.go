package services

import (
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/klauspost/compress/zstd"
)

//go:embed global.css index.js
var fdefault embed.FS

type encoderWeight struct {
	name   string
	weight float32
}

// Structure that handles serving all assets, treating caching and compression.
type AssetHandler struct {
	start     time.Time
	fs        fs.FS
	gzipFiles map[string]bytes.Buffer
	zstdFiles map[string]bytes.Buffer
}

// Creates and compresses the files from the filesystem specified.
// If fs is nil, the embed one is used instead.
func NewAssetHandler(fs interface {
	fs.ReadDirFS
	fs.ReadFileFS
}) (*AssetHandler, error) {
	if fs == nil {
		fs = fdefault
	}

	handler := &AssetHandler{
		start:     time.Now().UTC().Round(time.Second),
		fs:        fs,
		gzipFiles: make(map[string]bytes.Buffer),
		zstdFiles: make(map[string]bytes.Buffer),
	}

	files, err := fs.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to open directory: %v", err)
	}

	for _, file := range files {
		if err := handler.storeFile(file); err != nil {
			return nil, fmt.Errorf("failed to store file: %v", err)
		}
	}

	return handler, nil
}


func (ah *AssetHandler) storeFile(file fs.DirEntry) error {
	if file.Type().IsRegular() {
		name := file.Name()
		switch mime.TypeByExtension(filepath.Ext(name)) {
		case "text/javascript; charset=utf-8":
			fallthrough
		case "text/css; charset=utf-8":
			err := compress(name, ah.fs, compressGzip, ah.gzipFiles)
			if err != nil {
				return fmt.Errorf("failed to compress to gzip: %v", err)
			}

			err = compress(name, ah.fs, compressZstd, ah.zstdFiles)
			if err != nil {
				return fmt.Errorf("failed to compress to gzip: %v", err)
			}
		}
	}

	return nil
}

func (ah *AssetHandler) writeResponse(weighted []encoderWeight, filename string, w http.ResponseWriter) {
	mimetype := mime.TypeByExtension(filepath.Ext(filename))
	w.Header().Set("Content-Type", mimetype)

	switch mimetype {
	case "text/css; charset=utf-8":
		fallthrough
	case "text/javascript; charset=utf-8":
		enc_best := weighted[len(weighted)-1].name
		switch enc_best {
		case "gzip":
			w.Header().Set("Content-Encoding", "gzip")
			if file, ok := ah.gzipFiles[filename]; ok {
				w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Len()))
				io.Copy(w, &file)
				fmt.Println("")
			} else {
				w.WriteHeader(500)
			}
		case "zstd":
			w.Header().Set("Content-Encoding", "zstd")
			if file, ok := ah.zstdFiles[filename]; ok {
				w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Len()))
				io.Copy(w, &file)
			} else {
				w.WriteHeader(500)
			}
		default:
			if file, err := ah.fs.Open(filename); err == nil {
				stat, err := file.Stat()
				if err != nil {
					w.WriteHeader(404)
					io.WriteString(w, "404 file not found")
					return
				}

				w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
				io.Copy(w, file)
				file.Close()
			}
		}
	}
}

func (ah *AssetHandler) getWeightedEncoders(encodings []string) []encoderWeight {
	encs := make([]encoderWeight, 0, 5)
	for i := range encodings {
		enc := strings.Split(encodings[i], ";q=")
		weight := float32(0)
		if len(enc) > 1 {
			w, err := strconv.ParseFloat(enc[1], 32)
			if err == nil {
				weight = float32(w)
			}
		}
		encs = append(encs, encoderWeight{
			enc[0], float32(weight),
		})
	}

	sort.Slice(encs, func(i, j int) bool {
		// adjust weights based on server preference if they're equal
		if encs[i].weight == encs[j].weight {
			if encs[i].name == "gzip" && encs[j].name == "zstd" {
				encs[j].weight += 0.1
			} else if encs[i].name == "zstd" && encs[j].name == "gzip" {
				encs[i].weight += 0.1
			}
		}

		return encs[i].weight < encs[j].weight
	})

	return encs
}

// Function that handles HTTP requests.
// **You are required to have an URL parameter of `filename`, else the function
// will return `404 file not found`!**
func (ah *AssetHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if f, err := ah.fs.Open(filename); err != nil {
		w.WriteHeader(404)
		io.WriteString(w, "404 file not found")
		return
	} else {
		f.Close()
	}

	encodings := strings.Split(r.Header.Get("Accept-Encoding"), ", ")
	weighted := ah.getWeightedEncoders(encodings)
	file_time, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since"))

	if (ah.start.Before(file_time) || ah.start.Equal(file_time)) && err == nil {
		w.WriteHeader(304)
	} else {
		w.Header().Set("Last-Modified", ah.start.UTC().Format(http.TimeFormat))
		w.Header().Set("Cache-Control", "max-age=0")
		ah.writeResponse(weighted, filename, w)
	}
}

type compressFunc func(file fs.File) (bytes.Buffer, error)

func compress(path string, fs fs.FS, cf compressFunc, store map[string]bytes.Buffer) error {
			file, err := fs.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}

			data, err := cf(file)
			if err != nil {
				return fmt.Errorf("failed to compress file: %v", err)
			}

			store[path] = data
			err = file.Close()
			if err != nil {
				return fmt.Errorf("failed to close file: %v", err)
			}

			return nil
}

func compressGzip(file fs.File) (bytes.Buffer, error) {
	var gzip_data bytes.Buffer
	enc := gzip.NewWriter(&gzip_data)

	_, err := io.Copy(enc, file)
	if err != nil {
		return gzip_data, fmt.Errorf("failed to write gzip: %v", err)
	}

	err = enc.Close()
	if err != nil {
		return gzip_data, fmt.Errorf("failed to close gzip writer: %v", err)
	}

	return gzip_data, nil
}

func compressZstd(file fs.File) (bytes.Buffer, error) {
	var zstd_data bytes.Buffer

	enc, err := zstd.NewWriter(&zstd_data, zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(19)))
	if err != nil {
		return zstd_data, fmt.Errorf("failed to create encoder: %v", err)
	}

	_, err = io.Copy(enc, file)
	if err != nil {
		return zstd_data, fmt.Errorf("failed to write zstd: %v", err)
	}

	err = enc.Close()
	if err != nil {
		return zstd_data, fmt.Errorf("failed to close zstd writer: %v", err)
	}

	return zstd_data, nil
}
