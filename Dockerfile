
FROM ubuntu:20.04

ARG DEBIAN_FRONTEND=noninteractive
WORKDIR /app
# install all dependencies
RUN apt-get -y update && apt-get install -y python3.9 \
ffmpeg \
python3-pip \
curl \
portaudio19-dev \
alsa-base \
alsa-utils \
libsndfile1-dev && \
apt-get clean

RUN pip3 install -U openai-whisper
# test openai whisper
RUN whisper test/test.mp3 --language Spanish --output_dir output
# install go
RUN apt-get install -y wget
RUN wget https://golang.org/dl/go1.22.0.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
RUN go version
# install aws cli
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN apt-get install -y unzip
RUN unzip awscliv2.zip
RUN ./aws/install
RUN aws configure set aws_access_key_id $AWS_ACCESS_KEY_ID
RUN aws configure set aws_secret_access_key $AWS_SECRET_ACCESS_KEY
RUN aws configure set default.region $AWS_DEFAULT_REGION

COPY . .
# install go packages
RUN go mod download
COPY *.go ./
RUN go build -o jarvis main.go
CMD ["./jarvis"]
