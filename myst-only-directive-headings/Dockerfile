FROM public.ecr.aws/d3j8x8q7/olympus-base-python:latest

WORKDIR /app
COPY . .

RUN pip install --no-cache-dir \
    alabaster==1.0.0 \
    babel==2.18.0 \
    beautifulsoup4==4.15.0 \
    certifi==2026.7.22 \
    charset-normalizer==3.5.2 \
    coverage==7.16.2 \
    defusedxml==0.7.1 \
    docutils==0.22.4 \
    idna==3.20 \
    imagesize==2.0.1 \
    iniconfig==2.3.1 \
    Jinja2==3.1.6 \
    linkify-it-py==2.2.0 \
    markdown-it-py==4.2.0 \
    MarkupSafe==3.0.4 \
    mdit-py-plugins==0.6.1 \
    mdurl==0.1.2 \
    packaging==26.3 \
    pluggy==1.6.0 \
    Pygments==2.20.0 \
    pytest==9.1.1 \
    pytest-cov==7.1.0 \
    pytest-datadir==1.8.0 \
    pytest-regressions==2.11.0 \
    pytest_param_files==0.6.0 \
    PyYAML==6.0.3 \
    requests==2.34.2 \
    roman-numerals==4.1.0 \
    ruamel.yaml==0.19.1 \
    snowballstemmer==3.1.1 \
    soupsieve==3.0.3 \
    Sphinx==9.1.0 \
    sphinx_pytest==0.3.0 \
    sphinxcontrib-applehelp==2.0.0 \
    sphinxcontrib-devhelp==2.0.0 \
    sphinxcontrib-htmlhelp==2.1.0 \
    sphinxcontrib-jsmath==1.0.1 \
    sphinxcontrib-qthelp==2.0.0 \
    sphinxcontrib-serializinghtml==2.0.0 \
    typing_extensions==4.16.0 \
    urllib3==2.8.0

RUN pip install --no-cache-dir --no-deps -e .

CMD ["/bin/bash"]
