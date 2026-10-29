FROM public.ecr.aws/d3j8x8q7/olympus-base-python:latest

WORKDIR /app
COPY . .

RUN pip install --no-cache-dir \
    asgiref==3.12.1 \
    attrs==26.1.0 \
    colorama==0.4.6 \
    Django==6.1.2 \
    docopt==0.6.2 \
    iniconfig==2.3.1 \
    packaging==26.3 \
    parso==0.8.7 \
    pluggy==1.6.0 \
    Pygments==2.21.0 \
    pytest==8.4.2 \
    sqlparse==0.6.0 \
    typing_extensions==4.16.0

# The typeshed and django-stubs submodules are required by setup.py; fetch
# them at their pinned commits when the source comes without submodules.
RUN set -e; \
    fetch() { \
      if [ ! -f "$1/$4" ]; then \
        rm -rf "$1" && git init -q "$1" \
          && git -C "$1" fetch -q --depth 1 "$2" "$3" \
          && git -C "$1" checkout -q FETCH_HEAD \
          && rm -rf "$1/.git"; \
      fi; \
    }; \
    fetch jedi/third_party/typeshed https://github.com/davidhalter/typeshed.git \
      4bb9d8351d0795fbc7dee3da069013a48641689d LICENSE; \
    fetch jedi/third_party/django-stubs https://github.com/davidhalter/django-stubs.git \
      fd057010f6cbf176f57d1099e82be46d39b99cb9 LICENSE.txt

RUN pip install --no-cache-dir --no-deps -e .

CMD ["/bin/bash"]
