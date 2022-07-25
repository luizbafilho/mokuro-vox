package mokurovox

import (
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
)

var (
	//go:embed audio-play.html
	audioPlayScript string
	//go:embed audio-tag.html
	audioTagHtml string
)

func GenerateAudio(htmlPath string, speakerId string) error {
	var htmlFile, err = os.OpenFile(htmlPath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer htmlFile.Close()

	doc, err := goquery.NewDocumentFromReader(htmlFile)
	if err != nil {
		return err
	}

	textBoxes := doc.Find(".textBox")

	bar := progressbar.Default(int64(textBoxes.Length()))

	audioDir, err := createAudioDir(htmlPath)
	if err != nil {
		return err
	}

	textBoxes.Each(func(i int, s *goquery.Selection) {
		audioPath, err := generateAudio(audioDir, s.Text(), speakerId)
		if err != nil {
			fmt.Printf("failed generating audio for text=%s err=%s\n", s.Text(), err)
		}

		setAudioTag(s, relativePath(htmlPath, audioPath))

		bar.Add(1)
	})

	setAudioPlayScript(doc)
	html, err := doc.Html()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("mangas/Volume1-updated.html", []byte(html), 0644)
	if err != nil {
		return err
	}

	return nil
}

func relativePath(htmlPath string, audioPath string) string {
	parent := path.Dir(htmlPath)

	path, err := filepath.Rel(parent, audioPath)
	if err != nil {
		log.Fatal(err)
	}

	return path
}

func createAudioDir(htmlPath string) (string, error) {
	parent := path.Dir(htmlPath)

	volumeName := strings.TrimSuffix(htmlPath, filepath.Ext(htmlPath))
	volumeName = strings.TrimPrefix(volumeName, parent)

	path := filepath.Join(parent, "audio", volumeName)
	err := os.MkdirAll(path, os.ModePerm)

	return path, err
}

func generateAudio(audioDir string, text string, speakerId string) (string, error) {
	audioPath := filepath.Join(audioDir, randomString(5)+".wav")

	audioQuery, err := getAudioQuery(speakerId, text)
	if err != nil {
		return "", fmt.Errorf("failed getting audio_query: %w", err)
	}
	defer audioQuery.Close()

	if err := synthesizeAudio(speakerId, audioQuery, audioPath); err != nil {
		return "", err
	}

	return audioPath, nil
}

func setAudioPlayScript(doc *goquery.Document) {
	doc.Find("body").AppendHtml(audioPlayScript)
}

func setAudioTag(s *goquery.Selection, audioPath string) {
	s.AppendHtml(fmt.Sprintf(audioTagHtml, audioPath))
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func getAudioQuery(speakerId string, text string) (io.ReadCloser, error) {
	data := url.Values{}
	data.Set("speaker", speakerId)
	data.Set("text", text)
	encodedData := data.Encode()

	response, err := http.Post(fmt.Sprintf("http://localhost:50021/audio_query?%s", encodedData), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func synthesizeAudio(speakerId string, audioQuery io.ReadCloser, audioPath string) error {
	resp, err := http.Post(fmt.Sprintf("http://localhost:50021/synthesis?speaker=%s", speakerId), "application/json", audioQuery)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(audioPath)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed copying audio=%s to disk: %w", audioPath, err)
	}

	return nil
}
