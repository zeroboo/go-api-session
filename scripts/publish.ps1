$VERSION="v1.0.5"

go mod tidy
go test 
# git commit -m "New version"
git tag $VERSION
git push origin $VERSION
go list -m github.com/zeroboo/go-api-session@$VERSION