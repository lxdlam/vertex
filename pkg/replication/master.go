package replication

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/lxdlam/vertex/pkg/common"
	"github.com/lxdlam/vertex/pkg/log"
	"github.com/lxdlam/vertex/pkg/util"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func init() {
	// Ensure the dir is exist
	_ = os.Mkdir("/tmp/vertex", 0755)
}

type Master interface {
	Start()
	Stop()

	SetFile(*log.PersistentFile, string)

	master()
}

type master struct {
	listener *net.TCPListener
	conn     map[string]net.Conn
	mutex    sync.Mutex
	file     *log.PersistentFile
	filePath string
	started  int32
	shutdown int32
}

func NewMaster(addr *net.TCPAddr) Master {
	if addr == nil {
		return nil
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil
	}

	return &master{
		listener: l,
		conn:     make(map[string]net.Conn),
		file:     nil,
		started:  0,
		shutdown: 0,
	}
}

func (m *master) start() {
	if atomic.CompareAndSwapInt32(&m.started, 0, 1) && atomic.LoadInt32(&m.shutdown) != 1 {
		common.Infof("start master. addr=%s", m.listener.Addr().String())
	Outer:
		for {
			conn, err := m.listener.Accept()
			if err != nil {
				if strings.HasSuffix(err.Error(), "use of closed network connection") {
					break Outer
				} else {
					common.Warnf("tcp listen error. err=%s", err)
					break
				}
			}

			// Not too much connection here, so we just start a new go routine
			go m.handleConn(conn)
		}
	}
}

func (m *master) stop() {
	if atomic.CompareAndSwapInt32(&m.shutdown, 0, 1) {
		m.mutex.Lock()
		m.mutex.Unlock()

		for _, conn := range m.conn {
			_ = conn.Close()
		}
	}
}

func (m *master) handleConn(conn net.Conn) {
	if m.file == nil {
		_ = common.Errorf("no file provided.")
		return
	}

	addr := conn.RemoteAddr().String()
	defer func() { _ = conn.Close() }()

	m.mutex.Lock()
	m.conn[addr] = conn
	m.mutex.Unlock()

	filePath := fmt.Sprintf("/tmp/vertex/%s", util.GenNewUUID())
	tmpFile, err := os.Create(filePath)
	if err != nil {
		_ = common.Errorf("create temporary file failed. addr=%s, path=%s, err=%s", addr, filePath, err.Error())
		return
	}

	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(filePath)
	}()

	err = m.file.Flush()
	if err != nil {
		_ = common.Errorf("save to file failed. addr=%s, path=%s, err=%s", addr, filePath, err.Error())
		return
	}

	rawFile, err := os.Open(m.filePath)
	if err != nil {
		_ = common.Errorf("open raw file failed. addr=%s, path=%s, raw_path=%s, err=%s", addr, filePath, rawFile, err.Error())
		return
	}

	n, err := io.Copy(tmpFile, rawFile)
	if err != nil {
		_ = common.Errorf("copy file failed. addr=%s, path=%s, err=%s", addr, filePath, err.Error())
		return
	}

	ret, err := tmpFile.Seek(0, 0)
	if ret != 0 {
		_ = common.Errorf("seek file failed. addr=%s, path=%s, ret=%d", addr, filePath, ret)
		return
	} else if err != nil {
		_ = common.Errorf("seek file failed. addr=%s, path=%s, ret=%d, err=%s", addr, filePath, ret, err.Error())
		return
	}

	s, err := tmpFile.Stat()
	if err != nil {
		_ = common.Errorf("cannot obtain status of file. addr=%s, path=%s, err=%s", addr, filePath, err.Error())
		return
	}

	// TODO: Only use int32 now. Chunk the file if the size is too large.
	lengthByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthByte, uint32(s.Size()))

	n, err = io.Copy(conn, bytes.NewReader(lengthByte))
	if n != 4 {
		_ = common.Errorf("write length of file mismatch. written=%d, addr=%s, path=%s", n, addr, filePath)
		return
	} else if err != nil {
		_ = common.Errorf("write length of file failed. written=%d, addr=%s, path=%s, err=%s", n, addr, filePath, err.Error())
		return
	}

	n, err = io.Copy(conn, tmpFile)
	if n != s.Size() {
		_ = common.Errorf("write length of file mismatch. written=%d, addr=%s, path=%s", n, addr, filePath)
		return
	} else if err != nil {
		_ = common.Errorf("write length of file failed. written=%d, addr=%s, path=%s, err=%s", n, addr, filePath, err.Error())
		return
	}

	m.mutex.Lock()
	delete(m.conn, addr)
	m.mutex.Unlock()
}

func (m *master) Start() {
	go m.start()
}

func (m *master) Stop() {
	m.stop()
}

func (m *master) SetFile(file *log.PersistentFile, filePath string) {
	m.file = file
	m.filePath = filePath
}

func (m *master) master() {
}
