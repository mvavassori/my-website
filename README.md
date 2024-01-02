# My Website

This is the source code for my Go-based website. It's dockerized for easy setup and deployment.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have installed Docker on your machine. Get Docker

## Running the Application with Docker

Follow these steps to run the application using Docker:

### Step 1: Clone the Repository

First, clone this repository to your local machine:

bashCopy code

`git clone https://github.com/mvavassori/my-website.git cd my-website`

### Step 2: Build the Docker Image

Build the Docker image using the following command:

bashCopy code

`docker build -t my-website .`

This command builds a Docker image named `my-website` based on the instructions in the `Dockerfile`.

### Step 3: Run the Docker Container

Now, you can run your application in a Docker container.

If you have a `.env` file with environment variables:

bashCopy code

`docker run --env-file .env -p 8080:8080 my-website`

If you don't have a `.env` file, simply run:

bashCopy code

`docker run -p 8080:8080 my-website`

These commands start a Docker container based on the `my-website` image. The `-p 8080:8080` option maps port 8080 of the container to port 8080 on your host machine, making the application accessible via `http://localhost:8080`.

### Step 4: Accessing the Application

Open a web browser and go to `http://localhost:8080` to access the application.

## Stopping the Application

To stop the Docker container, you can press `CTRL+C` in the terminal if you're running it in the foreground. If it's running in the background, use:

bashCopy code

`docker ps`

to find the container ID and then:

bashCopy code

`docker stop [container_id]`

to stop it.

## Additional Information

- Ensure you have the latest version of Docker installed for a smooth experience.
- For any additional configuration or environment-specific setups, refer to the application's main documentation or contact the development team.
