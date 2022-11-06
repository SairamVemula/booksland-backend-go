package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 2 // 2MB

func saveFileinAssests(fileHeader *multipart.FileHeader) (filelocation string, err error) {
	// truncated for brevity

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, err := fileHeader.Open()
	if err != nil {
		return
	}
	defer file.Close()

	//Takeing first 521 btye to determine MIME
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return "", err
	}

	//using DetectContentType to ge filetype
	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/jpg" && filetype != "image/webp" {
		return "", errors.New(fmt.Sprintf("The provided file format %s is not allowed. Please upload a JPEG or PNG or WEBP image", filetype))
	}

	//The file.Seek() method is used to return the pointer back to the start of the file so that io.Copy() starts from the beginning
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./assets", os.ModePerm)
	if err != nil {
		return
	}
	// Create a new file in the uploads directory
	filelocation = fmt.Sprintf("/assets/%s-%d%s", strings.Split(strings.ReplaceAll(fileHeader.Filename, " ", "-"), ".")[0], time.Now().UnixMilli(), filepath.Ext(fileHeader.Filename))
	dst, err := os.Create("." + filelocation)
	if err != nil {
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		return
	}

	log.Println("Upload successful  ->  ", "."+filelocation)

	return
}

func UploadFile(r *http.Request, key string) (string, error) {
	log.Println(key)
	_, fileHeader, err := r.FormFile(key)
	if err != nil || fileHeader == nil {
		return "empty", nil
	}
	filelocation, err := saveFileinAssests(fileHeader)
	if err != nil {
		return "", err
	}
	return filelocation, nil
}

func UploadFiles(r *http.Request, key string) ([]string, error) {
	var filelocations []string
	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}
	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File[key]
	// if !ok {
	// 	log.Println(key, len(files))
	// }
	// log.Println(ok, key, len(files))
	for _, fileHeader := range files {
		// Restrict the size of each uploaded file to 1MB.
		// To prevent the aggregate size from exceeding
		// a specified value, use the http.MaxBytesReader() method
		// before calling ParseMultipartForm()
		if fileHeader.Size > MAX_UPLOAD_SIZE {
			return []string{}, errors.New(fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 2MB in size", fileHeader.Filename))
		}
		filelocation, err := saveFileinAssests(fileHeader)
		if err != nil {
			return nil, err
		}
		filelocations = append(filelocations, filelocation)
	}
	return filelocations, nil
}
