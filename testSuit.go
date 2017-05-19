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
	request     *http.Request
	response    *httptest.ResponseRecorder
	Header      map[string]string
	Router      *gin.Engine
}

// Do start test
func (ts *TestSuit) Do() *httptest.ResponseRecorder {
	if ts.ContentType == "" {
		ts.ContentType = "application/x-www-form-urlencoded"
	}

	addHeader := func() {
		ts.request.Header.Add("Content-Type", ts.ContentType)

		for k, v := range ts.Header {
			ts.request.Header.Add(k, v)
		}

	}

	if ts.Data != nil {
		encodedData := ts.Data.Encode()

		switch ts.Method {
		case "ALL":
			ts.request, _ = http.NewRequest("POST", ts.URL,
				strings.NewReader(encodedData))
			addHeader()
		case "GET":
			ts.request, _ = http.NewRequest(ts.Method, ts.URL, nil)
			ts.request.URL.RawQuery = ts.Data.Encode()
			addHeader()
		case "POST":
			ts.request, _ = http.NewRequest(ts.Method, ts.URL,
				strings.NewReader(ts.Data.Encode()))
			addHeader()
		}

	} else if ts.Buffer != nil {
		ts.request, _ = http.NewRequest(ts.Method, ts.URL, ts.Buffer)
		addHeader()
	} else {
		ts.request, _ = http.NewRequest(ts.Method, ts.URL, nil)
		addHeader()
	}

	ts.response = httptest.NewRecorder()
	ts.Router.ServeHTTP(ts.response, ts.request)
	return ts.response
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
