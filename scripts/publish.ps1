$VERSION="v1.6.1"

go mod tidy
go test 
# git commit -m "New version"
git tag $VERSION
git push origin $VERSION
go list -m github.com/zeroboo/go-api-session@$VERSION