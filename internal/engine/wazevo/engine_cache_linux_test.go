package wazevo

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"unsafe"

	"github.com/tetratelabs/wazero/internal/platform"
	"github.com/tetratelabs/wazero/internal/testing/require"
)

// On Linux, a file cache hit maps the code from the cache file rather than
// copying it into anonymous memory.
func TestDeserializeCompiledModule_mapsCodeFromFile(t *testing.T) {
	code := []byte{1, 2, 3, 4, 5}
	entry, err := io.ReadAll(serializeCompiledModule(testVersion, &compiledModule{
		executables:     &executables{executable: code},
		functionOffsets: []int{0},
	}))
	require.NoError(t, err)
	f := cacheFile(t, entry)

	cm, staleCache, err := deserializeCompiledModule(testVersion, f)
	require.NoError(t, err)
	require.False(t, staleCache)
	require.Equal(t, code, cm.executable)
	defer func() { require.NoError(t, platform.MunmapCodeSegment(cm.executable)) }()

	require.Equal(t, f.Name(), mappingPath(t, uintptr(unsafe.Pointer(&cm.executable[0]))))
}

// mappingPath returns the file backing the mapping that contains addr, from
// /proc/self/maps.
func mappingPath(t *testing.T, addr uintptr) string {
	maps, err := os.ReadFile("/proc/self/maps")
	require.NoError(t, err)
	for _, line := range strings.Split(string(maps), "\n") {
		var start, end uintptr
		if _, err := fmt.Sscanf(line, "%x-%x", &start, &end); err != nil || addr < start || addr >= end {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			return ""
		}
		return fields[5]
	}
	t.Fatalf("no mapping contains %#x", addr)
	return ""
}
