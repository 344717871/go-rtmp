//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package message

type testCase struct {
	Name string
	TypeID
	Value  Message
	Binary []byte
}

var testCases = []testCase{
	testCase{
		Name:   "SetChunkSize",
		TypeID: TypeIDSetChunkSize,
		Value: &SetChunkSize{
			ChunkSize: 1024,
		},
		Binary: []byte{
			// ChunkSize 1024 (*31bit*, BigEndian)
			0x00, 0x00, 0x04, 0x00,
		},
	},
	testCase{
		Name:   "AbortMessage",
		TypeID: TypeIDAbortMessage,
		Value: &AbortMessage{
			ChunkStreamID: 1024,
		},
		Binary: []byte{
			// ChunkStreamID 1024 (32bit, BigEndian)
			0x00, 0x00, 0x04, 0x00,
		},
	},
	testCase{
		Name:   "Ack",
		TypeID: TypeIDAck,
		Value: &Ack{
			SequenceNumber: 1024,
		},
		Binary: []byte{
			// SequenceNumber 1024 (32bit, BigEndian)
			0x00, 0x00, 0x04, 0x00,
		},
	},
	// TODO: UserCtrl
	testCase{
		Name:   "WinAckSize",
		TypeID: TypeIDWinAckSize,
		Value: &WinAckSize{
			Size: 1024,
		},
		Binary: []byte{
			// Size 1024 (32bit, BigEndian)
			0x00, 0x00, 0x04, 0x00,
		},
	},
	testCase{
		Name:   "SetPeerBandwidth",
		TypeID: TypeIDSetPeerBandwidth,
		Value: &SetPeerBandwidth{
			Size:  1024,
			Limit: LimitTypeSoft,
		},
		Binary: []byte{
			// Size 1024 (32bit, BigEndian)
			0x00, 0x00, 0x04, 0x00,
			// Limit Type 1(LimitTypeSoft, 8bit)
			0x01,
		},
	},
	testCase{
		Name:   "AudioMessage",
		TypeID: TypeIDAudioMessage,
		Value: &AudioMessage{
			Payload: []byte("audio data"),
		},
		Binary: []byte("audio data"),
	},
	testCase{
		Name:   "VideoMessage",
		TypeID: TypeIDVideoMessage,
		Value: &VideoMessage{
			Payload: []byte("video data"),
		},
		Binary: []byte("video data"),
	},
	// TODO: DataMessageAMF3
	// TODO: TypeIDSharedObjectMessageAMF3
	// TODO: TypeIDCommandMessageAMF3
	testCase{
		Name:   "DataMessageAMF0",
		TypeID: TypeIDDataMessageAMF0,
		Value: &DataMessageAMF0{
			DataMessage: DataMessage{
				Name: "test",
				Data: nil,
			},
		},
		Binary: []byte{
			// Name: AMF0 / string marker
			0x02,
			// Name: AMF0 / string Length 4
			0x00, 0x04,
			// Name: AMF0 / "test" string
			0x74, 0x65, 0x73, 0x74,
		},
	},
	// TODO: TypeIDSharedObjectMessageAMF0
	testCase{
		Name:   "CommandMessageAMF0",
		TypeID: TypeIDCommandMessageAMF0,
		Value: &CommandMessageAMF0{
			CommandMessage: CommandMessage{
				CommandName:   "_result",
				TransactionID: 10,
				Command:       nil,
			},
		},
		Binary: []byte{
			// CommandName: AMF0 / string marker
			0x02,
			// CommandName: AMF0 / string Length
			0x00, 0x07,
			// CommandName: AMF0 / "_result" string
			0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
			// TransactionID: AMF0 / number marker
			0x00,
			// TransactionID: AMF0 / 10 number
			0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},
	// TODO: TypeIDAggregateMessage
}
