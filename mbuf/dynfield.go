package mbuf

/*
#include <rte_config.h>
#include <rte_mbuf.h>
#include <rte_mbuf_dyn.h>

// Read a uint64 from a dynamic field at the given byte offset in the mbuf.
static inline uint64_t go_mbuf_dynfield_get_uint64(struct rte_mbuf *m, int offset)
{
	return *RTE_MBUF_DYNFIELD(m, offset, uint64_t *);
}

// Register the standard RX timestamp dynamic field.
// On success, *offset receives the byte offset into rte_mbuf where the
// NIC stores the HW timestamp (as rte_mbuf_timestamp_t / uint64).
// Returns 0 on success, negative errno on failure.
static inline int go_mbuf_dyn_rx_timestamp_register(int *offset, uint64_t *flag)
{
	return rte_mbuf_dyn_rx_timestamp_register(offset, flag);
}
*/
import "C"

import "fmt"

// DynfieldGetUint64 reads a uint64 value from a dynamic metadata field
// at the given byte offset inside the mbuf structure.
// The offset is obtained from DynRxTimestampRegister (or a similar
// rte_mbuf_dynfield_register call).
func (m *Mbuf) DynfieldGetUint64(offset int) uint64 {
	return uint64(C.go_mbuf_dynfield_get_uint64(mbuf(m), C.int(offset)))
}

// DynRxTimestampRegister registers the standard DPDK RX timestamp
// dynamic field ("rte_dynfield_timestamp") and the associated dynamic
// flag ("rte_dynflag_rx_timestamp").
//
// On success it returns (fieldOffset, rxFlag, nil).
// fieldOffset is the byte offset within rte_mbuf where the NIC writes
// the hardware timestamp.  rxFlag is a bitmask for ol_flags that
// indicates the field is valid.
func DynRxTimestampRegister() (fieldOffset int, rxFlag uint64, err error) {
	var off C.int
	var flag C.uint64_t
	rc := C.go_mbuf_dyn_rx_timestamp_register(&off, &flag)
	if rc < 0 {
		return -1, 0, fmt.Errorf("rte_mbuf_dyn_rx_timestamp_register failed: %d", int(rc))
	}
	return int(off), uint64(flag), nil
}
