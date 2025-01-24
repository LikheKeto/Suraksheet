FROM python:3.9-slim

WORKDIR /app

ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/
RUN apt-get update && apt-get install -y tesseract-ocr tesseract-ocr-eng libgl1-mesa-glx zlib1g-dev

ADD requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY . .


CMD ["python", "main.py"]
