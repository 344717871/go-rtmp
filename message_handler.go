//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package rtmp

import (
	"github.com/yutopp/go-rtmp/message"
)

type messageHandler interface {
	Handle(chunkStreamID int, timestamp uint32, msg message.Message, stream *Stream) error
	HandleCommand(chunkStreamID int, timestamp uint32, encTy message.EncodingType, cmdMsg *message.CommandMessage, stream *Stream) error
	HandleData(chunkStreamID int, timestamp uint32, encTy message.EncodingType, dataMsg *message.DataMessage, stream *Stream) error
}
