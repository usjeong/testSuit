# testSuit
personal test tool for gin based http server

[![Go Report Card](https://goreportcard.com/badge/github.com/usjeong/testSuit)](https://goreportcard.com/report/github.com/usjeong/testSuit)

## Usage
```go
var (
	caseOne = conf.NewCaseOne("develop")
	App     = setApp()
)

func setApp() *gin.Engine {
	r := testSuit.GetGinEngine()
	api.NewApp(caseOne)
	api.SetRouter(r)
	return r
}

func TestPing(*testing.T) {
	suit := &testSuit.TestSuit{
		Router: App,
		Method: "GET",
		URL:    "/ping",
	}

	resp := suit.Do()
	assert.Equal(t, 200, resp.Code)
}

```

## Other functions
```go
// generate the random color image
func GenImage(ext string, width, height int) (*bytes.Buffer, string, error) {}

// return gin engine with test mode
func GetGinEngine() *gin.Engine {}

```




