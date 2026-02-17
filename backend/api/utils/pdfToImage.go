package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	. "backend/api/config"
)

// POST request onto pdf creation Locat !
// api/utils/pdfToImage.go
func RequestToPdfToImage(path string, name string) bool {
	fmt.Println("📄 PDF source:", path)
	fmt.Println("🖼 Thumbnail name:", name)
	//	sendRequest(path, name, "https://pdf2png.sheetable.net/createthumbnail")
	sendRequest(path, name, "http://localhost:5000/createthumbnail")
	return true
}

func sendRequest(pdfPath string, name string, remoteURL string) bool {
	file, err := os.Open(pdfPath)
	if err != nil {
		log.Println("open pdf:", err)
		return false
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(pdfPath))
	if err != nil {
		log.Println(err)
		return false
	}
	io.Copy(part, file)

	writer.WriteField("name", name)
	writer.Close()

	req, err := http.NewRequest("POST", remoteURL, body)
	if err != nil {
		log.Println(err)
		return false
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("request failed:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("pdf2png returned:", resp.Status)
		return false
	}

	thumbnailPath := path.Join(
		Config().ConfigPath,
		"sheets/thumbnails",
		name+".png",
	)

	out, err := os.Create(thumbnailPath)
	if err != nil {
		log.Println(err)
		return false
	}
	defer out.Close()

	io.Copy(out, resp.Body)
	return true
}

func Upload(client *http.Client, url string, values map[string]io.Reader, name string) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	/*
		Don't forget to close the multipart writer.
		If you don't close it, your request will be missing the terminating boundary.
	*/
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}

	// Save response
	defer res.Body.Close()
	out, err := os.Create(path.Join(Config().ConfigPath, "sheets/thumbnails", name+".png"))
	if err != nil {
		return
	}
	defer out.Close()
	io.Copy(out, res.Body)

	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
