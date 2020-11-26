FROM python:3.8-stretch as builder

RUN mkdir /app
WORKDIR /app

COPY requirements.txt requirements.txt

RUN pip install -r requirements.txt

FROM builder

COPY --from=builder /install /usr/local
COPY src .

CMD ["gunicorn", "-w 4", "main:app"]