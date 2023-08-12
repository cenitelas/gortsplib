package main

import (
	"log"
	"time"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/formats"
	"github.com/bluenviron/gortsplib/v4/pkg/media"
	"github.com/bluenviron/gortsplib/v4/pkg/url"
	"github.com/pion/rtp"
)

const (
	existingStream = "rtsp://x.x.x.x:8554/mystream"
	reconnectPause = 2 * time.Second
)

type client struct {
	s *server
}

func newClient(s *server) *client {
	c := &client{
		s: s,
	}

	// start a separated routine
	go c.run()

	return c
}

func (c *client) run() {
	for {
		err := c.read()
		log.Printf("ERR: %s\n", err)

		time.Sleep(reconnectPause)
	}
}

func (c *client) read() error {
	rc := gortsplib.Client{}

	// parse URL
	u, err := url.Parse(existingStream)
	if err != nil {
		return err
	}

	// connect to the server
	err = rc.Start(u.Scheme, u.Host)
	if err != nil {
		return err
	}
	defer rc.Close()

	// find published medias
	medias, baseURL, _, err := rc.Describe(u)
	if err != nil {
		return err
	}

	// setup all medias
	err = rc.SetupAll(medias, baseURL)
	if err != nil {
		return err
	}

	stream := c.s.setStreamReady(medias)
	defer c.s.setStreamUnready()

	log.Printf("stream is ready and can be read from the server at rtsp://localhost:8554/stream\n")

	// called when a RTP packet arrives
	rc.OnPacketRTPAny(func(medi *media.Media, forma formats.Format, pkt *rtp.Packet) {
		// route incoming packets to the server stream
		stream.WritePacketRTP(medi, pkt)
	})

	// start playing
	_, err = rc.Play(nil)
	if err != nil {
		return err
	}

	// wait until a fatal error
	return rc.Wait()
}
