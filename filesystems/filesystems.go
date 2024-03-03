package filesystems

import "time"

// FS is the interface that wraps the basic methods for a filesystem
// In order to satisfy this interface, a filesystem must implement the following methods:
type FS interface {
	Put(filename, folder string) error
	Get(destination string, items ...string) error
	List(prefix string) ([]Listing, error)
	Delete(itemsToDelete []string) bool
}

// Listing is a struct that contains the information of a file or folder
type Listing struct {
	Etag         string
	LastModified time.Time
	Key          string
	Size         float64
	IsDir        bool
}
