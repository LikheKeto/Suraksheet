FROM python:3.9-slim

WORKDIR /app

# Set environment variables
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/

# Install system dependencies
RUN apt-get update && apt-get install -y \
    tesseract-ocr \
    tesseract-ocr-eng \
    tesseract-ocr-nep \ 
    libgl1-mesa-glx \
    zlib1g-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the application code
COPY . .

# Run the application
CMD ["python", "main.py"]