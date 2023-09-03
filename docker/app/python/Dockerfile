FROM python:3.11
ENV LC_ALL=C.UTF-8 \
    LANG=C.UTF-8

WORKDIR /app
RUN apt-get update \
    && apt-get -y install default-mysql-client \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /data

COPY admin/ishocon2.dump.tar.bz2 /data
COPY docker/app/entrypoint.sh /
COPY webapp/python/requirements.txt ./
RUN pip3 install --no-cache-dir -r ./requirements.txt

COPY webapp/python/ /app/

CMD ["uwsgi", "--ini", "app.ini", "--python-autoreload"]
