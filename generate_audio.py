import core
import os

import random
import string
from pathlib import Path


def get_random_string(length):
    # choose from all lowercase letter
    letters = string.ascii_lowercase
    result_str = ''.join(random.choice(letters) for i in range(length))
    return result_str


def generate_audio(text: str, speaker_id: int, out_dir: str) -> str:
    filename = get_random_string(8)

    core.initialize(True, 0)
    core.voicevox_load_openjtalk_dict("open_jtalk_dic_utf_8-1.11")

    wavefmt = core.voicevox_tts(text, speaker_id)

    filepath = f"{os.path.join(out_dir, filename)}.wav"

    with open(filepath, "wb") as f:
        f.write(wavefmt)

    core.finalize()

    return filepath
