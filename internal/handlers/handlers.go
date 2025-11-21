package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"project/internal/models"
	"project/internal/pdf"
	"sync"
	"time"
)

type Storage interface {
	SaveLinks(models.LinksResponce)
	GetLinks(linksNum []int32) []models.LinksResponce
	GetLinksNum() int
}

func CheckLinks(s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var linksRequest models.LinksRequest
		linksResponce := models.LinksResponce{
			Links:    make(map[string]string),
			LinksNum: s.GetLinksNum(),
		}

		err := json.NewDecoder(r.Body).Decode(&linksRequest)
		if errors.Is(err, io.EOF) {
			http.Error(w, "request body is empty", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "something wrong please retry", http.StatusBadRequest)
			return
		}
		if linksRequest.Links == nil {
			http.Error(w, "please enter some links", http.StatusBadRequest)
			return
		}

		wg := sync.WaitGroup{}

		client := &http.Client{
			Timeout: 2 * time.Second,
		}

		for _, url := range linksRequest.Links {
			if _, ok := linksResponce.Links[url]; !ok {
				wg.Add(1)
				go func() {
					defer wg.Done()
					resp, err := client.Get("https://" + url)
					if err != nil {
						linksResponce.Links[url] = "not available"
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode == http.StatusOK {
						linksResponce.Links[url] = "available"
					} else {
						linksResponce.Links[url] = "not available"
					}
				}()
			} else {
				continue
			}
		}
		wg.Wait()

		err = json.NewEncoder(w).Encode(linksResponce)
		if err != nil {
			http.Error(w, "something wrong please retry", http.StatusBadRequest)
			return
		}
		s.SaveLinks(linksResponce)
	}
}

func PastLinks(s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var pastLinksRequest models.PastLinksRequest

		err := json.NewDecoder(r.Body).Decode(&pastLinksRequest)
		if errors.Is(err, io.EOF) {
			http.Error(w, "request body is empty", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "something wrong please retry", http.StatusBadRequest)
			return
		}
		if pastLinksRequest.LinksNum == nil {
			http.Error(w, "please enter some numbers", http.StatusBadRequest)
			return
		}

		pdfBytes := pdf.CreateBeautifulStructPDF(s.GetLinks(pastLinksRequest.LinksNum))

		filename := "links.pdf"

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

		pdfBytes.Output(w)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "go.html")
}
