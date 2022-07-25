import argparse
from email.mime import audio
import os
from pathlib import Path
from bs4 import BeautifulSoup
from tqdm import tqdm


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
    audio = soup.new_tag("audio", hidden=True, preload="none")
    audio.append(soup.new_tag("source", src=filepath, type="audio/mpeg"))

    audio_icon = soup.new_tag('p', attrs={"class": "audio-icon"})
    audio_icon.append(soup.new_tag('img', src="https://super.so/icon/dark/play-circle.svg"))
    audio_icon.append(audio)

    item.append(audio_icon)

def strip_absolute_path(html_file, audio_path):
    parent_path = Path(html_file).parent.absolute()
    return os.path.relpath(audio_path, parent_path)

def update_html(html_file):
    with open(html_file, encoding='utf-8') as f:
        html = f.read()

    soup = BeautifulSoup(html, "html.parser")

    set_audio_play_script(soup)

    for item in tqdm(soup.select(".textBox")):
        audiopath = generate_audio(item.get_text(), 13, html_file)
        set_audio_tags(item, soup, strip_absolute_path(html_file, audiopath))


    path = Path(html_file)
    file_name = os.path.basename(html_file)
    volume_name = os.path.splitext(file_name)[0]

    updatedFilename = os.path.join(path.parent.absolute(), f"{volume_name}-with-audio.html")

    with open(updatedFilename, 'w') as writer:
        writer.write(soup.prettify())



if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--html-file", type=str)

    args = parser.parse_args()

    update_html(args.html_file)
