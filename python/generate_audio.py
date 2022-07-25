from urllib import response
import core
import os
import requests


import random
import string
from pathlib import Path
import subprocess



def get_random_string(length):
    # choose from all lowercase letter
    letters = string.ascii_lowercase
    result_str = ''.join(random.choice(letters) for i in range(length))
    return result_str

def audio_dir(html_file):
    path = Path(html_file)

    file_name = os.path.basename(html_file)
    volume_name = os.path.splitext(file_name)[0]

    audio_volume_dir = os.path.join(path.parent, "audio", volume_name)
    Path(audio_volume_dir).mkdir(parents=True, exist_ok=True)

    return audio_volume_dir


def generate_audio(text: str, speaker_id: int, html_file: str) -> str:
    filepath = f"{os.path.join(audio_dir(html_file), get_random_string(8))}.wav"

    params = dict(
        speaker=str(speaker_id),
        text=text
    )

    audio_query_resp = requests.post('http://localhost:50021/audio_query', params=params)

    params = dict(
        speaker=str(speaker_id)
    )

    audio_resp = requests.post('http://localhost:50021/synthesis', params=params, data=audio_query_resp.content)

    with open(filepath, "wb") as f:
        f.write(audio_resp.content)

    return filepath
