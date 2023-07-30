# Speech-To-Text-guide



**Workshop Guide Outline**

**1. Introduction to CLI Applications using `cli` package in Go**

In this part, we'll cover how to set up a Command Line Interface (CLI) application using the "cli" package in Go. We will start with a simple program that accepts the "get" command for fetching transcriptions using a supplied YouTube link. We'll just print back the youtube link for now.

**Task**: Create a new Go file in a new git branch named "cli-setup". Now, write a basic program that defines a CLI application which accepts the "get" command.

**Depedencies**
`go get github.com/urfave/cli/v2`


**Code Snippets**

```go
package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2" //this imports the cli package
)

func main() {
	app := &cli.App{
		Name:  "ytt",
		Usage: "Transcribe YouTube videos",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get transcriptions by ytlink",
				Action: func(c *cli.Context) error {
					//Print YouTube link
					youtubelink := c.Args().Get(0)
					if youtubelink == "" {
						return cli.NewExitError("Please provide a YouTube link", 1)
				    }
					fmt.Println("YouTube link:", youtubelink)
					return nil
				},
			},
		},
	
   	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

```

**2. Downloading YouTube Video and Audio Streams**

Now, we'll use the `youtube/v2` package to fetch video details and download the audio stream for the supplied YouTube link. 

**Task**: In a new git branch named "downloading-video", append code to your program from Step 1 for downloading a video using a YouTube link. Add a "download" command to ur CLI for testing this functionality.

**Depedencies**

`go get github.com/kkdai/youtube/v2`

**Code Snippets**

```go
//youtube.go
package main

import (
	"context"
	"fmt"

	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
)

func YoutubeDL(ytlink string) error {
	client := youtube.Client{
		Debug: true,
	}
	ctx := context.Background()


	video, err := client.GetVideoContext(ctx, ytlink)
	if err != nil {
		return fmt.Errorf("Error getting video: %w", err)
	}
	downloader := downloader.Downloader{Client: client, OutputDir: "./"}

	return downloader.DownloadComposite(ctx, "", video, "hd1080", "mp4")

}
```

```go
//main.go
err := YoutubeDL(youtubelink)
if err != nil {
    return err
}
```
**3. Audio Conversion Using FFmpeg**

Learn the use of the FFmpeg tool to convert the audio file to a suitable format for transcription. This module may also cover how to install FFmpeg, if not already available on participants' machines. 

**Task**: In a new git branch "ffmpeg-conversion", append the current program you have with the file conversion using FFmpeg code. Add a "convert" command to your CLI for testing this step.

**Code Snippets**

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func ConvertFile(inputFile string, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-vn", "-ac", "1", "-ar", "16000", "-codec:a", "pcm_s16le", "-f", "wav", outputFile)
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	return nil
}
```
```go
//main.go
fmt.Println("Audio file:", audieofilename)
err = ConvertFile(audieofilename+".mp4", audieofilename+".wav")
if err != nil {
    return err
}
```
**4. Building Docker Images and running Docker Containers**

This part of the workshop dives into Docker, as Docker will simplify the setup process for the application development and distribution.

**Understanding the Dockerfile**

```Dockerfile
# Use the official Golang base image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Install whisper
RUN git clone https://github.com/ggerganov/whisper.cpp.git &&\
    cd whisper.cpp && make &&\
    make libwhisper.so libwhisper.a &&\
    cp whisper.h /usr/local/include &&\
    cp ggml.h /usr/local/include &&\
    cp libwhisper.a /usr/local/lib &&\
    cp libwhisper.so /usr/local/lib &&\
    cd ..

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download all necessary dependencies
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Go app
RUN go build -o ytt

# Install ffmpeg
RUN apt-get update && apt-get install -y ffmpeg

# execute cli help to check if everything is ok
RUN ./ytt -h


# Run the compiled binary with a default command
ENTRYPOINT ["/app/ytt"]
```
    
We will go through the Dockerfile and explain each directive:

- `FROM` specifies the base image. Here, we are using the official Golang image.
- `WORKDIR` sets the current working directory inside the Docker container.
- `RUN` is used to execute commands. In this Dockerfile, we use it to download and compile the whisper library, download and install FFmpeg, and build our Go application.
- `COPY` copies new files or directories from "<source>" and adds them to the filesystem of the image at the path "<destination>".
- `ENTRYPOINT` sets the command and parameters that will be executed first when a container is run.

Before we build the image check if the tiny model is present. If not download it. 
 ```bash
sh download_model tiny.en
```

**Task: Building the Docker Image**
    
We will use Docker build command to create a Docker image from the Dockerfile. Here's the command we'll use:

```bash
docker buildx build --platform linux/amd64 -t ytt-amd64 --load -f Dockerfile .
```

This command tells Docker to build an image using the Dockerfile in the current directory (the "." at the end). The flag `--platform linux/amd64` specifies the platform the image is being built for. The `-t ytt-amd64` flag tags the image with the name "ytt-amd64". The `--load` flag tells Docker to load the built image into Docker's locally accessible image store.
    
Run this command so that Docker can build your image.

**Task: Running the Docker Container**

After successfully building the Docker image, we can now create and run a Docker container from it.

```bash
docker run --platform linux/amd64 -v "$(pwd)":/data -it --rm ytt-amd64 get "JzPfMbG1vrE"
```

The `docker run` command creates and starts a Docker container. The `--platform linux/amd64` flag specifies the platform of the container. The `-v "$(pwd)":/data` flag mounts the current directory from the host into the container at "/data". The `-it` flag ensures that we can interact with the container via the terminal. `--rm` tells Docker to automatically clean up the container and remove the file system when the container exits. The `get "JzPfMbG1vrE"` at the end of the command is the command arguments that will be passed to the ENTRYPOINT command inside the Kubernetes container.

Run this command to interact with your image. Check to see that it works as expected, given its output, and familiarize yourself with the process of building and running Docker containers. This should conclude our workshop.
