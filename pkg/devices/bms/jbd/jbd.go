package jbd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
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

	// EEPROM  registers
	regCycleCount   = 0x17
	regSerialNumber = 0xA0
	regModel        = 0xA1
	regBarCode      = 0xA2
	regFETCtrl      = 0xE1
)

type Device struct {
	conn    Conn
	timeout time.Duration

	config *pb.BatteryConfig
}

func NewDevice(addr string) (*Device, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	d := &Device{
		conn:    Conn{c},
		timeout: 1 * time.Second,
	}

	d.config, err = d.ReadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	log.Printf("Config: %+#v", d.config)

	return d, nil
}

func (d *Device) send(cmd byte, payload []byte) error {
	plen := len(payload)
	buf := make([]byte, plen+5)
	buf[0] = byteStart
	buf[1] = cmd
	buf[plen+4] = byteEnd
	copy(buf[2:], payload)
	binary.BigEndian.PutUint16(buf[plen+2:], checksum(payload))

	if err := d.conn.SetWriteDeadline(time.Now().Add(d.timeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	if _, err := d.conn.Write(buf); err != nil {
		return fmt.Errorf("failed to send cmd: %w", err)
	}

	return nil
}

func (d *Device) recv(cmd byte) ([]byte, error) {
	buf := make([]byte, 128)

	if err := d.conn.SetReadDeadline(time.Now().Add(d.timeout)); err != nil {
		return nil, err
	}

	m, err := d.conn.Read(buf)

	plen := int(buf[3])

	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	} else if buf[0] != byteStart {
		return nil, errors.New("invalid start byte")
	} else if buf[1] != cmd {
		return nil, errors.New("not a response")
	} else if buf[2] != 0x00 {
		return nil, errors.New("non-zero status")
	} else if m != plen+7 {
		return nil, fmt.Errorf("received invalid number of bytes: %d != %d", m, plen+7)
	} else if buf[plen+6] != byteEnd {
		return nil, errors.New("invalid end byte")
	}

	chks := binary.BigEndian.Uint16(buf[plen+4:])
	payload := buf[4 : 4+plen]

	if checksum(buf[2:4+plen]) != chks {
		return nil, errors.New("invalid checksum")
	}

	return payload, nil
}

func (d *Device) ReadRegister(off int) ([]byte, error) {
	payload := []byte{byte(off), 0}

	if err := d.send(cmdRead, payload); err != nil {
		return nil, err
	}

	response, err := d.recv(byte(off))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *Device) WriteRegister(off int, data []byte) error {
	payload := []byte{byte(off), byte(len(data))}
	payload = append(payload, data...)

	if err := d.send(cmdWrite, payload); err != nil {
		return err
	}

	_, err := d.recv(byte(off))
	if err != nil {
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

func (d *Device) GetState() (*pb.BatteryState, error) {
	b, err := d.ReadRegister(regBasicInfo)
	if err != nil {
		return nil, err
	}

	buf := Buf(b)

	state := &pb.BatteryState{
		Soc:            uint32(buf[0x13]),
		PackVoltage:    float32(buf.Uint16(0x00)) * 10e-3,
		PackCurrent:    float32(buf.Int16(0x02)) * 10e-3,
		CycleCapacity:  float32(buf.Uint16(0x04)) * 10e-3,
		DesignCapacity: float32(buf.Uint16(0x06)) * 10e-3,

		Balancing:          buf.Uint32(0x0C),
		Errors:             uint32(buf.Uint16(0x10)),
		CycleCount:         uint32(buf.Uint16(0x08)),
		ChargeFetEnable:    buf.Bit(0x14, 0),
		DischargeFetEnable: buf.Bit(0x14, 1),
		Temperatures:       []float32{},
		Voltages:           []float32{},
	}

	d.config.NtcCount = 2

	for i := 0; i < int(d.config.NtcCount); i++ {
		state.Temperatures = append(state.Temperatures, float32(buf.Uint16(0x17+i*2))*0.1-273.1)
	}

	b, err = d.ReadRegister(regCellVoltages)
	if err != nil {
		return nil, err
	}

	buf = Buf(b)

	for i := 0; i < int(d.config.CellCount); i++ {
		state.Voltages = append(state.Voltages, float32(buf.Uint16(i<<1))*1e-3)
	}

	return state, nil
}

func (d *Device) ReadDeviceName() (string, error) {
	b, err := d.ReadRegister(regDeviceName)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (d *Device) ReadConfig() (*pb.BatteryConfig, error) {
	b, err := d.ReadRegister(regBasicInfo)
	if err != nil {
		return nil, err
	}

	buf := Buf(b)

	mfcDate := uint32(buf.Uint16(0x0A))

	snBuf, err := d.ReadEEPROM(regSerialNumber)
	if err != nil {
		return nil, err
	}

	barCodeBuf, err := d.ReadEEPROM(regBarCode)
	if err != nil {
		return nil, err
	}

	modelBuf, err := d.ReadEEPROM(regModel)
	if err != nil {
		return nil, err
	}

	devName, err := d.ReadDeviceName()
	if err != nil {
		return nil, err
	}

	config := &pb.BatteryConfig{
		ProductionYear:  (mfcDate >> 9) & 0x7F,
		ProductionMonth: (mfcDate >> 5) & 0x0F,
		ProductionDay:   (mfcDate >> 0) & 0x1F,
		SoftwareVersion: uint32(buf[0x12]),
		CellCount:       uint32(buf[0x15]),
		NtcCount:        uint32(buf[0x16]),
		DeviceName:      devName,
		SerialNumber:    string(snBuf[1:]),
		Model:           string(modelBuf[1:]),
		Barcode:         string(barCodeBuf[1:]),
	}

	return config, nil
}

func (d *Device) SetFET(charge, discharge bool) error {
	if err := d.EnterFactory(); err != nil {
		return err
	}
	defer d.ExitFactory(false)

	var val uint16 = 0
	if !charge {
		val |= 1 << 0
	}

	if !discharge {
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
	return d.WriteRegister(0x00, []byte{0x56, 0x78})
}

func (d *Device) ExitFactory(save bool) error {
	if save {
		return d.WriteRegister(0x01, []byte{0x28, 0x28})
	} else {
		return d.WriteRegister(0x01, []byte{0x00, 0x00})
	}
}

func (d *Device) ReadEEPROM(off int) ([]byte, error) {
	if err := d.EnterFactory(); err != nil {
		return nil, err
	}
	defer d.ExitFactory(false)

	return d.ReadRegister(off)
}
