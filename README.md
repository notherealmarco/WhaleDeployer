## WhaleDeployer
*a container, to deploy them all!*

<img src="https://www.docker.com/wp-content/uploads/2022/03/horizontal-logo-monochromatic-white.png" height="35" />

**WhaleDeployer is a CI/CD tool** that automates the deployment process of your **Docker applications**.

When a new **commit** is available on the Git repository, WhaleDeployer will **pull** the new version of the application, **build** the Docker image and **deploy** it on your machine.

Everything can be set up from a **web interface**, and then you can forget about it.

### Who is it meant for?

WhaleDeployer best suits your needs if:

- You want to **automate** the deployment process on your machine
- You **build** your Docker images **on your machine**
- You want to manage **SSH deploy keys** in a painless way
- *You want this all to be without any human intervention*

## To do list

- Built-in support for authentication
- Ability to reset SSH keys
- Improve code readability
- *Your idea(s) here!*

## Screenshots

<img src="https://i.imgur.com/VKNMAgx.png" alt="Home" width="300">
<img src="https://i.imgur.com/GxaQHfa.png" alt="Home" width="300">
<img src="https://i.imgur.com/c5WSkuo.png" alt="Home" width="300">

## Installation

Edit the given `docker-compose.yml` (if needed) to change ports and volumes.

By default, the web interface and the API will be available on port `3333`, and two persistent voluemes will be created to store the configuration and your repositories.

Then run:
```bash
docker-compose up -d
```

The application should be up and running on `http://<your server IP>:3333`.

## Usage

### Configure a new project

First, make sure your repository contains a `docker-compose.yml` file.

1. Click on the `+` button in the top right corner
2. Fill the form with the following information:
    - **Name**: the name of your project
    - **Path**: the path where the Git repository will be cloned
    - **Git URI**: the URI *(SSH or HTTPS)** of your Git repository
    - **Branch**: the branch that will be pulled
    - **Dockerfile name**: optional, if you want to use `docker build` instead of `docker-compose build`
    - **Image name**: the name of the Docker image *(only if you specify a Dockerfile)*
    - **Image tag**: the tag of the Docker image *(only if you specify a Dockerfile)*
    - **I have a private repository**: check if you want to use SSH* with a deploy key

_**(*)** For now SSH only works if you generate and use a deploy key. If the repository is public, you can temporarily use HTTPS._

3. Save and deploy!
4. Configure a GitHub Action to send a POST request to the following webhook: 
`http://<your server domain and port>/api/projects/<project name>`

### Configure the GitHub Action

Here's an example using Basic authentication provided by a reverse proxy (see later):
```yaml
name: WhaleDeployer
on:
  push:
    branches: [ "main" ]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: curl
        uses: wei/curl@v1
        env:
          TOKEN: ${{ secrets.AUTH }}
        with:
          args: "-X POST https://<your server>/api/projects/mysupercoolproject -H 'Authorization: Basic ${{ secrets.AUTH }}' --max-time 900"
```
## FAQs

### Can I use a reverse proxy?

Sure! You should use it as the application does not (yet) support authentication and HTTPS.

### How can I set up the authentication?

For now, you can configure an HTTP Basic Auth in your reverse proxy. _You should really do it!_

## License notice

This projects uses code from Enrico Bassetti's [Fantastic coffee decaffeinated](https://github.com/sapienzaapps/fantastic-coffee-decaffeinated) project, which is licensed under the below MIT License.

```text
MIT License

Copyright (c) 2022 Enrico Bassetti

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
```
