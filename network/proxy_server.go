package network

import (
	"io"
	"math/rand"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"supernet.tools/tcp-proxy-server/config"
	"supernet.tools/tcp-proxy-server/service"
)

type transferRes chan bool
type applyFunc func([]byte)

// Buffer size for transfer operations
const buffSize = 0xffff

// ProxyServer will provide
type ProxyServer struct {
	listenAddr   string
	destinations []string

	encr service.Encryptor
}

func NewProxy(conf *config.AppConf, encr service.Encryptor) *ProxyServer {
	destinations := conf.Destinations
	if len(destinations) <= 0 {
		destinations = []string{"google.com:80"} // at least some default destination...
	}

	log.Info().Msgf("Registered destinations: %v", destinations)

	return &ProxyServer{
		listenAddr:   conf.ListenAddress(),
		encr:         encr,
		destinations: destinations,
	}
}

func (ps *ProxyServer) Start() {
	listener, err := net.Listen("tcp", ps.listenAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize listener")
	}

	log.Info().Msgf("Listening on: %s", ps.listenAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("Failed to process connection")
			continue
		}

		// Handle connection in separate go routine
		go ps.handle(conn)
	}
}

func (ps *ProxyServer) handle(conn net.Conn) {
	defer conn.Close()

	log.Info().Msg("Processing request")

	// Select destination address and open remote connection
	dest := ps.fetchDest()
	remote, err := net.Dial("tcp", dest)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to establish remote connection to (%s)", dest)
		return
	}
	defer remote.Close()

	destDone := make(transferRes)
	clientDone := make(transferRes)

	go ps.transfer(conn, remote, destDone, ps.encr.Decrypt)
	go ps.transfer(remote, conn, clientDone, ps.encr.Encrypt)

	// Wait until transfer finishes
	<-destDone
	<-clientDone

	log.Info().Msg("Processing finished")
}

func (ps *ProxyServer) transfer(from, to io.ReadWriter, done transferRes, apply applyFunc) {
	buff := make([]byte, buffSize)
	for {
		n, err := from.Read(buff)
		if err != nil {
			ps.handleError(err, done)
			return
		}

		log.Debug().Msgf("Source data read: %d bytes", n)

		data := buff[:n]

		// Apply data manipulation functions
		apply(data)

		_, err = to.Write(data)
		if err != nil {
			ps.handleError(err, done)
			return
		}
	}
}

func (ps *ProxyServer) fetchDest() string {
	// Randomly select one of the registered destinations...
	// Could be applied more interesting logic for load balancing for example...
	// But who cares :)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	idx := r.Intn(len(ps.destinations))

	return ps.destinations[idx]
}

func (ps *ProxyServer) handleError(err error, done transferRes) {
	if err != io.EOF {
		log.Error().Err(err).Msg("Transfer IO operation failed")
	}

	done <- true
}
