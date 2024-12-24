# Example Dockerized Go Application

[![Build Status](https://semaphore-demos.semaphoreci.com/badges/golang-mathapp/branches/main.svg?key=3c5ab555-0c9b-44df-806e-cacc1b5b416b)](https://semaphore-demos.semaphoreci.com/projects/golang-mathapp)

An example demo project for [How To Deploy a Go Web Application with Docker](https://semaphoreci.com/community/tutorials/how-to-deploy-a-go-web-application-with-docker)

## Docker Images

Build development image:

```bash
$ docker build -t mathapp-development .
```

Run container:

```bash
$ docker run -it --rm \
    -p 8010:8010 \
    -v $PWD/src:/go/src/mathapp \
    mathapp-development
```

## License

MIT License

Copyright (c) 2022 Rendered Text

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.


