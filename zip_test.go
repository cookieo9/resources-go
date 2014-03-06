package resources

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"reflect"
	"strings"
	. "testing"
)

var ZipBundle_Is_A_Bundle Bundle = &zipBundle{}
var ZipBundle_Is_A_Searcher Searcher = &zipBundle{}
var ZipBundle_Is_A_Lister Lister = &zipBundle{}

func base64_decode(b64 string) []byte {
	d := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64))
	if data, err := ioutil.ReadAll(d); err != nil {
		panic(err)
	} else {
		return data
	}
	panic("unreachable")
}

var files = []struct {
	Path     string
	Contents []byte
}{
	{"foo.txt", []byte("foo is foo")},
	{"subfolder/bar.txt", []byte("bar is not foo")},
	{"logo.ico", base64_decode(`AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD///8AVE44//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2/1ROOP////8A////AFROOP/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv9UTjj/////AP///wBUTjj//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb/VE44/////wD///8AVE44//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2/1ROOP////8A////AFROOP/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv9UTjj/////AP///wBUTjj//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb//uF2//7hdv/+4Xb/VE44/////wD///8AVE44//7hdv/+4Xb//uF2//7hdv/+4Xb/z7t5/8Kyev/+4Xb//993///dd///3Xf//uF2/1ROOP////8A////AFROOP/+4Xb//uF2//7hdv//4Hn/dIzD//v8///7/P//dIzD//7hdv//3Xf//913//7hdv9UTjj/////AP///wBUTjj//uF2///fd//+4Xb//uF2/6ajif90jMP/dIzD/46Zpv/+4Xb//+F1///feP/+4Xb/VE44/////wD///8AVE44//7hdv/z1XT////////////Is3L/HyAj/x8gI//Is3L////////////z1XT//uF2/1ROOP////8A19nd/1ROOP/+4Xb/5+HS//v+//8RExf/Liwn//7hdv/+4Xb/5+HS//v8//8RExf/Liwn//7hdv9UTjj/19nd/1ROOP94aDT/yKdO/+fh0v//////ERMX/y4sJ//+4Xb//uF2/+fh0v//////ERMX/y4sJ//Ip07/dWU3/1ROOP9UTjj/yKdO/6qSSP/Is3L/9fb7//f6///Is3L//uF2//7hdv/Is3L////////////Is3L/qpJI/8inTv9UTjj/19nd/1ROOP97c07/qpJI/8inTv/Ip07//uF2//7hdv/+4Xb//uF2/8zBlv/Kv4//pZJU/3tzTv9UTjj/19nd/////wD///8A4eLl/6CcjP97c07/e3NO/1dOMf9BOiX/TkUn/2VXLf97c07/e3NO/6CcjP/h4uX/////AP///wD///8A////AP///wD///8A////AP///wDq6/H/3N/j/9fZ3f/q6/H/////AP///wD///8A////AP///wD///8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==`)},
}

func CreateTestZip(t *T) *bytes.Reader {
	t.Log("Creating in-memory zip file")
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for _, file := range files {
		t.Logf("Adding file %q (%d bytes) to zip file", file.Path, len(file.Contents))

		fw, err := zw.Create(file.Path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fw.Write(file.Contents); err != nil {
			t.Fatal(err)
		}
	}

	t.Logf("Creating Manifest in zip file")
	if fw, err := zw.Create("MANIFEST"); err != nil {
		t.Fatal(err)
	} else {
		for _, file := range files {
			mimestring := http.DetectContentType(file.Contents)
			mediatype, params, err := mime.ParseMediaType(mimestring)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Fprintf(fw, "%s (%d bytes): %s %v\n", file.Path, len(file.Contents), mediatype, params)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	t.Log("Finished zip file,", buf.Len(), "bytes written.")
	return bytes.NewReader(buf.Bytes())
}

func TestZip(t *T) {
	t.Parallel()

	zip := CreateTestZip(t)
	zb, err := OpenZipReader(zip, int64(zip.Len()))
	if err != nil {
		t.Fatal(err)
	}

	t.Log()
	t.Log("--- Check Contents --- ")
	for _, file := range files {
		if rdr, err := zb.Open(file.Path); err != nil {
			t.Error(err)
		} else {
			if data, err := ioutil.ReadAll(rdr); err != nil {
				t.Error(err)
			} else {
				if reflect.DeepEqual(data, file.Contents) {
					t.Log(file.Path, "- Contents Match")
				} else {
					t.Error(file.Path, "- Contents Differ")
				}
			}
		}
	}

	t.Log()
	t.Log("--- Zip Listing ---")
	lister := zb.(Lister)
	if list, err := lister.List(); err != nil {
		t.Error(err)
	} else {
		t.Log("Resources:", list)
	}

	searcher := zb.(Searcher)
	t.Log()
	t.Log("--- Find everything ---")
	for _, file := range files {
		if r, err := searcher.Find(file.Path); err != nil {
			t.Error(err)
		} else {
			t.Log("Found:", r)
		}
	}

	t.Log()
	t.Log("--- Glob root level stuff ---")
	if globset, err := searcher.Glob("*"); err != nil {
		t.Error(err)
	} else {
		t.Log("Found:", globset)
	}
}

func TestWriteZip(t *T) {
	if Short() {
		return
	}
	t.Parallel()

	zip := CreateTestZip(t)
	if f, err := os.Create("testdata.zip"); err != nil {
		t.Fatal(err)
	} else {
		if _, err := io.Copy(f, zip); err != nil {
			t.Fatal(err)
		}
	}
}
