TestProcessAlbumDirectory {
    *p1="/tempZone/home/rods/gopher.jpg";
    *p2="/home/john/irods-ugm/irods-ugm-2017/go-microservice/iRODS-UGM-Demo-eda4ee05c91f.json";

    msiprocess_album_directory(*p1, *p2);
}

INPUT null
OUTPUT ruleExecOut