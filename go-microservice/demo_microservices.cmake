
set(
  GOMICROSERVICES
  msiextract_image_metadata
  msibasic_example
  )

set(
  CPACK_RPM_EXCLUDE_FROM_AUTO_FILELIST_ADDITION
  /etc/irods
  /usr/lib/irods
  /usr/lib/irods/plugins
  /usr/lib/irods/plugins/microservices
)

set(
  IRODS_COMPILE_DEFINITIONS_DEMO_MICROSERVICES
  RODS_SERVER
  ENABLE_RE
  )

set(ENV{GOPATH} ${CMAKE_SOURCE_DIR})

set(ENV{CXX} "/opt/irods-externals/clang3.8-0/bin/clang++")
set(ENV{CC} "/opt/irods-externals/clang3.8-0/bin/clang")

# Install golang dependencies
execute_process(COMMAND go get github.com/jjacquay712/GoRODS
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})

execute_process(COMMAND go get golang.org/x/net/context
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})

execute_process(COMMAND go get cloud.google.com/go/vision/apiv1
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})

execute_process(COMMAND go get golang.org/x/text/language
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})

execute_process(COMMAND go get cloud.google.com/go/translate
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})

execute_process(COMMAND go get github.com/rwcarlsen/goexif/exif
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})

execute_process(COMMAND go get github.com/rwcarlsen/goexif/mknote
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})

execute_process(COMMAND go get github.com/rwcarlsen/goexif/tiff
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})


foreach(GOMICROSERVICE ${GOMICROSERVICES})

  message(STATUS "Building golang static archive for " ${GOMICROSERVICE})

  execute_process(COMMAND go build -buildmode=c-archive ${GOMICROSERVICE}.go
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}/${GOMICROSERVICE})

  add_library(
    ${GOMICROSERVICE}
    MODULE
    ${CMAKE_SOURCE_DIR}/${GOMICROSERVICE}/lib${GOMICROSERVICE}.cpp
    )

  target_include_directories(
    ${GOMICROSERVICE}
    PRIVATE
    ${IRODS_INCLUDE_DIRS}
    ${IRODS_EXTERNALS_FULLPATH_BOOST}/include
    ${IRODS_EXTERNALS_FULLPATH_JANSSON}/include
    )

  target_link_libraries(
    ${GOMICROSERVICE}
    PRIVATE
    pthread
    irods_server
    irods_client
    irods_common
    ${IRODS_EXTERNALS_FULLPATH_BOOST}/lib/libboost_filesystem.so
    ${IRODS_EXTERNALS_FULLPATH_BOOST}/lib/libboost_system.so
    ${CMAKE_SOURCE_DIR}/${GOMICROSERVICE}/${GOMICROSERVICE}.a
    ${CMAKE_DL_LIBS}
    )

  target_compile_definitions(${GOMICROSERVICE} PRIVATE ${IRODS_COMPILE_DEFINITIONS_DEMO_MICROSERVICES} ${IRODS_COMPILE_DEFINITIONS} BOOST_SYSTEM_NO_DEPRECATED)
  set_property(TARGET ${GOMICROSERVICE} PROPERTY CXX_STANDARD ${IRODS_CXX_STANDARD})

  install(
    TARGETS
    ${GOMICROSERVICE}
    LIBRARY
    DESTINATION /usr/lib/irods/plugins/microservices
    )

endforeach()

