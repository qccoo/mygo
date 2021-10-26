package main

import (
    "encoding/binary"
    "fmt"
)

const (
    // Constants for sizes
    PACKET_LEN_SIZE  = 4
    HEADER_LEN_SIZE  = 2
    VERSION_SIZE     = 2
    OP_SIZE          = 4
    SEQ_SIZE         = 4
    HEADER_SIZE      = PACKET_LEN_SIZE + HEADER_LEN_SIZE + VERSION_SIZE + OP_SIZE + SEQ_SIZE

    PACKET_LEN_START = 0
    HEADER_LEN_START = PACKET_LEN_START + PACKET_LEN_SIZE
    VERSION_START    = HEADER_LEN_START + HEADER_LEN_SIZE
    OP_START         = VERSION_START + VERSION_SIZE
    SEQ_START        = OP_START + OP_SIZE
    BODY_START       = SEQ_START + SEQ_SIZE
)


type Message struct {
    Version   uint16
    Op        uint32
    Seq       uint32
    Body      []byte
}


// Decodes goim message.
func decode(pack []byte) (*Message, error) {
    if (len(pack) < HEADER_SIZE) {
        return nil, fmt.Errorf(
            "Invalid bytes length %v, should be at least %v.", len(pack), HEADER_SIZE)
    }
    m := &Message{}
    packetLen := binary.BigEndian.Uint32(pack[:HEADER_LEN_START])
    headerLen := binary.BigEndian.Uint16(pack[HEADER_LEN_START:VERSION_START])
    bodyLen := packetLen - uint32(headerLen)
    if (headerLen != HEADER_SIZE) {
        return nil, fmt.Errorf(
            "Invalid header_length value %v, should be %v.", headerLen, HEADER_SIZE)
    }
    m.Version = binary.BigEndian.Uint16(pack[VERSION_START:OP_START])
    m.Op = binary.BigEndian.Uint32(pack[OP_START:SEQ_START])
    m.Seq = binary.BigEndian.Uint32(pack[SEQ_START:BODY_START])
    if bodyLen > 0 {
        if packetLen > uint32(len(pack)) {
            m.Body = pack[BODY_START:]
            return m, fmt.Errorf(
                "Packet bytes length %v is shorter than given packet_length %v.", 
                    len(pack), packetLen)
        }
        m.Body = pack[BODY_START:(BODY_START + bodyLen)]
    }
    return m, nil
}

// Encodes goim message.
func encode(m *Message) []byte {
    packetLen := HEADER_SIZE + len(m.Body)
    pack := make([]byte, packetLen)
    binary.BigEndian.PutUint32(pack[:HEADER_LEN_START], uint32(packetLen))
    binary.BigEndian.PutUint16(pack[HEADER_LEN_START:VERSION_START], uint16(HEADER_SIZE))
    binary.BigEndian.PutUint16(pack[VERSION_START:OP_START], m.Version)
    binary.BigEndian.PutUint32(pack[OP_START:SEQ_START], m.Op)
    binary.BigEndian.PutUint32(pack[SEQ_START:BODY_START], m.Seq)
    copy(pack[BODY_START:], m.Body)
    return pack
}

func main() {
    im := &Message{
        Version: 123,
        Op: 10,
        Seq: 3,
        Body: []byte("Body 123"),
    }
    om, err := decode(encode(im))
    if err != nil {
        fmt.Println("Exited with error:", err)
        return
    }
    fmt.Printf("Decoded message: %v\n", om)
    fmt.Printf("Decoded message body as string: %v\n", string(om.Body))
}
