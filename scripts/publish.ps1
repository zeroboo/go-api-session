$VERSION="v0.0.1"

go mod tidy
go test 
git commit -m ""
git tag $VERSION
git push origin $VERSION
go list -m github.com/zeroboo/go-api-session@v0.1.0