package model

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/viorelyo/tlsExperiment/constants"
	"github.com/viorelyo/tlsExperiment/helpers"
	"os"
)

type ClientHello struct {
	RecordHeader             RecordHeader
	HandshakeHeader          HandshakeHeader
	ClientVersion            [2]byte
	ClientRandom             [32]byte
	SessionID                [1]byte
	CipherSuiteLength        [2]byte
	CipherSuite              []byte
	CompressionMethodsLength [1]byte
	CompressionMethods       []byte
}

func MakeClientHello() ClientHello {
	clientHello := ClientHello{}

	recordHeader := RecordHeader{}
	recordHeader.Type = 0x16
	recordHeader.ProtocolVersion = constants.GTlsVersions.GetByteCodeForVersion("TLS 1.0")

	handshakeHeader := HandshakeHeader{}
	handshakeHeader.MessageType = constants.HandshakeClientHello

	clientHello.ClientVersion = constants.GTlsVersions.GetByteCodeForVersion("TLS 1.2")
	// TODO Create random array
	clientHello.ClientRandom = [32]byte{}

	clientHello.SessionID = [1]byte{0x00}

	suitesByteCode := constants.GCipherSuites.GetSuiteByteCodes(constants.GCipherSuites.GetAllSuites())

	clientHello.CipherSuite = suitesByteCode
	clientHello.CipherSuiteLength = helpers.ConvertIntToByteArray(uint16(len(suitesByteCode)))

	clientHello.CompressionMethods = []byte{0x00}
	clientHello.CompressionMethodsLength = [1]byte{0x01}

	handshakeHeader.MessageLength = clientHello.getHandshakeHeaderLength()
	clientHello.HandshakeHeader = handshakeHeader

	recordHeader.Length = clientHello.getRecordLength()
	clientHello.RecordHeader = recordHeader

	return clientHello
}

func (clientHello ClientHello) getHandshakeHeaderLength() [3]byte {
	var length [3]byte
	var k int

	k = len(clientHello.ClientVersion)
	k += len(clientHello.ClientRandom)
	k += len(clientHello.SessionID)
	k += len(clientHello.CipherSuiteLength)
	k += len(clientHello.CipherSuite)
	k += len(clientHello.CompressionMethodsLength)
	k += len(clientHello.CompressionMethods)

	tmp := helpers.ConvertIntToByteArray(uint16(k))
	length[0] = 0x00
	length[1] = tmp[0]
	length[2] = tmp[1]

	return length
}

func (clientHello ClientHello) getRecordLength() [2]byte {
	tmp := int(helpers.Convert3ByteArrayToUInt32(clientHello.HandshakeHeader.MessageLength))
	tmp += 1 // size of MessageType
	tmp += len(clientHello.HandshakeHeader.MessageLength)

	return helpers.ConvertIntToByteArray(uint16(tmp))
}

func (clientHello ClientHello) GetClientHelloPayload() []byte {
	var payload []byte

	payload = append(payload, clientHello.RecordHeader.Type)
	payload = append(payload, clientHello.RecordHeader.ProtocolVersion[:]...)
	payload = append(payload, clientHello.RecordHeader.Length[:]...)
	payload = append(payload, clientHello.HandshakeHeader.MessageType)
	payload = append(payload, clientHello.HandshakeHeader.MessageLength[:]...)
	payload = append(payload, clientHello.ClientVersion[:]...)
	payload = append(payload, clientHello.ClientRandom[:]...)
	payload = append(payload, clientHello.SessionID[:]...)
	payload = append(payload, clientHello.CipherSuiteLength[:]...)
	payload = append(payload, clientHello.CipherSuite[:]...)
	payload = append(payload, clientHello.CompressionMethodsLength[:]...)
	payload = append(payload, clientHello.CompressionMethods...)

	return payload
}

func (clientHello ClientHello) SaveJSON() {
	file, _ := os.OpenFile("ClientHello.json", os.O_CREATE, os.ModePerm)
	defer file.Close()
	_ = json.NewEncoder(file).Encode(&clientHello)
}

func (clientHello ClientHello) String() string {
	out := fmt.Sprintf("Client Hello\n")
	out += fmt.Sprint(clientHello.RecordHeader)
	out += fmt.Sprint(clientHello.HandshakeHeader)
	out += fmt.Sprintf("  Client Version.....: %6x - %s\n", clientHello.ClientVersion, constants.GTlsVersions.GetVersionForByteCode(clientHello.ClientVersion))
	out += fmt.Sprintf("  Client Random......: %6x\n", clientHello.ClientRandom)
	out += fmt.Sprintf("  Session ID.........: %6x\n", clientHello.SessionID)
	out += fmt.Sprintf("  CipherSuite Len....: %6x\n", clientHello.CipherSuiteLength)
	out += fmt.Sprintf("  CipherSuites.......:\n")
	for _, c := range helpers.ConvertByteArrayToCipherSuites(clientHello.CipherSuite) {
		out += fmt.Sprintf("       %s\n", c)
	}
	out += fmt.Sprintf("  CompressionMethods Len..: %6x\n", clientHello.CompressionMethodsLength)
	out += fmt.Sprintf("  CompressionMethods..: %6x\n", clientHello.CompressionMethods)
	return out
}

func (clientHello *ClientHello) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		RecordHeader             RecordHeader    `json:"RecordHeader"`
		HandshakeHeader          HandshakeHeader `json:"HandshakeHeader"`
		ClientVersion            string          `json:"ClientVersion"`
		ClientRandom             string          `json:"ClientRandom"`
		SessionID                uint8           `json:"SessionID"`
		CipherSuiteLength        uint16          `json:"CipherSuiteLength"`
		CipherSuites             []string        `json:"CipherSuites"`
		CompressionMethodsLength uint8           `json:"CompressionMethodsLength"`
		CompressionMethods       string          `json:"CompressionMethods"`
	}{
		RecordHeader:             clientHello.RecordHeader,
		HandshakeHeader:          clientHello.HandshakeHeader,
		ClientVersion:            constants.GTlsVersions.GetVersionForByteCode(clientHello.ClientVersion),
		ClientRandom:             hex.EncodeToString(clientHello.ClientRandom[:]),
		SessionID:                clientHello.SessionID[0],
		CipherSuiteLength:        helpers.ConvertByteArrayToUInt16(clientHello.CipherSuiteLength),
		CipherSuites:             helpers.ConvertByteArrayToCipherSuites(clientHello.CipherSuite),
		CompressionMethodsLength: clientHello.CompressionMethodsLength[0],
		CompressionMethods:       hex.EncodeToString(clientHello.CompressionMethods),
	})
}
