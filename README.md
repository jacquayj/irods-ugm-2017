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

4. Configure iRODS with rule set
```
$ sudo vi /etc/irods/server_config.json 
```

Add `extract_image_metadata` to `re_rulebase_set` JSON array:
```
"re_rulebase_set": [
    "core", "extract_image_metadata"
],
```

5. Restart iRODS
```
$ sudo service irods restart
```

## Usage

```
$ iput gopher.jpg
$ imeta ls -d gopher.jpg
AVUs defined for dataObj gopher.jpg:
attribute: tags_english
value: nature,mammal,vertebrate,wildlife,fauna,grass,whiskers,domestic rabbit,prairie dog,squirrel
units: 
----
attribute: tags_dutch
value: natuur,zoogdier,gewerveld,dieren in het wild,fauna,gras,bakkebaarden,Huiselijk konijn,prairiehond,eekhoorn
units: 
```
