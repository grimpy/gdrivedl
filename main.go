package main

import "flag"
import "strings"
import "log"
import "io"
import "os"
import "mime"
import "fmt"
import "strconv"
import "net/url"
import "net/http"
import "net/http/cookiejar"
import "github.com/schollz/progressbar"

func downloadByID(downloadid string, destination string) {
	downloadurl, _ := url.Parse("https://docs.google.com/uc")
	query := downloadurl.Query()
	query.Add("id", downloadid)
	query.Add("export", "download")
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}
	downloadurl.RawQuery = query.Encode()
	resp, err := client.Get(downloadurl.String())
	if err != nil {
		log.Fatalf("Failed to download %s", err)
	}
	defer resp.Body.Close()
	token := ""
	for _, cookie := range resp.Cookies() {
		if strings.HasPrefix(cookie.Name, "download_warning") {
			token = cookie.Value
			break
		}
	}
	if len(token) == 0 {
		log.Println("No token found")
	} else {
		log.Printf("Found token %s\n", token)
		query.Add("confirm", token)
		downloadurl.RawQuery = query.Encode()
	}

	log.Printf("Download url %s\n", downloadurl.String())
	req, err := http.NewRequest("GET", downloadurl.String(), nil)
	if err != nil {
		log.Fatal("Could not make request\n")
	}
	req.Header.Add("Range", "bytes=0-")
	resp, err = client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("Could not make request\n")
	}
	totalsize := resp.Header.Get("Content-Length")
	if len(totalsize) == 0 {
		// let's try get it from Content-Range
		contentrange := resp.Header.Get("Content-Range")
		if len(contentrange) != 0 {
			items := strings.Split(contentrange, "/")
			totalsize = items[len(items)-1]
		}
	}
	var bytesize int64
	if len(totalsize) != 0 {
		bytesize, err = strconv.ParseInt(totalsize, 10, 64)

	}
	log.Printf("File size %d\n", bytesize)
	if len(destination) == 0 {
		disposition := resp.Header.Get("Content-Disposition")
		_, params, err := mime.ParseMediaType(disposition)
		if err != nil {
			log.Fatalf("Failed to find desination name")
		}
		destination = params["filename"]
	}
	partfile := destination + ".part"
	stat, err := os.Stat(partfile)
	var startsize int64
	if err == nil {
		startsize = stat.Size()
		resp.Body.Close()
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-", startsize))
		resp, err = client.Do(req)
	}
	file, err := os.OpenFile(partfile, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatalf("Failed to open file %s\n", destination)
	}
	file.Seek(startsize, 0)
	defer file.Close()
	bar := progressbar.NewOptions(
		int(bytesize),
		progressbar.OptionSetBytes(int(bytesize)),
	)
	bar.Set(int(startsize))
	out := io.MultiWriter(file, bar)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Failed to download %s\n", err)
	}
	os.Rename(partfile, destination)

}

func main() {
	var downloadurl string
	var dest string
	var downloadid string
	flag.StringVar(&downloadurl, "url", "", "Url to download from")
	flag.StringVar(&dest, "dest", "", "Destination of file")
	flag.Parse()
	parsedurl, err := url.Parse(downloadurl)
	if err != nil {
		downloadid = downloadurl
	} else {
		if strings.HasPrefix(parsedurl.Scheme, "http") {
			downloadid = parsedurl.Query().Get("id")
			if len(downloadid) == 0 {
				log.Fatal("Invalid url")
			}
		} else {
			downloadid = downloadurl
		}
	}
	log.Printf("Download id %s\n", downloadid)
	downloadByID(downloadid, dest)
}
