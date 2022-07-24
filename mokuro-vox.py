import argparse
from email.mime import audio
import os
from pathlib import Path
from bs4 import BeautifulSoup

from generate_audio import generate_audio

def set_audio_play_script(soup):
    audio_play_script = """const audioIcons = document.querySelectorAll('.audio-icon');
audioIcons.forEach(audio => {
  audio.addEventListener('click', function handleClick(event) {
    event.target.parentNode.querySelector('audio').play();
  });
});"""
    for body in soup.select("body"):
        script = soup.new_tag("script")
        script.string = audio_play_script
        body.append(script)

def set_audio_tags(item, soup, filepath):
    count = 0
    print(filepath)

    audio = soup.new_tag("audio", hidden=True, preload="none")
    audio.append(soup.new_tag("source", src=filepath, type="audio/mpeg"))

    audio_icon = soup.new_tag('p', attrs={"class": "audio-icon"})
    audio_icon.append(soup.new_tag('img', src="https://super.so/icon/dark/play-circle.svg"))
    audio_icon.append(audio)

    item.append(audio_icon)
    count+=1

def update_html(html_file):
    with open(html_file, encoding='utf-8') as f:
        html = f.read()

    soup = BeautifulSoup(html, "html.parser")

    set_audio_play_script(soup)

    count = 0
    for item in soup.select(".textBox"):
        if count == 100:
            break

        filepath = generate_audio(item.get_text, 11, audio_dir(html_file))
        set_audio_tags(item, soup, filepath)


    with open('mangas/updated.html', 'w') as writer:
        writer.write(soup.prettify())

def audio_dir(html_file):
    path = Path(html_file)

    file_name = os.path.basename(html_file)
    volume_name = os.path.splitext(file_name)[0]

    audio_volume_dir = os.path.join(path.parent, "audio", volume_name)
    Path(audio_volume_dir).mkdir( parents=True, exist_ok=True )

    return audio_volume_dir

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--html-file", type=str)
    parser.add_argument("--audio-dir", type=str)

    args = parser.parse_args()

    update_html(args.html_file)
