package photos

import (
	"crypto/md5"
	"fmt"
	"os"
	"path"
)

var (
	BaseStoragePath = "photos_storage"
)

func CreateUserStorageDirectory(email string) (storagePath string, err error) {
	sum := md5.Sum([]byte(email))
	storagePath = path.Join(BaseStoragePath, fmt.Sprintf("%x", sum[:]))
	err = os.MkdirAll(BaseStoragePath, 0777)
	return
}
