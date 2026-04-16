package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>

static const struct rte_flow_item_ipv6 *get_item_ipv6_mask() {
	return &rte_flow_item_ipv6_mask;
}

*/
import "C"
import (
	"runtime"
	"unsafe"
)

// IPv6 represents a raw IPv6 address (16 bytes).
type IPv6 [16]byte

// IPv6Header is the IPv6 header raw format.
type IPv6Header struct {
	VtcFlow    uint32 /* IP version, traffic class & flow label. */
	PayloadLen uint16 /* IP payload size, including ext. headers. */
	Proto      uint8  /* Protocol, next header. */
	HopLimits  uint8  /* Hop limits. */
	SrcAddr    IPv6   /* Source address. */
	DstAddr    IPv6   /* Destination address. */
}

// ItemIPv6 matches an IPv6 header.
type ItemIPv6 struct {
	cPointer

	Header IPv6Header
}

var _ ItemStruct = (*ItemIPv6)(nil)

// Reload implements ItemStruct interface.
func (item *ItemIPv6) Reload() {
	cptr := (*C.struct_rte_flow_item_ipv6)(item.createOrRet(C.sizeof_struct_rte_flow_item_ipv6))
	cvtIPv6Header(&cptr.hdr, &item.Header)
	runtime.SetFinalizer(item, nil)
	runtime.SetFinalizer(item, (*ItemIPv6).free)
}

func cvtIPv6Header(dst *C.struct_rte_ipv6_hdr, src *IPv6Header) {
	beU32(src.VtcFlow, unsafe.Pointer(&dst.vtc_flow))
	beU16(src.PayloadLen, unsafe.Pointer(&dst.payload_len))
	dst.proto = C.uint8_t(src.Proto)
	dst.hop_limits = C.uint8_t(src.HopLimits)

	for i := 0; i < 16; i++ {
		dst.src_addr[i] = C.uint8_t(src.SrcAddr[i])
		dst.dst_addr[i] = C.uint8_t(src.DstAddr[i])
	}
}

// Type implements ItemStruct interface.
func (item *ItemIPv6) Type() ItemType {
	return ItemTypeIPv6
}

// Mask implements ItemStruct interface.
func (item *ItemIPv6) Mask() unsafe.Pointer {
	return unsafe.Pointer(C.get_item_ipv6_mask())
}
