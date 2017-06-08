# iRODS UGM 2017 Demo

This repository contains resources for John Jacquay's iRODS UGM 2017 talk/demo. The following instructions are for building and testing the Golang microservice demo.

## Build / Install

1. **Install iRODS and dependencies (In my case, CentOS 7)**
```
$ sudo rpm --import https://packages.irods.org/irods-signing-key.asc
$ wget -qO - https://packages.irods.org/renci-irods.yum.repo | sudo tee /etc/yum.repos.d/renci-irods.yum.repo
$ yum install irods-externals* irods-runtime irods-devel
```

2. **Fetch project files**
```
$ git clone https://github.com/jjacquay712/irods-ugm-2017.git
```

3. **Setup $PATH**
```
$ export PATH=$PATH:/opt/irods-externals/cmake3.5.2-0/bin
$ export PATH=$PATH:/opt/irods-externals/clang3.8-0/bin
```

4. **Build and install**
```
$ cd irods-ugm-2017/go-microservice
$ mkdir build && cd build
$ cmake .. && make
$ make install
```

5. **Configure iRODS `core.re`**
```
$ sudo vi /etc/irods/core.re
```

Add the following contents to `core.re` file:
```
acPostProcForPut {
	msiextract_image_metadata($objPath);
}
```

6. **You're all set!**

## Usage

From root `irods-ugm-2017/` repo directory:
```
$ iput gopher.jpg
```

![Gopher Picture](/gopher.jpg?raw=true "Gophers are cool")

```
$ imeta ls -d gopher.jpg
AVUs defined for dataObj gopher.jpg:
attribute: tags_english
value: mammal,vertebrate,squirrel,wildlife,fauna,whiskers,prairie dog,marmot,rodent,fox squirrel
units: 
----
attribute: tags_dutch
value: zoogdier,gewerveld,eekhoorn,dieren in het wild,fauna,bakkebaarden,prairiehond,marmot,knaagdier,Vos eekhoorn
units: 
```

## Developer Resources

* [msiextract_image_metadata Microservice Source Code](/go-microservice/msiextract_image_metadata/msiextract_image_metadata.go)
* [GoRODS/msi Package Documentation](https://godoc.org/github.com/jjacquay712/GoRODS/msi)

