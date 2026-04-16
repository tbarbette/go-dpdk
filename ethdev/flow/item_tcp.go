package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>

static const struct rte_flow_item_tcp *get_item_tcp_mask() {
	return &rte_flow_item_tcp_mask;
}

static void set_tcp_data_off(struct rte_tcp_hdr *h, uint8_t v) {
	h->data_off = v;
}

static void set_tcp_flags(struct rte_tcp_hdr *h, uint8_t v) {
	h->tcp_flags = v;
}

*/
import "C"
import (
	"runtime"
	"unsafe"
)

// TCPHeader represents TCP header format.
type TCPHeader struct {
	SrcPort  uint16 /* TCP source port. */
	DstPort  uint16 /* TCP destination port. */
	SentSeq  uint32 /* TX data sequence number. */
	RecvAck  uint32 /* RX data acknowledgment sequence number. */
	DataOff  uint8  /* Data offset. */
	TCPFlags uint8  /* TCP flags. */
	RxWin    uint16 /* RX flow control window. */
	Cksum    uint16 /* TCP checksum. */
	TCPUrp   uint16 /* TCP urgent pointer, if any. */
}

// ItemTCP matches a TCP header.
type ItemTCP struct {
	cPointer

	Header TCPHeader
}

var _ ItemStruct = (*ItemTCP)(nil)

// Reload implements ItemStruct interface.
func (item *ItemTCP) Reload() {
	cptr := (*C.struct_rte_flow_item_tcp)(item.createOrRet(C.sizeof_struct_rte_flow_item_tcp))
	cvtTCPHeader(&cptr.hdr, &item.Header)
	runtime.SetFinalizer(item, nil)
	runtime.SetFinalizer(item, (*ItemTCP).free)
}

func cvtTCPHeader(dst *C.struct_rte_tcp_hdr, src *TCPHeader) {
	beU16(src.SrcPort, unsafe.Pointer(&dst.src_port))
	beU16(src.DstPort, unsafe.Pointer(&dst.dst_port))
	beU32(src.SentSeq, unsafe.Pointer(&dst.sent_seq))
	beU32(src.RecvAck, unsafe.Pointer(&dst.recv_ack))
	C.set_tcp_data_off(dst, C.uint8_t(src.DataOff))
	C.set_tcp_flags(dst, C.uint8_t(src.TCPFlags))
	beU16(src.RxWin, unsafe.Pointer(&dst.rx_win))
	beU16(src.Cksum, unsafe.Pointer(&dst.cksum))
	beU16(src.TCPUrp, unsafe.Pointer(&dst.tcp_urp))
}

// Type implements ItemStruct interface.
func (item *ItemTCP) Type() ItemType {
	return ItemTypeTCP
}

// Mask implements ItemStruct interface.
func (item *ItemTCP) Mask() unsafe.Pointer {
	return unsafe.Pointer(C.get_item_tcp_mask())
}
