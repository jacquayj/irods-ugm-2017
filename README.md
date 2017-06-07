# iRODS UGM 2017 Demo

This repository contains resources for John Jacquay's iRODS UGM 2017 talk/demo. The following instructions are for building and testing the Golang microservice demo.

## Build / Install

1. Install iRODS and dependencies (In my case, CentOS 7)
```
$ sudo rpm --import https://packages.irods.org/irods-signing-key.asc
$ wget -qO - https://packages.irods.org/renci-irods.yum.repo | sudo tee /etc/yum.repos.d/renci-irods.yum.repo
$ yum install irods-externals* irods-runtime irods-devel
```

2. Fetch project files
```
git clone https://github.com/jjacquay712/irods-ugm-2017.git
```

3. Build and install
```
$ cd irods-ugm-2017/go-microservice
$ mkdir build && cd build
$ cmake .. && make
$ make install
```

## Usage

Working directory: `irods-ugm-2017/go-microservice/build`

```

```
