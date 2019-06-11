# tcs-cli

**tcs-cli** is a command line tool which interacts with telestream cloud API. Tool is written in go and uses auto-generated [telestream cloud go sdk](github.com/Telestream/telestream-cloud-go-sdk).
	
## configure credentials

Before any interraction with telestream cloud - client's **X-Api-Key** must be set:

```sh
$ tcs configure -api_key CLIENTS_X_API_KEY
```

After execution of this line, user credentials are saved in .tcs-credentials file.

## flip service

### factories

#### - factories list 

To list all client's factories, call:

```sh
$ tcs flip factories list
```

#### - factories describe

To print given factory description:

```sh
$ tcs flip factories describe -factory_id FACTORY_ID
```

### profiles 

#### - profiles list

To list all client's profiles in given factory:

```sh
$ tcs flip profiles list -factory_id FACTORY_ID
```

#### - profiles describe 

To print given profile description:

```sh
$ tcs flip profiles describe -factory_id FACTORY_ID -profile_name PROFILE_NAME
```
or

```sh
$ tcs flip profiles describe -factory_id FACTORY_ID -profile_id PROFILE_ID
```


#### - profiles create 

To create new profile in given factory with given preset name (also all profile parameters are available):

```sh
$ tcs flip profiles create -factory_id FACTORY_ID -preset_name PRESET_NAME -width WIDTH -height HEIGHT ...
```

#### - profiles delete

To delete given profile in given factory:

```sh
$ tcs flip profiles delete -factory_id FACTORY_ID -profile_id PROFILE_ID
```

#### - profiles update 

To update existing profile in given factory with given preset name (also all profile parameters are available):

```sh
$ tcs flip profiles update -factory_id FACTORY_ID -preset_name PRESET_NAME -width WIDTH -height HEIGHT ...
```

### videos

#### - videos list

To list all client's videos in given factory:

```sh
$ tcs flip videos list -factory_id FACTORY_ID
```

#### - videos describe 

To print given video description:

```sh
$ tcs flip videos describe -factory_id FACTORY_ID -video_id VIDEO_ID
```

#### - videos create 

To create new video in given factory with given source url (also all video parameters are available):

```sh
$ tcs flip profiles create -factory_id FACTORY_ID -source_url SOURCE_URL -profiles PROFILE1 PROFILE2 ...
```

#### - video delete

To delete given video in given factory:

```sh
$ tcs flip profiles delete -factory_id FACTORY_ID -video_id VIDEO_ID
```

#### - video cancel

To cancel given video in given factory:

```sh
$ tcs flip profiles cancel -factory_id FACTORY_ID -video_id VIDEO_ID
```

### encodings

#### - encodings list

To list all client's encodings in given factory:

```sh
$ tcs flip encodings list -factory_id FACTORY_ID
```

#### - encodings describe

To print description of given encoding in given factory:

```sh
$ tcs flip encodings describe -factory_id FACTORY_ID -encoding_id ENCODING_ID
```

#### - encodings delete

To delete given encoding in given factory:

```sh
$ tcs flip encodings delete -factory_id FACTORY_ID -encoding_id ENCODING_ID
```

#### - encodings cancel

To cancel given encoding in given factory:

```sh
$ tcs flip encodings cancel -factory_id FACTORY_ID -encoding_id ENCODING_ID
```

#### - encodings signed_urls

To list encoding signed urls in given factory:

```sh
$ tcs flip encodings signed-urls -factory_id FACTORY_ID -encoding_id ENCODING_ID
```

## tts service

### projects

#### - projects list

To list all client's project:

```sh
$ tcs tts projects list
```

#### - projects describe

To print description of given project:

```sh
$ tcs tts projects describe -project_id PROJECT_ID
```

#### - project create 

To create new profile with given name, language and description (also all project parameters are available):

```sh
$ tcs tts project create -name NAME -description DESCRIPTION -language LANGUAGE ...
```

#### - projects delete

To delete given project

```sh
$ tcs tts projects delete -project_id PROJECT_ID
```

#### - project update

To create new profile with given name, language and description (also all project parameters are available):

```sh
$ tcs tts project update -id ID -description DESCRIPTION -language LANGUAGE ...
```

### jobs

#### - jobs list

To list all client's jobs in given project:

```sh
$ tcs tts jobs list -project_id PROJECT_ID
```

#### - jobs describe

To print description of given job in given project:

```sh
$ tcs tts jobs describe -project_id PROJECT_ID -job_id JOB_ID
```

#### - jobs create

To create job in given project with given source url (also all job parameters are available):

```sh
$ tcs tts jobs create -project_id PROJECT_ID -source_url SOURCE_URL ...
```

#### - jobs delete

To delete job in given project:

```sh
$ tcs tts jobs delete -project_id PROJECT_ID -job_id JOB_ID
```

#### - jobs result

To print result of given job in given project:

```sh
$ tcs tts jobs result -project_id PROJECT_ID -job_id JOB_ID
```

#### - jobs outputs

To print outputs of given job in given project:

```sh
$ tcs tts jobs outputs -project_id PROJECT_ID -job_id JOB_ID
```

### corpora


#### - corpora list

To list all client's corpora in given project:

```sh
$ tcs tts corpora list -project_id PROJECT_ID
```

#### - corpora describe

To print description of given corpus in given project:

```sh
$ tcs tts corpora describe -project_id PROJECT_ID -corpus_name CORPUS_NAME
```

### - copora create

To create corpus in given project with given corpus name and corpus body:

```sh
$ tcs tts corpora create -project_id PROJECT_ID -corpus_name CORPUS_NAME - corpus_body CORPUS_BODY ...
```

#### - copora delete

To delete corpus in given project:

```sh
$ tcs tts corpora delete -project_id PROJECT_ID -corpus_name CORPUS_NAME
```

## how to pass additional header key and value

```sh
$ tcs ... -header_key HEADER_KEY -header_value HEADER_KEY_VALUE
```
