<h1>J.A.R.V.I.S</h1>


<h1>Installation</h1>
<h2>Step 1:</h2>

Download [Whisper](https://github.com/openai/whisper#setup)

<h2>Step 2:</h2>

Download [AWS-cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)

<h2>Step 3:</h2>

Create an [Access key](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_root-user_manage_add-key.html) and authenticate


```shell
$ aws configure
AWS Access Key ID [None]: AKIAIOSFODNN7EXAMPLE
AWS Secret Access Key [None]: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
Default region name [None]: us-west-2
Default output format [None]: json
```

<h2>Step 4:</h2>

Create an OpenAI [AUTH Token](https://platform.openai.com/api-keys)

<h2>Step 5:</h2>

Downnload J.A.R.V.I.S.

```sh

$ wget https://raw.githubusercontent.com/drpaij0se/J.A.R.V.I.S./main/init.sh ; bash init.sh
```

<h1>Installation with Docker</h1>

<h2>Step 1:</h2>
Clone the repository

```sh
git clone https://github.com/drpaij0se/J.A.R.V.I.S ; cd J.A.R.V.I.S/ 
```

<h2>Step 2:</h2>
Configure the Dockerfile with your credentials

```Dockerfile
RUN aws configure set aws_access_key_id $AWS_ACCESS_KEY_ID
RUN aws configure set aws_secret_access_key $AWS_SECRET_ACCESS_KEY
RUN aws configure set default.region $AWS_DEFAULT_REGION
```

<h2>Step 3:</h2>
Build the image

```sh
docker build --tag jarvis . 
```

<h2>Step 4:</h2>
Run the container

```sh
docker run -ti --rm --device /dev/snd:/dev/snd jarvis  
```

<h1>Configuration</h1>

- The configuration file is located in: `$HOME/.config/jarvis/jarvis.yml`

```yml
auth: sk-**** # OpenAI Key
language: Spanish 
model: gpt-3.5-turbo
voiceId: Lupe
```
