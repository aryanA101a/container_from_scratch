# Container From Scratch
A container written in go from scratch.

### IMPORTANT
- #### To use this program you have to be in sudo mode
- #### Download Ubuntu Filesystem
  1.  Pull the ubuntu image from docker and create a container, then use docker export.
  2. `sudo docker export ub > ufs.tar`
  3. `tar xvf ufs.tar`
  
### Environment
Ubuntu/Fedora

### Requirements
- Go
- Docker

### Usage
`go run main.go run /bin/bash`
