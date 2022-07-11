package jbd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/stv0g/vand/pkg/pb"
)

const (
	byteStart byte = 0xDD
	byteEnd   byte = 0x77
	cmdRead   byte = 0xA5
	cmdWrite  byte = 0x5A

	regBasicInfo    = 3
	regCellVoltages = 4
	regDeviceName   = 5

	regFETCtrl = 0xE1
)

type Device struct {
	conn    net.Conn
	timeout time.Duration

	config *pb.BatteryConfig
}

func NewDevice(addr string) (*Device, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	d := &Device{
		conn:    c,
		timeout: time.Second,
	}

	return d, nil
}

func (d *Device) send(payload []byte) error {
	buf := make([]byte, len(payload)+4)
	buf = append(buf, byteStart)
	buf = append(buf, payload...)
	binary.BigEndian.PutUint16(buf, checksum(payload))
	buf = append(buf, byteEnd)

	if err := d.conn.SetWriteDeadline(time.Now().Add(d.timeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	if _, err := d.conn.Write(buf); err != nil {
		return fmt.Errorf("failed to send cmd: %w", err)
	}

	return nil
}

func (d *Device) recv(n int) ([]byte, error) {
	buf := make([]byte, n+4)

	if err := d.conn.SetReadDeadline(time.Now().Add(d.timeout)); err != nil {

	}

	m, err := d.conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	if m != n+4 {
		return nil, fmt.Errorf("received invalid number of bytes: %d != %d", m, n+4)
	}

	if buf[0] != byteStart {
		return nil, errors.New("invalid start byte")
	}

	if buf[n+3] != byteEnd {
		return nil, errors.New("invalid end byte")
	}

	chks := binary.BigEndian.Uint16(buf[n+1:])
	payload := buf[1 : 1+n]

	if checksum(payload) != chks {
		return nil, errors.New("invalid checksum")
	}

	return payload, nil
}

func (d *Device) ReadRegister(off, n int) ([]byte, error) {
	payload := []byte{cmdRead, byte(off), byte(n)}

	if err := d.send(payload); err != nil {
		return nil, err
	}

	response, err := d.recv(n)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *Device) WriteRegister(off int, data []byte) error {
	payload := []byte{cmdWrite, byte(off), byte(len(data))}
	payload = append(payload, data...)

	if err := d.send(payload); err != nil {
		return err
	}

	return nil
}

func checksum(data []byte) uint16 {
	sum := 0x10000

	for _, b := range data {
		sum -= int(b)
	}

	return uint16(sum)
}

func (d *Device) ReadState() (*pb.BatteryState, error) {
	ntc_count := 2

	b, err := d.ReadRegister(regBasicInfo, 0x16+ntc_count*2)
	if err != nil {
		return nil, err
	}

	buf := Buf(b)

	state := &pb.BatteryState{
		PackVoltage:    float32(buf.Uint16(0x00)) * 10e-3,
		PackCurrent:    float32(buf.Int16(0x02)) * 10e-3,
		CycleCapacity:  float32(buf.Uint16(0x04)) * 10e-3,
		DesignCapacity: float32(buf.Uint16(0x06)) * 10e-3,

		Balancing:          buf.Uint32(0x03),
		Errors:             uint32(buf.Uint16(0x10)),
		CycleCount:         uint32(buf.Uint16(0x08)),
		ChargeFetEnable:    buf.Bit(0x13, 0),
		DischargeFetEnable: buf.Bit(0x13, 1),
	}

	return state, nil
}

func (d *Device) ReadDeviceName() (string, error) {
	b, err := d.ReadRegister(regDeviceName, 1)
	if err != nil {
		return "", err
	}

	length := int(b[0])

	b, err = d.ReadRegister(regDeviceName, length)
	if err != nil {
		return "", err
	}

	return string(b[1:]), nil
}

func (d *Device) ReadConfig() (*pb.BatteryConfig, error) {
	b, err := d.ReadRegister(regBasicInfo, 0x16)
	if err != nil {
		return nil, err
	}

	devName, err := d.ReadDeviceName()
	if err != nil {
		return nil, err
	}

	buf := Buf(b)

	mfcDate := uint32(buf.Uint16(0xA))

	state := &pb.BatteryConfig{
		ProductionYear:  (mfcDate >> 9) & 0x7F,
		ProductionMonth: (mfcDate >> 5) & 0xF,
		ProductionDay:   (mfcDate >> 0) & 0x1F,
		SoftwareVersion: uint32(buf[0x11]),
		CellCount:       uint32(buf[0x14]),
		NtcCount:        uint32(buf[0x15]),
		DeviceName:      devName,
	}

	return state, nil
}

func (d *Device) SetFET(charge, discharge bool) error {
	if err := d.EnterFactory(); err != nil {
		return err
	}
	defer d.ExitFactory()

	var val uint16 = 0
	if charge {
		val |= 1 << 0
	}

	if discharge {
		val |= 1 << 1
	}

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)

	if err := d.WriteRegister(regFETCtrl, buf); err != nil {
		return err
	}

	return nil
}

func (d *Device) EnterFactory() error {
	// TODO

	return nil
}

func (d *Device) ExitFactory() error {
	// TODO

	return nil
}
