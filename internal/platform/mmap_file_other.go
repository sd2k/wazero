//go:build !linux

package platform

import (
	"errors"
	"os"
)

// MapsCodeSegmentsFromFiles reports whether MapCodeSegmentFromFile is
// implemented on this platform.
const MapsCodeSegmentsFromFiles = false

// MapCodeSegmentFromFile is only implemented on Linux, see mmap_linux.go.
// Elsewhere, callers copy the code into MmapCodeSegment instead.
func MapCodeSegmentFromFile(*os.File, int64, int) ([]byte, error) {
	return nil, errors.ErrUnsupported
}
