set -e

VERSION=$(grep current_version .bumpversion.cfg | cut -d " " -f 3)
PROJECT_NAME=$(basename $(pwd))

# build project
cd Godeps/_workspace
mkdir -p src/github.com/trustedanalytics/
cd src/github.com/trustedanalytics/
ln -s ../../../../.. $PROJECT_NAME
cd ../../../../..

GOPATH=`godep path`:$GOPATH go test ./...
godep go build

rm Godeps/_workspace/src/github.com/trustedanalytics/$PROJECT_NAME

# assemble the artifact
PACKAGE_CATALOG=${PROJECT_NAME}-${VERSION}

# prepare build manifest
echo "commit_sha=$(git rev-parse HEAD)" > build_info.ini

# create zip package
zip -r ${PROJECT_NAME}-${VERSION}.zip * -x ${PROJECT_NAME}-${VERSION}.zip


echo "Zip package for $PROJECT_NAME project in version $VERSION has been prepared."
