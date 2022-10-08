# Container From Scratch
A container written in go from scratch.

### IMPORTANT
- #### To use this program you have to be in sudo mode
- #### Download Ubuntu Filesystem
  1. `sudo docker pull ubuntu:latest`
  2. `sudo docker create --name ub ubuntu`
  3. `sudo docker export ub > ufs.tar`
  4. `tar xvf ufs.tar`
  
### Environment
Ubuntu/Fedora

### Requirements
- Go
- Docker

### Usage
`go run main.go run /bin/bash`
