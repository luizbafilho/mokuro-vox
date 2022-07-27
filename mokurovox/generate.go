package mokurovox

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

func GenerateAudio(htmlPath string, speakerId string, overrideHtml bool, speed float64) error {
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
		audioPath, err := generateAudio(audioDir, s.Text(), speakerId, speed)
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

	updatedPath := htmlPath
	if !overrideHtml {
		updatedPath = filepath.Join(path.Dir(htmlPath), volumeName(htmlPath)+"-with-audio.html")
	}

	err = ioutil.WriteFile(updatedPath, []byte(html), 0644)
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

func volumeName(htmlPath string) string {
	volumeName := strings.TrimSuffix(htmlPath, filepath.Ext(htmlPath))
	return strings.TrimPrefix(volumeName, path.Dir(htmlPath))

}

func createAudioDir(htmlPath string) (string, error) {
	parent := path.Dir(htmlPath)

	volumeName := volumeName(htmlPath)

	path := filepath.Join(parent, "audio", volumeName)
	err := os.MkdirAll(path, os.ModePerm)

	return path, err
}

func generateAudio(audioDir string, text string, speakerId string, speed float64) (string, error) {
	audioQuery, err := getAudioQuery(speakerId, text)
	if err != nil {
		return "", fmt.Errorf("failed getting audio_query: %w", err)
	}
	defer audioQuery.Close()

	audioQuery, err = setSpeed(audioQuery, speed)
	if err != nil {
		return "", fmt.Errorf("failed setting speed: %w", err)
	}

	audioPath, err := synthesizeAudio(speakerId, audioQuery, audioDir)
	if err != nil {
		return "", err
	}

	return audioPath, nil
}

func setSpeed(body io.ReadCloser, speed float64) (io.ReadCloser, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed audio query response: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil { // Parse []byte to go struct pointer
		return nil, fmt.Errorf("failed audio query JSON: %w", err)
	}

	result["speedScale"] = speed
	data, err = json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed marshalling updated audio query: %w", err)
	}

	return io.NopCloser(bytes.NewReader(data)), nil
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

func synthesizeAudio(speakerId string, audioQuery io.ReadCloser, audioDir string) (string, error) {
	resp, err := http.Post(fmt.Sprintf("http://localhost:50021/synthesis?speaker=%s", speakerId), "application/json", audioQuery)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpAudioPath := filepath.Join(audioDir, randomString(5))
	out, err := os.Create(tmpAudioPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed copying audio=%s to disk: %w", tmpAudioPath, err)
	}

	audioPath := tmpAudioPath + ".mp3"
	err = exec.Command("ffmpeg", "-i", tmpAudioPath, audioPath).Run()
	if err != nil {
		return "", fmt.Errorf("failed converting wav to mp3: %w", err)
	}

	return audioPath, os.Remove(tmpAudioPath)
}
