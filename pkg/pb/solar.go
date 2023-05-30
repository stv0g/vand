// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package pb

const (
	CircuitChargeMOSShortCircuit      = (1 << 30)
	AntiReverseMOSShort               = (1 << 29)
	SolarPanelReverselyConnected      = (1 << 28)
	SolarPanelWorkingPointOverVoltage = (1 << 27)
	SolarPanelCounterCurrent          = (1 << 26)
	PhotovoltaicInputSideOverVoltage  = (1 << 25)
	PhotovoltaicInputSideShortCircuit = (1 << 24)
	PhotovoltaicInputOverpower        = (1 << 23)
	AmbientTemperatureTooHigh         = (1 << 22)
	ControllerTemperatureTooHigh      = (1 << 21)
	LoadOverPowerOrLoadOverCurrent    = (1 << 20)
	LoadShortCircuit                  = (1 << 19)
	BatteryUnderVoltageWarning        = (1 << 18)
	BatteryOverVoltage                = (1 << 17)
	BatteryOverDischarge              = (1 << 16)
)
