# iRODS UGM 2017 Demo

This repository contains resources for John Jacquay's iRODS UGM 2017 talk/demo. The following instructions are for building and testing the Golang microservice demo.

Two microservices are included: `msiextract_image_metadata` and `msibasic_example`.

`msiextract_image_metadata` extracts metadata information from image files using [Google's Cloud Vision API](https://cloud.google.com/vision/) and a [Golang EXIF package](https://github.com/rwcarlsen/goexif).

`msibasic_example` accepts a two column CSV string as the first argument, and returns an output `KeyValPair_MS_T` data structure in the second argument.

## Vagrant + VirtualBox Install (easy mode)

```
$ git clone https://github.com/jjacquay712/irods-ugm-2017.git
$ cd irods-ugm-2017
$ vagrant up && vagrant ssh
$ sudo service irods start
```

## Manual Build / Install Microservices (hard mode)

1. **Install iRODS 4.2.1 and dependencies (In my case, CentOS 7)**
```
$ sudo rpm --import https://packages.irods.org/irods-signing-key.asc
$ wget -qO - https://packages.irods.org/renci-irods.yum.repo | sudo tee /etc/yum.repos.d/renci-irods.yum.repo
$ yum install irods-externals* irods-runtime irods-devel
```

2. [Install Golang](https://golang.org/doc/install)

3. **Fetch project files**
```
$ git clone https://github.com/jjacquay712/irods-ugm-2017.git
```

4. **Setup $PATH**
```
$ export PATH=$PATH:/opt/irods-externals/cmake3.5.2-0/bin
$ export PATH=$PATH:/opt/irods-externals/clang3.8-0/bin
```

5. **Build and install**
```
$ cd irods-ugm-2017/go-microservice
$ mkdir build && cd build
$ cmake .. && make
$ sudo make install
```

6. **You're all set!**

## Usage of `msiextract_image_metadata`

1. Edit iRODS `core.re`
```
$ sudo vi /etc/irods/core.re
```

2. Add the following contents to `core.re` file (if you chose manual install, not required for Vagrant):
```
acPostProcForPut {
	# Second parameter enables/disables gzip compression
	msiextract_image_metadata($objPath, 0);
}
```

3. Run `iput` on image file from root `irods-ugm-2017/` repo directory:
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

## Usage of `msibasic_example`

1. Run `irule -F go-microservice/msibasic_example/test.r` from root `irods-ugm-2017/` repo directory.

test.r:
```
TestBasicExample {
    msibasic_example("keytest,valuetest", *outKVP);
    msiPrintKeyValPair("stderr", *outKVP)
}

INPUT null
OUTPUT ruleExecOut
```

Output:
```
keytest = valuetest
```

### Testing `msibasic_example`

From `irods-ugm-2017/` root repo directory:
```
$ go-microservice/run_tests.sh 
```

Output:
```
=== RUN   TestBasicExample
--- PASS: TestBasicExample (0.00s)
	msibasic_example_test.go:19: Success!
PASS
ok  	command-line-arguments	0.023s
```

## Developer Resources

* [msibasic_example Microservice Source Code](/go-microservice/msibasic_example/msibasic_example.go)
* [msibasic_example Microservice Testing Code](/go-microservice/msibasic_example/msibasic_example_test.go)
* [msiextract_image_metadata Microservice Source Code](/go-microservice/msiextract_image_metadata/msiextract_image_metadata.go)
* [GoRODS/msi Package Documentation](https://godoc.org/github.com/jjacquay712/GoRODS/msi)


