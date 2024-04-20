package mode

import "time"

type SimpleFileInfo struct {
	Name           string
	Size           int64
	HashCode       string
	LastUpdateTime time.Time
}
