package filestorageimpl

import "github.com/rfancn/airpress/injection"

func init() {
	injection.Provide(
		NewMinIO,
		NewLocalFileStorage,
		NewAliyun,
	)
}
