package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/fsouza/go-dockerclient"
)

// Job is a job
type Job struct {
	ID         int
	StepID     int
	BranchName string
	CommitHash string
	step       Step
}

// Cmd is a command to run in the job
type Cmd []string

// Run runs the Job
func (job *Job) Run() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Printf("Error creating client %s", err)
		return err
	}
	log.Printf("%+v", client)
	err = client.Ping()
	if err != nil {
		log.Printf("Failed to ping: %s", err)
		return err
	}

	commands := job.step.Cmds

	var prevImage *docker.Image

	for i, cmd := range commands {
		imageName := job.step.ImageName
		if prevImage != nil {
			imageName = prevImage.ID
		}
		config := docker.Config{
			Cmd:   cmd,
			Image: imageName,
		}

		createOpts := docker.CreateContainerOptions{
			Name:   fmt.Sprintf("%d_%d", job.ID, i),
			Config: &config,
		}

		container, err := client.CreateContainer(createOpts)
		if err != nil {
			log.Printf("Failed to create container: %s", err)
			return err
		}

		log.Printf("%+v", container)

		hostConfig := &docker.HostConfig{}
		err = client.StartContainer(container.ID, hostConfig)
		if err != nil {
			log.Printf("Failed to start container: %s", err)
			return err
		}
		log.Println("Container started, waiting for it to finish")
		exitCode, err := client.WaitContainer(container.ID)
		if err != nil {
			log.Printf("Error waiting for container: %s", err)
			return err
		}
		log.Printf("Container exited with code %d", exitCode)
		log.Printf("Getting logs")
		stdOut := bytes.NewBuffer([]byte{})
		stdErr := bytes.NewBuffer([]byte{})
		err = client.Logs(docker.LogsOptions{
			Container:    container.Name,
			OutputStream: stdOut,
			ErrorStream:  stdErr,
			Stdout:       true,
			Stderr:       true,
		})

		if err != nil {
			log.Printf("Error getting logs: %s", err)
			return err
		}

		log.Printf("Stdout: %s", string(stdOut.Bytes()))
		log.Printf("Stderr: %s", string(stdErr.Bytes()))

		image, err := client.CommitContainer(docker.CommitContainerOptions{
			Container:  container.ID,
			Repository: "temp",
			Tag:        strconv.Itoa(i),
		})
		if err != nil {
			log.Printf("Error committing image: %s", err)
			return err
		}
		log.Printf("%+v", image)
		prevImage = image
	}
	return nil
}
