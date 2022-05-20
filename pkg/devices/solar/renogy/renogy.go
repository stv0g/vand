package renogy

import (
	"fmt"
	"strings"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/stv0g/vand/pkg/pb"
)

type Device struct {
	client *modbus.ModbusClient
}

func NewDevice(addr string) (*Device, error) {
	c, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "rtuovertcp://" + addr,
		Speed:   9600,
		Timeout: 2 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	if err = c.Open(); err != nil {
		return nil, err
	}

	return &Device{
		client: c,
	}, nil
}

func (d *Device) GetConfig() (*pb.SolarConfig, error) {
	var off uint16 = 0xa
	regs, err := d.client.ReadRegisters(off, 0x12, modbus.HOLDING_REGISTER)
	if err != nil {
		return nil, err
	}

	var eoff uint16 = 0xe001
	eeprom, err := d.client.ReadRegisters(eoff, 0x14, modbus.HOLDING_REGISTER)
	if err != nil {
		return nil, err
	}

	str := func(r uint16, len uint16) string {
		str := ""
		for i := uint16(0); i < len/2; i++ {
			str += string(byte(regs[r+i-off] >> 8))
			str += string(byte(regs[r+i-off] & 0xff))
		}
		return strings.Trim(str, " ")
	}

	ver := func(r uint16) string {
		major := regs[r-off] & 0xff
		minor := regs[r+1-off] >> 8
		patch := regs[r+1-off] & 0xff

		return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	}

	return &pb.SolarConfig{
		MaxSystemVoltage:        uint32(regs[0xa-off] >> 8),
		RatedChargingCurrent:    uint32(regs[0xa-off] & 0xff),
		RatedDischargingCurrent: uint32(regs[0xb-off] >> 8),
		ProductType:             pb.ProductType(regs[0xa-off] & 0xf),
		BatteryType:             pb.BatteryType(eeprom[0xe004-eoff] & 0xf),
		ProductModel:            str(0xc, 16),
		SoftwareVersion:         ver(0x14),
		HardwareVersion:         ver(0x16),
		ProductSerialNumber:     uint32(regs[0x18+1-off]),
		ProductionYear:          2000 + uint32(regs[0x18-off]>>8),
		ProductionMonth:         uint32(regs[0x18-off] & 0xff),
	}, nil
}

func (d *Device) GetState() (*pb.SolarState, error) {
	var err error

	var off uint16 = 0x100
	regs, err := d.client.ReadRegisters(off, 0x23, modbus.HOLDING_REGISTER)
	if err != nil {
		return nil, err
	}

	temp := func(b uint16) float32 {
		v := float32(b & 0x7f)
		if b&0x80 != 0 {
			v *= -1
		}
		return v
	}

	float := func(r uint16) float32 {
		return float32(regs[r-off])
	}

	word := func(r uint16) uint32 {
		return uint32(regs[r-off])
	}

	dword := func(r uint16) uint32 {
		return uint32(regs[r-off]) << +uint32(regs[r+1-off])
	}

	return &pb.SolarState{
		BatterySoc:               float(0x100),
		BatteryVoltage:           float(0x101) * 0.1,
		ChargingCurrent:          float(0x102) * 0.01,
		TemperatureController:    temp(regs[0x103-off] >> 8),
		TemperatureBattery:       temp(regs[0x103-off] & 0xff),
		PanelVoltage:             float(0x107) * 0.1,
		PanelCurrent:             float(0x108) * 0.01,
		ChargingPower:            float(0x109),
		BatteryVoltageDayMin:     float(0x10b) * 0.1,
		BatteryVoltageDayMax:     float(0x10c) * 0.1,
		ChargingCurrentDayMax:    float(0x10d) * 0.01,
		DischargingCurrentDayMax: float(0x10e) * 0.01,
		ChargingPowerDayMax:      float(0x10f),
		DischargingPowerDayMax:   float(0x110),
		ChargingDayAmpHours:      float(0x111),
		DischargingDayAmpHours:   float(0x112),
		GenerationDay:            float(0x113),
		ConsumptionDay:           float(0x114),
		OperationDays:            word(0x115),
		BatteryOverDischargeCnt:  word(0x116),
		BatteryFullChargeCnt:     word(0x117),
		ChargingTotalAmpHours:    dword(0x118),
		DischargingTotalAmpHours: dword(0x11a),
		GenerationTotal:          dword(0x11c) / 1e4,
		ConsumptionTotal:         dword(0x11e) / 1e4,
		ChargingState:            pb.SolarChargingState(regs[0x120-off] & 0xff),
		Faults:                   dword(0x121),
	}, nil
}
