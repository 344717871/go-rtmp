package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yutopp/go-flv"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
	"log"
	"os"
	"path/filepath"
)

var _ rtmp.Handler = (*Handler)(nil)

type Handler struct {
	rtmp.DefaultHandler
	flvFile *os.File
	flvEnc  *flv.Encoder
}

func (h *Handler) OnServe() {
}

func (h *Handler) OnConnect(timestamp uint32, cmd *rtmpmsg.NetConnectionConnect) error {
	log.Printf("OnConnect: %#v", cmd)
	return nil
}

func (h *Handler) OnCreateStream(timestamp uint32, cmd *rtmpmsg.NetConnectionCreateStream) error {
	log.Printf("OnCreateStream: %#v", cmd)
	return nil
}

func (h *Handler) OnPublish(timestamp uint32, cmd *rtmpmsg.NetStreamPublish) error {
	log.Printf("OnPublish: %#v", cmd)

	// record streams as FLV!
	p := filepath.Join(
		os.TempDir(),
		filepath.Clean(filepath.Join("/", fmt.Sprintf("%s.flv", cmd.PublishingName))),
	)
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return errors.Wrap(err, "Failed to create flv file")
	}
	h.flvFile = f

	enc, err := flv.NewEncoder(f, flv.FlagsAudio|flv.FlagsVideo)
	if err != nil {
		_ = f.Close()
		return errors.Wrap(err, "Failed to create flv encoder")
	}
	h.flvEnc = enc

	return nil
}

func (h *Handler) OnSetDataFrame(timestamp uint32, data *rtmpmsg.NetStreamSetDataFrame) error {
	r := bytes.NewReader(data.Payload)

	var script flvtag.ScriptData
	if err := flvtag.DecodeScriptData(r, &script); err != nil {
		log.Printf("Failed to decode script data: Err = %+v", err)
		return nil // ignore
	}

	log.Printf("SetDataFrame: Script = %#v", script)

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeScriptData,
		Timestamp: timestamp,
		Data:      &script,
	}); err != nil {
		log.Printf("Failed to write script data: Err = %+v", err)
	}

	return nil
}

func (h *Handler) OnAudio(timestamp uint32, payload []byte) error {
	r := bytes.NewReader(payload)

	var audio flvtag.AudioData
	if err := flvtag.DecodeAudioData(r, &audio); err != nil {
		return err
	}

	log.Printf("FLV Audio Data: Timestamp = %d, SoundFormat = %+v, SoundRate = %+v, SoundSize = %+v, SoundType = %+v, AACPacketType = %+v, Data length = %+v",
		timestamp,
		audio.SoundFormat,
		audio.SoundRate,
		audio.SoundSize,
		audio.SoundType,
		audio.AACPacketType,
		len(audio.Data),
	)

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeAudio,
		Timestamp: timestamp,
		Data:      &audio,
	}); err != nil {
		log.Printf("Failed to write audio: Err = %+v", err)
	}

	return nil
}

func (h *Handler) OnVideo(timestamp uint32, payload []byte) error {
	r := bytes.NewReader(payload)

	var video flvtag.VideoData
	if err := flvtag.DecodeVideoData(r, &video); err != nil {
		return err
	}

	log.Printf("FLV Video Data: Timestamp = %d, FrameType = %+v, CodecID = %+v, AVCPacketType = %+v, CT = %+v, Data length = %+v",
		timestamp,
		video.FrameType,
		video.CodecID,
		video.AVCPacketType,
		video.CompositionTime,
		len(video.Data),
	)

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeVideo,
		Timestamp: timestamp,
		Data:      &video,
	}); err != nil {
		log.Printf("Failed to write video: Err = %+v", err)
	}

	return nil
}

func (h *Handler) OnClose() {
	log.Printf("OnClose")

	if h.flvFile != nil {
		_ = h.flvFile.Close()
	}
}
