import core
import os

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

    subprocess.run([
        'python', '/content/voicevox_core/example/python/run.py',
        '--root_dir_path', "/content/voicevox_core/release",
        '--use_gpu',
        '--speaker_id', str(speaker_id),
        '--text', text,
        '--output_file', filepath
    ])
    # subprocess.run([
    #     'echo', filepath
    # ])

    return filepath
