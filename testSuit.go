package testSuit

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TestSuit http request test form
type TestSuit struct {
	Method      string
	URL         string
	Data        url.Values
	Buffer      io.Reader
	ContentType string
	Router      *gin.Engine
}

// Do start test
func (ts *TestSuit) Do(handler gin.HandlerFunc) (*httptest.ResponseRecorder, *httptest.ResponseRecorder) {
	var req *http.Request

	if ts.ContentType == "" {
		ts.ContentType = "application/x-www-form-urlencoded"
	}

	if ts.Method == "ALL" {
		if ts.Data != nil {
			encodedData := ts.Data.Encode()

			reqGet, _ := http.NewRequest("GET", ts.URL, nil)
			reqGet.Header.Add("Content-Type", ts.ContentType)
			reqGet.URL.RawQuery = encodedData

			reqPost, _ := http.NewRequest("POST", ts.URL,
				strings.NewReader(encodedData))
			reqPost.Header.Add("Content-Type", ts.ContentType)

			respGet := httptest.NewRecorder()
			respPost := httptest.NewRecorder()
			ts.Router.ServeHTTP(respGet, reqGet)
			ts.Router.ServeHTTP(respPost, reqPost)

			return respGet, respPost
		}

		req, _ = http.NewRequest(ts.Method, ts.URL, nil)
		req.Header.Add("Content-Type", ts.ContentType)

	} else {
		if ts.Data != nil {

			switch ts.Method {
			case "GET":
				req, _ = http.NewRequest(ts.Method, ts.URL, nil)
				req.URL.RawQuery = ts.Data.Encode()
				req.Header.Add("Content-Type", ts.ContentType)
			case "POST":
				req, _ = http.NewRequest(ts.Method, ts.URL,
					strings.NewReader(ts.Data.Encode()))
				req.Header.Add("Content-Type", ts.ContentType)
			}

		} else if ts.Buffer != nil {
			req, _ = http.NewRequest(ts.Method, ts.URL, ts.Buffer)
			req.Header.Add("Content-Type", ts.ContentType)
		} else {
			req, _ = http.NewRequest(ts.Method, ts.URL, nil)
			req.Header.Add("Content-Type", ts.ContentType)
		}
	}

	resp := httptest.NewRecorder()
	ts.Router.ServeHTTP(resp, req)
	return resp, nil
}

// GetGinEngine create gin.Engine with test mode
func GetGinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	return r
}

// GenImage generate randomize image
func GenImage(ext string, width, height int) (*bytes.Buffer, string, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
			})
		}
	}

	nameFix := uuid.New().URN()
	buffer := &bytes.Buffer{}
	bufferWriter := multipart.NewWriter(buffer)
	formWriter, err := bufferWriter.CreateFormFile("file", nameFix+".png")
	contentType := bufferWriter.FormDataContentType()

	switch ext {
	case "png":
		png.Encode(formWriter, img)
		bufferWriter.Close()
	case "jpg":
		jpeg.Encode(formWriter, img, nil)
		bufferWriter.Close()
	case "gif":
		gif.Encode(formWriter, img, nil)
		bufferWriter.Close()
	}
	return buffer, contentType, err
}
