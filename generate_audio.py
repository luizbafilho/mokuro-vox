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


def generate_audio(text: str, speaker_id: int, out_dir: str) -> str:
    filename = get_random_string(8)
    filepath = f"{os.path.join(out_dir, filename)}.wav"

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
