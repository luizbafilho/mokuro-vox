FROM python:3.7

ENV VIRTUAL_ENV=/opt/venv
RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

RUN mkdir /content

WORKDIR /content

RUN git clone https://github.com/VOICEVOX/voicevox_core && \
    cd /content/voicevox_core &&\
    git checkout release-0.11

RUN apt update -y && apt install -y cmake libsndfile1-dev

RUN cd /content/voicevox_core && \
    echo "3" | python configure.py --voicevox_version "0.11.4" &&\
    pip install -r requirements.txt && \
    pip install . && \
    cd /content/voicevox_core/example/python && \
    pip install -r requirements.txt
