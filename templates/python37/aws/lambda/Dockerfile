FROM public.ecr.aws/lambda/python:3.7

# Copy local code to the container image.
WORKDIR /var/task
COPY . ./

# Install the depenedencies
RUN bin/install_dependencies.sh

CMD ["main.{{.FunctionName}}"]
