package modem

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/stv0g/vand/pkg/pb"
)

type AutoConnectMode string
type FailOverMode string

const (
	AutoConnectNever       AutoConnectMode = "Never"
	AutoConnectHomeNetwork AutoConnectMode = "HomeNetwork"
	AutoConnectAlways      AutoConnectMode = "Always"

	FailOverAuto FailOverMode = "Auto"
	FailOverWAN  FailOverMode = "WAN"
	FailOverLTE  FailOverMode = "LTE"
)

type Model struct {
	Custom        Custom        `json:"custom"`
	WebDaemon     WebDaemon     `json:"webd"`
	LCD           LCD           `json:"lcd"`
	LED           LED           `json:"led"`
	SIM           SIM           `json:"sim"`
	SMS           SMS           `json:"sms"`
	Session       Session       `json:"session"`
	General       General       `json:"general"`
	Power         Power         `json:"power"`
	WWAN          WWAN          `json:"wwan"`
	WWANAdvanced  WWANAdvanced  `json:"wwanadv"`
	Ethernet      Ethernet      `json:"ethernet"`
	Wifi          Wifi          `json:"wifi"`
	Router        Router        `json:"router"`
	FirmwareOTA   FirmwareOTA   `json:"fota"`
	Failover      FailOver      `json:"failover"`
	Cradle        Cradle        `json:"cradle"`
	EventLog      EventLog      `json:"eventlog"`
	UI            UI            `json:"ui"`
	AccessControl AccessControl `json:"accesscontrol"`
	Ready         Ready         `json:"ready"`
	Insight       Insight       `json:"insight"`
}

type FirmwareOTA struct {
	FWupdater struct {
		Available     bool   `json:"available"`
		Chkallow      bool   `json:"chkallow"`
		Chkstatus     string `json:"chkstatus"`
		DloadProg     int    `json:"dloadProg"`
		Error         bool   `json:"error"`
		LastChkDate   int    `json:"lastChkDate"`
		State         string `json:"state"`
		IsPostponable bool   `json:"isPostponable"`
		StatusCode    int    `json:"statusCode"`
		ChkTimeLeft   int    `json:"chkTimeLeft"`
		DloadSize     int    `json:"dloadSize"`
		IsRejectable  bool   `json:"isRejectable"`
		Description   string `json:"description"`
	} `json:"fwupdater"`
}

type FailOver struct {
	Enabled         bool         `json:"enabled"`
	Backhaul        string       `json:"backhaul"`
	Supported       bool         `json:"supported"`
	MonitorPeriod   int          `json:"monitorPeriod"`
	WANConnected    bool         `json:"wanConnected"`
	KeepaliveEnable bool         `json:"keepaliveEnable"`
	KeepaliveSleep  int          `json:"keepaliveSleep"`
	IPv4Targets     []IPv4Target `json:"ipv4Targets"`
	IPv6Targets     []IPv6Target `json:"ipv6Targets"`
}

type IPv4Target struct {
	ID     string `json:"id,omitempty"`
	String string `json:"string,omitempty"`
}

type IPv6Target struct {
}

type UI struct {
	ServerDaysLeftHide   bool `json:"serverDaysLeftHide"`
	PromptActivation     bool `json:"promptActivation"`
	StealthEnabled       bool `json:"stealthEnabled"`
	DisableSettingsOnLCD bool `json:"disableSettingsOnLCD"`
}

type Power struct {
	PMState          string  `json:"PMState"`
	SmState          string  `json:"SmState"`
	BattLowThreshold int     `json:"battLowThreshold"`
	AutoOff          AutoOff `json:"autoOff"`
	Standby          struct {
		Ethernet struct {
			Timer struct {
				OnBattery int `json:"onBattery"`
			} `json:"timer"`
		} `json:"ethernet"`
		OnIdle struct {
			Timer struct {
				OnAC      int `json:"onAC"`
				OnBattery int `json:"onBattery"`
				OnUSB     int `json:"onUSB"`
			} `json:"timer"`
		} `json:"onIdle"`
	} `json:"standby"`
	AutoOn struct {
		Enable bool `json:"enable"`
	} `json:"autoOn"`
	BatteryTemperature  int    `json:"batteryTemperature"`
	BatteryVoltage      int    `json:"batteryVoltage"`
	BattChargeLevel     int    `json:"battChargeLevel"`
	BattChargeSource    string `json:"battChargeSource"`
	BatteryState        string `json:"batteryState"`
	BattChargeAlgorithm string `json:"battChargeAlgorithm"`
	Charging            bool   `json:"charging"`
	ButtonHoldTime      int    `json:"buttonHoldTime"`
	DeviceTempCritical  bool   `json:"deviceTempCritical"`
	Resetreason         int    `json:"resetreason"`
	ResetRequired       string `json:"resetRequired"`
	ChoiceOnUsb         string `json:"choiceOnUsb"`
	ActionOnUsb         string `json:"actionOnUsb"`
	BatteryTempState    string `json:"batteryTempState"`
	BatteryName         string `json:"batteryName"`
	WifiOff             struct {
		OnUsbConnect bool `json:"onUsbConnect"`
		OnTethered   bool `json:"onTethered"`
	} `json:"wifiOff"`
	Boost struct {
		CableConnected bool `json:"cableConnected"`
	} `json:"boost"`
	Lpm bool `json:"lpm"`
}

type AutoOff struct {
	OnUSBdisconnect struct {
		Enable         bool `json:"enable"`
		CountdownTimer int  `json:"countdownTimer"`
	} `json:"onUSBdisconnect"`
	OnIdle struct {
		Timer struct {
			OnAC      int `json:"onAC"`
			OnBattery int `json:"onBattery"`
		} `json:"timer"`
	} `json:"onIdle"`
}

type EventLog struct {
	Level       int `json:"level"`
	LogDuration int `json:"logDuration"`
	CatLevel    []struct {
		Category     string `json:"category,omitempty"`
		CatlevelMask int    `json:"catlevel_mask,omitempty"`
		Enabled      bool   `json:"enabled,omitempty"`
	} `json:"catLevel"`
}

type Cradle struct {
	Mode                bool   `json:"mode"`
	SmartMode           bool   `json:"smartMode"`
	URL                 string `json:"url"`
	PrimarySSID         string `json:"primarySSID"`
	PrimaryPassphrase   string `json:"primaryPassphrase"`
	SecondarySSID       string `json:"secondarySSID"`
	SecondaryPassphrase string `json:"secondaryPassphrase"`
	AutoIPT             bool   `json:"autoIPT"`
}

type Router struct {
	GatewayIP           string `json:"gatewayIP"`
	DMZaddress          string `json:"DMZaddress"`
	DMZenabled          bool   `json:"DMZenabled"`
	ForceSetup          bool   `json:"forceSetup"`
	DHCP                DHCP   `json:"DHCP"`
	UsbMode             string `json:"usbMode"`
	UsbNetworkTethering bool   `json:"usbNetworkTethering"`
	PortFwdEnabled      bool   `json:"portFwdEnabled"`
	UsbTetheringActive  bool   `json:"usbTetheringActive"`
	PortFwdList         []struct {
	} `json:"portFwdList"`
	PortFwdAllowEntry    int    `json:"portFwdAllowEntry"`
	PortFilteringEnabled bool   `json:"portFilteringEnabled"`
	PortFilteringMode    string `json:"portFilteringMode"`
	PortFilterWhiteList  []struct {
	} `json:"portFilterWhiteList"`
	PortFilterBlackList []struct {
	} `json:"portFilterBlackList"`
	ClientList struct {
		List []struct {
			IP     string `json:"IP,omitempty"`
			MAC    string `json:"MAC,omitempty"`
			Name   string `json:"name,omitempty"`
			Media  string `json:"media,omitempty"`
			Source string `json:"source,omitempty"`
		} `json:"list"`
		Count int `json:"count"`
	} `json:"clientList"`
	HostName               string `json:"hostName"`
	DomainName             string `json:"domainName"`
	IPPassThroughEnabled   bool   `json:"ipPassThroughEnabled"`
	IPPassThroughSupported bool   `json:"ipPassThroughSupported"`
	VPNpassthrough         bool   `json:"VPNpassthrough"`
	UsbTetherSupported     bool   `json:"usbTetherSupported"`
	NetMask                string `json:"netMask"`
	BridgeLanNetMask       string `json:"bridgeLanNetMask"`
	Ipv6Supported          bool   `json:"Ipv6Supported"`
	UPNPsupported          bool   `json:"UPNPsupported"`
	UPNPenabled            bool   `json:"UPNPenabled"`
	UPNPIgdVersion         int    `json:"UPNPIgdVersion"`
}

type DHCP struct {
	ServerEnabled bool   `json:"serverEnabled"`
	DNS1          string `json:"DNS1"`
	DNS2          string `json:"DNS2"`
	DNSmode       string `json:"DNSmode"`
	USBpcIP       string `json:"USBpcIP"`
	Range         struct {
		High string `json:"high"`
		Low  string `json:"low"`
	} `json:"range"`
}

type Wifi struct {
	Supported          string `json:"supported"`
	Enabled            bool   `json:"enabled"`
	Status             string `json:"status"`
	Mode               string `json:"mode"`
	MaxClientSupported int    `json:"maxClientSupported"`
	MaxClientLimit     int    `json:"maxClientLimit"`
	MaxClientCnt       int    `json:"maxClientCnt"`
	Encryption         string `json:"encryption"`
	Channel            int    `json:"channel"`
	TwoGBandwidth      string `json:"2gBandwidth"`
	FiveGBandwidth     string `json:"5gBandwidth"`
	TxPower            string `json:"txPower"`
	HiddenSSID         bool   `json:"hiddenSSID"`
	PassPhrase         string `json:"passPhrase"`
	PasswordReminder   bool   `json:"passwordReminder"`
	RTSthreshold       int    `json:"RTSthreshold"`
	FragThreshold      int    `json:"fragThreshold"`
	AccessControl      string `json:"accessControl"`
	SSID               string `json:"SSID"`
	WmmEnabled         bool   `json:"wmmEnabled"`
	MAC                string `json:"MAC"`
	SSIDreminder       bool   `json:"SSIDreminder"`
	ClientCount        int    `json:"clientCount"`
	Country            string `json:"country"`
	Privacy            bool   `json:"privacy"`
	TwoGMcsChanList    string `json:"2gMcsChanList"`
	FiveGMcsChanList   string `json:"5gMcsChanList"`
	WPS                struct {
		Supported string `json:"supported"`
		Enabled   string `json:"enabled"`
		Blocked   bool   `json:"blocked"`
		Mode      string `json:"mode"`
		Status    string `json:"status"`
	} `json:"wps"`
	Guest struct {
		MaxClientCnt       int    `json:"maxClientCnt"`
		Enabled            bool   `json:"enabled"`
		Status             string `json:"status"`
		Encryption         string `json:"encryption"`
		SSID               string `json:"SSID"`
		PassPhrase         string `json:"passPhrase"`
		GeneratePassphrase bool   `json:"generatePassphrase"`
		AccessProfile      string `json:"accessProfile"`
		HiddenSSID         bool   `json:"hiddenSSID"`
		TimerEnable        bool   `json:"timerEnable"`
		TimerTimestamp     int    `json:"timerTimestamp"`
		TimerValue         int    `json:"timerValue"`
		Privacy            bool   `json:"privacy"`
		Chan               int    `json:"chan"`
		Mode               string `json:"mode"`
		DHCP               struct {
			Range struct {
				High string `json:"high"`
				Low  string `json:"low"`
			} `json:"range"`
		} `json:"DHCP"`
	} `json:"guest"`
	Aux struct {
		Enabled      bool   `json:"enabled"`
		AuxMode      string `json:"AuxMode"`
		SSID         string `json:"SSID"`
		HiddenSSID   bool   `json:"hiddenSSID"`
		Encryption   string `json:"encryption"`
		PassPhrase   string `json:"passPhrase"`
		MaxClientCnt int    `json:"maxClientCnt"`
		Status       string `json:"status"`
		Chan         int    `json:"chan"`
		Mode         string `json:"mode"`
	} `json:"aux"`
	Offload struct {
		ActivationRequired bool   `json:"activationRequired"`
		Bars               int    `json:"bars"`
		Enabled            bool   `json:"enabled"`
		Rssi               int    `json:"rssi"`
		SecurityStatus     string `json:"securityStatus"`
		Status             string `json:"status"`
		Supported          bool   `json:"supported"`
		ConnectionSsid     string `json:"connectionSsid"`
		ScanProgress       string `json:"scanProgress"`
		StationIPv4        string `json:"stationIPv4"`
		TimeOn             int    `json:"timeOn"`
		DataTransferred    struct {
			Rx int `json:"rx"`
			Tx int `json:"tx"`
		} `json:"dataTransferred"`
		NetworkList []struct {
		} `json:"networkList"`
		ScanList []struct {
		} `json:"scanList"`
	} `json:"offload"`
	AccessBlackList struct {
		List []struct {
		} `json:"list"`
		Count int `json:"count"`
	} `json:"accessBlackList"`
	AccessWhiteList struct {
		List []struct {
		} `json:"list"`
		Count int `json:"count"`
	} `json:"accessWhiteList"`
}

type Ethernet struct {
	Mac     string `json:"mac"`
	Offload struct {
		Supported bool   `json:"supported"`
		Enabled   bool   `json:"enabled"`
		On        bool   `json:"on"`
		Ipv4Addr  string `json:"ipv4Addr"`
		Ipv6Addr  string `json:"ipv6Addr"`
		Rx        int    `json:"rx"`
		Tx        int    `json:"tx"`
		TimeOn    int    `json:"timeOn"`
	} `json:"offload"`
}

type WWANAdvanced struct {
	CurBand           string `json:"curBand"`
	RadioQuality      int    `json:"radioQuality"`
	Country           string `json:"country"`
	RAC               int    `json:"RAC"`
	LAC               int    `json:"LAC"`
	MCC               string `json:"MCC"`
	MNC               string `json:"MNC"`
	MNCFmt            int    `json:"MNCFmt"`
	CellID            int    `json:"cellId"`
	ChanID            int    `json:"chanId"`
	PrimScode         int    `json:"primScode"`
	PlmnSrvErrBitMask int    `json:"plmnSrvErrBitMask"`
	ChanIDUl          int    `json:"chanIdUl"`
	TxLevel           int    `json:"txLevel"`
	RxLevel           int    `json:"rxLevel"`
}

type General struct {
	DefaultLanguage           string `json:"defaultLanguage"`
	PRIid                     string `json:"PRIid"`
	ActivationDate            string `json:"activationDate"`
	GenericResetStatus        string `json:"genericResetStatus"`
	ReconditionDate           string `json:"reconditionDate"`
	Manufacturer              string `json:"manufacturer"`
	Model                     string `json:"model"`
	HWversion                 string `json:"HWversion"`
	FWversion                 string `json:"FWversion"`
	ModemFwVersion            string `json:"modemFwVersion"`
	BuildDate                 string `json:"buildDate"`
	BLversion                 string `json:"BLversion"`
	PRIversion                string `json:"PRIversion"`
	TruInstallAvailable       bool   `json:"truInstallAvailable"`
	TruInstallVersion         string `json:"truInstallVersion"`
	IMEI                      string `json:"IMEI"`
	SVN                       string `json:"SVN"`
	MEID                      string `json:"MEID"`
	ESN                       string `json:"ESN"`
	FSN                       string `json:"FSN"`
	ATTDeviceID               string `json:"ATTDeviceId"`
	PackageName               string `json:"packageName"`
	Activated                 bool   `json:"activated"`
	WebAppVersion             string `json:"webAppVersion"`
	AppVersion                string `json:"appVersion"`
	HIDenabled                bool   `json:"HIDenabled"`
	TruInstallEnabled         bool   `json:"truInstallEnabled"`
	TruInstallSupported       bool   `json:"truInstallSupported"`
	TCAaccepted               bool   `json:"TCAaccepted"`
	LEDenabled                bool   `json:"LEDenabled"`
	ShowAdvHelp               bool   `json:"showAdvHelp"`
	ShowWebInfo               bool   `json:"showWebInfo"`
	KeyLockState              string `json:"keyLockState"`
	DevTemperature            int    `json:"devTemperature"`
	VerMajor                  int    `json:"verMajor"`
	VerMinor                  int    `json:"verMinor"`
	Environment               string `json:"environment"`
	CurrTime                  int    `json:"currTime"`
	TimeZoneOffset            int    `json:"timeZoneOffset"`
	UserTzOffsetEnabled       bool   `json:"userTzOffsetEnabled"`
	UserTzOffset              int    `json:"userTzOffset"`
	DeviceName                string `json:"deviceName"`
	UseMetricSystem           bool   `json:"useMetricSystem"`
	FactoryResetButtonEnabled bool   `json:"factoryResetButtonEnabled"`
	FactoryResetLCDSupported  bool   `json:"factoryResetLCDSupported"`
	FactoryResetStatus        string `json:"factoryResetStatus"`
	SPClockStatus             string `json:"SPClockStatus"`
	SetupCompleted            bool   `json:"setupCompleted"`
	WarrantyDateCode          string `json:"warrantyDateCode"`
	LanguageSelected          bool   `json:"languageSelected"`
	UIDateFormat              string `json:"uiDateFormat"`
	UpTime                    int    `json:"upTime"`
	Use24HrTimeFormat         bool   `json:"use24HrTimeFormat"`
	SystemAlertList           struct {
		List []struct {
		} `json:"list"`
		Count int `json:"count"`
	} `json:"systemAlertList"`
	APIVersion        string `json:"apiVersion"`
	CompanyName       string `json:"companyName"`
	UsbDevicesURL     string `json:"usbDevicesURL"`
	NwFoldersURL      string `json:"nwFoldersURL"`
	ConfigURL         string `json:"configURL"`
	ProfileURL        string `json:"profileURL"`
	PinChangeURL      string `json:"pinChangeURL"`
	PortCfgURL        string `json:"portCfgURL"`
	PortFilterURL     string `json:"portFilterURL"`
	WifiACLURL        string `json:"wifiACLURL"`
	SupportedLangList []struct {
		ID        string `json:"id,omitempty"`
		IsCurrent string `json:"isCurrent,omitempty"`
		IsDefault string `json:"isDefault,omitempty"`
		Label     string `json:"label,omitempty"`
		Token1    string `json:"token1,omitempty"`
		Token2    string `json:"token2,omitempty"`
	} `json:"supportedLangList"`
}

type WWAN struct {
	PrlVersion               int    `json:"prlVersion"`
	BandDisablementMaskLTE   string `json:"bandDisablementMaskLTE"`
	BandDisablementMask      string `json:"bandDisablementMask"`
	LTEBandPriority          string `json:"LTEBandPriority"`
	NetScanStatus            string `json:"netScanStatus"`
	MTUSize                  int    `json:"mtuSize"`
	MTUChangeEnabled         bool   `json:"mtuChangeEnabled"`
	LTEeHRPDConfig           string `json:"LTEeHRPDConfig"`
	RoamingEnhancedIndicator int    `json:"roamingEnhancedIndicator"`
	RoamingMode              string `json:"roamingMode"`
	RoamingGuardDom          string `json:"roamingGuardDom"`
	RoamingGuardIntl         string `json:"roamingGuardIntl"`
	RoamingType              string `json:"roamingType"`
	RoamMenuDisplay          bool   `json:"roamMenuDisplay"`
	ERITestMode              bool   `json:"ERITestMode"`
	AutoBandRegionChanged    bool   `json:"autoBandRegionChanged"`
	InactivityCause          int    `json:"inactivityCause"`
	CurrentNWserviceType     string `json:"currentNWserviceType"`
	RegisterRejectCode       int    `json:"registerRejectCode"`
	NetSelEnabled            string `json:"netSelEnabled"`
	NetRegMode               string `json:"netRegMode"`
	Roaming                  bool   `json:"roaming"`
	IPv6                     string `json:"IPv6"`
	IP                       string `json:"IP"`
	RegisterNetworkDisplay   string `json:"registerNetworkDisplay"`
	RAT                      string `json:"RAT"`
	BandRegion               []struct {
		Index   int    `json:"index,omitempty"`
		Name    string `json:"name,omitempty"`
		Current bool   `json:"current,omitempty"`
	} `json:"bandRegion"`
	Autoconnect string `json:"autoconnect"`
	ProfileList []struct {
		Index          int    `json:"index,omitempty"`
		ID             string `json:"id,omitempty"`
		Name           string `json:"name,omitempty"`
		Apn            string `json:"apn,omitempty"`
		Username       string `json:"username,omitempty"`
		Password       string `json:"password,omitempty"`
		AuthType       string `json:"authtype,omitempty"`
		IPaddress      string `json:"ipaddr,omitempty"`
		AccessControl  int    `json:"access_control,omitempty"`
		Type           string `json:"type,omitempty"`
		PDPRoamingType string `json:"pdproamingtype,omitempty"`
	} `json:"profileList"`
	Profile struct {
		Default               string `json:"default"`
		DefaultLTE            string `json:"defaultLTE"`
		PromptForApnSelection bool   `json:"promptForApnSelection"`
	} `json:"profile"`
	PromptForPwd string `json:"promptForPwd"`
	DataUsage    struct {
		Total struct {
			LTEBillingTx  int `json:"lteBillingTx"`
			LTEBillingRx  int `json:"lteBillingRx"`
			CDMABillingTx int `json:"cdmaBillingTx"`
			CDMABillingRx int `json:"cdmaBillingRx"`
			GwBillingTx   int `json:"gwBillingTx"`
			GwBillingRx   int `json:"gwBillingRx"`
			LTELifeTx     int `json:"lteLifeTx"`
			LTELifeRx     int `json:"lteLifeRx"`
			CDMALifeTx    int `json:"cdmaLifeTx"`
			CDMALifeRx    int `json:"cdmaLifeRx"`
			GwLifeTx      int `json:"gwLifeTx"`
			GwLifeRx      int `json:"gwLifeRx"`
		} `json:"total"`
		Server struct {
			AccountType    string `json:"accountType"`
			SubAccountType string `json:"subAccountType"`
		} `json:"server"`
		Remote struct {
			Enabled bool `json:"enabled"`
		} `json:"remote"`
		ServerDataRemaining       int    `json:"serverDataRemaining"`
		ServerDataTransferred     int    `json:"serverDataTransferred"`
		ServerDataTransferredIntl int    `json:"serverDataTransferredIntl"`
		ServerDataValidState      string `json:"serverDataValidState"`
		ServerDaysLeft            int    `json:"serverDaysLeft"`
		ServerErrorCode           string `json:"serverErrorCode"`
		ServerLowBalance          bool   `json:"serverLowBalance"`
		ServerMSISDN              string `json:"serverMsisdn"`
		ServerRechargeURL         string `json:"serverRechargeUrl"`
		DataWarnEnable            bool   `json:"dataWarnEnable"`
		PlanSize                  int    `json:"planSize"`
		PlanDescription           string `json:"planDescription"`
		PrepaidStackedPlans       int    `json:"prepaidStackedPlans"`
		PrepaidStackedPlansIntl   int    `json:"prepaidStackedPlansIntl"`
		PrepaidAccountState       string `json:"prepaidAccountState"`
		AccountType               string `json:"accountType"`
		DisableAutoReset          bool   `json:"disableAutoReset"`
		Share                     struct {
			Enabled               bool   `json:"enabled"`
			DataTransferredOthers int    `json:"dataTransferredOthers"`
			LastSync              string `json:"lastSync"`
		} `json:"share"`
		Generic struct {
			DataLimitValid         bool   `json:"dataLimitValid"`
			UsageHighWarning       int    `json:"usageHighWarning"`
			FallbackSupported      bool   `json:"fallbackSupported"`
			LastSucceeded          string `json:"lastSucceeded"`
			BillingDay             int    `json:"billingDay"`
			NextBillingDate        string `json:"nextBillingDate"`
			LastSync               string `json:"lastSync"`
			BillingCycleRemainder  int    `json:"billingCycleRemainder"`
			BillingCycleLimit      int    `json:"billingCycleLimit"`
			DataTransferred        int    `json:"dataTransferred"`
			DataTransferredRoaming int    `json:"dataTransferredRoaming"`
			LastReset              string `json:"lastReset"`
			UserDisplayFormat      string `json:"userDisplayFormat"`
		} `json:"generic"`
	} `json:"dataUsage"`
	DataTransferred struct {
		Totalb string `json:"totalb"`
		Rxb    string `json:"rxb"`
		Txb    string `json:"txb"`
	} `json:"dataTransferred"`
	NetManualNoCvg       bool   `json:"netManualNoCvg"`
	Connection           string `json:"connection"`
	ConnectionType       string `json:"connectionType"`
	CurrentPSserviceType string `json:"currentPSserviceType"`
	CA                   struct {
		SCCcount int `json:"SCCcount"`
		SCClist  []struct {
		} `json:"SCClist"`
	} `json:"ca"`
	ConnectionText   string `json:"connectionText"`
	SessionDuration  int    `json:"sessDuration"`
	SessionStartTime int    `json:"sessStartTime"`
	SignalStrength   struct {
		RSSI int `json:"rssi"`
		RSCP int `json:"rscp"`
		ECIO int `json:"ecio"`
		RSRP int `json:"rsrp"`
		RSRQ int `json:"rsrq"`
		BARS int `json:"bars"`
		SINR int `json:"sinr"`
	} `json:"signalStrength"`
	DiagInfo []struct {
		LTEAttached       bool   `json:"lteAttached,omitempty"`
		NR5GAttached      bool   `json:"nr5gAttached,omitempty"`
		EndcEnabledConfig bool   `json:"endcEnabledConfig,omitempty"`
		LTESigValid       bool   `json:"ltesigValid,omitempty"`
		LTESigRssi        string `json:"ltesigRssi,omitempty"`
		LTESigRsrp        string `json:"ltesigRsrp,omitempty"`
		LTESigRsrq        string `json:"ltesigRsrq,omitempty"`
		LTESigSnr         string `json:"ltesigSnr,omitempty"`
		NR5GsigValid      bool   `json:"nr5gsigValid,omitempty"`
		NR5GsigRsrp       string `json:"nr5gsigRsrp,omitempty"`
		NR5GsigRsrq       string `json:"nr5gsigRsrq,omitempty"`
		NR5GsigSnr        string `json:"nr5gsigSnr,omitempty"`
	} `json:"diagInfo"`
	LTEBandInfo []struct {
		IsPcc          bool   `json:"isPcc,omitempty"`
		SccID          int    `json:"sccId,omitempty"`
		SccState       int    `json:"sccState,omitempty"`
		Band           int    `json:"band,omitempty"`
		Channel        string `json:"channel,omitempty"`
		PhyCid         string `json:"phyCid,omitempty"`
		DlBandwidth    string `json:"dlBandwidth,omitempty"`
		SccUlCaEnabled string `json:"sccUlCaEnabled,omitempty"`
		SccUlChannel   int    `json:"sccUlChannel,omitempty"`
		SccUlBandwidth string `json:"sccUlBandwidth,omitempty"`
	} `json:"lteBandInfo"`
	Laa struct {
		Supported bool `json:"supported"`
		Enabled   bool `json:"enabled"`
	} `json:"laa"`
}

type Session struct {
	UserRole            string `json:"userRole"`
	Lang                string `json:"lang"`
	HintDisplayPassword string `json:"hintDisplayPassword"`
	SecToken            string `json:"secToken"`
	ClientIP            string `json:"clientIP"`
	SupportedLangList   []struct {
		ID        string `json:"id,omitempty"`
		IsCurrent bool   `json:"isCurrent,omitempty"`
		Label     string `json:"label,omitempty"`
	} `json:"supportedLangList"`
}

type SIM struct {
	PIN PIN `json:"pin"`
	PUK PUK `json:"puk"`
	MEP struct {
	} `json:"mep"`
	PhoneNumber   string `json:"phoneNumber"`
	ICCID         string `json:"iccid"`
	IMSI          string `json:"imsi"`
	SPN           string `json:"SPN"`
	Status        string `json:"status"`
	SprintSimLock int    `json:"sprintSimLock"`
}

type SMS struct {
	Ready          bool      `json:"ready"`
	SendSupported  bool      `json:"sendSupported"`
	SendEnabled    bool      `json:"sendEnabled"`
	UnreadMessages int       `json:"unreadMsgs"`
	AlertSupported bool      `json:"alertSupported"`
	AlertEnabled   bool      `json:"alertEnabled"`
	AlertNumList   string    `json:"alertNumList"`
	MessageCount   int       `json:"msgCount"`
	Messages       []Message `json:"msgs"`
	Trans          []struct {
	} `json:"trans"`
	SendMessage []struct {
		ClientID  string `json:"clientId,omitempty"`
		Enc       string `json:"enc,omitempty"`
		ErrorCode int    `json:"errorCode,omitempty"`
		MsgID     int    `json:"msgId,omitempty"`
		Receiver  string `json:"receiver,omitempty"`
		Status    string `json:"status,omitempty"`
		Text      string `json:"text,omitempty"`
		TxTime    string `json:"txTime,omitempty"`
	} `json:"sendMsg"`
}

type Message struct {
	ID     string `json:"id,omitempty"`
	RxTime string `json:"rxTime,omitempty"`
	Text   string `json:"text,omitempty"`
	Sender string `json:"sender,omitempty"`
	Read   bool   `json:"read,omitempty"`
}

type PUK struct {
	Retry int `json:"retry"`
}

type PIN struct {
	Mode  string `json:"mode"`
	Retry int    `json:"retry"`
}

type AccessControl struct {
	Nlpc struct {
		FilterLevel string `json:"filterLevel"`
		Enabled     bool   `json:"enabled"`
		Username    string `json:"username"`
		MacList     []struct {
		} `json:"macList"`
	} `json:"nlpc"`
	BlockSites              BlockSites `json:"blocksites"`
	SchedulerBlockingActive bool       `json:"schedulerBlockingActive"`
	SchedulerSupported      bool       `json:"schedulerSupported"`
	SchedulerEnabled        bool       `json:"schedulerEnabled"`
	SchedulerDefaultMode    string     `json:"schedulerDefaultMode"`
	SchedulerList           []struct {
	} `json:"schedulerList"`
}

type WebDaemon struct {
	AdminPassword              string `json:"adminPassword"`
	OwnerModeEnabled           bool   `json:"ownerModeEnabled"`
	HideAdminPassword          bool   `json:"hideAdminPassword"`
	HintAnswer                 string `json:"hintAnswer"`
	HintNumber                 int    `json:"hintNumber"`
	OwnerCustomizationsChanged bool   `json:"ownerCustomizationsChanged"`
	OwnerCustomizationsList    []struct {
	} `json:"ownerCustomizationsList"`
}

type BlockSites struct {
	Enabled bool   `json:"enabled"`
	MACMode string `json:"macMode"`
	MACList []struct {
	} `json:"macList"`
	PatternMode string `json:"patternMode"`
	PatternList []struct {
	} `json:"patternList"`
}

type LED struct {
}

type LCD struct {
	BacklightEnabled       bool   `json:"backlightEnabled"`
	BacklightActive        bool   `json:"backlightActive"`
	InactivityTimer        int    `json:"inactivityTimer"`
	InactivityTimerAC      int    `json:"inactivityTimerAC"`
	InactivityTimerUSB     int    `json:"inactivityTimerUSB"`
	BrightnessOnBatt       string `json:"brightnessOnBatt"`
	BrightnessOnUSB        string `json:"brightnessOnUSB"`
	BrightnessOnAC         string `json:"brightnessOnAC"`
	BrightnessOnBattIntVal int    `json:"brightnessOnBattIntVal"`
	BrightnessOnUSBIntVal  int    `json:"brightnessOnUSBIntVal"`
	BrightnessOnACIntVal   int    `json:"brightnessOnACIntVal"`
	BacklightOverride      string `json:"backlightOverride"`
	LockScreenEnable       bool   `json:"lockScreenEnable"`
	LockScreenUsePin       bool   `json:"lockScreenUsePin"`
}

type Insight struct {
	Supported bool   `json:"supported"`
	Enabled   bool   `json:"enabled"`
	Active    bool   `json:"active"`
	Version   string `json:"version"`
}

type Ready struct {
	ShareSupported bool        `json:"shareSupported"`
	ShareEnabled   bool        `json:"shareEnabled"`
	CloudSupported bool        `json:"cloudSupported"`
	CloudEnabled   bool        `json:"cloudEnabled"`
	DeviceShare    DeviceShare `json:"deviceShare"`
	Cloud          struct {
		RegistrationStatus string `json:"registrationStatus"`
	} `json:"cloud"`
}

type DeviceShare struct {
	NetworkDeviceName     string `json:"networkDeviceName"`
	WorkgroupName         string `json:"workgroupName"`
	RemoveUsbDeviceResult string `json:"removeUsbDeviceResult"`
	ExtMediaSupported     bool   `json:"extMediaSupported"`
	USBDevicesInfo        []struct {
	} `json:"usbDevicesInfo"`
	NwFoldersInfo []struct {
	} `json:"nwFoldersInfo"`
}

type Custom struct {
	LastWifiChan      int    `json:"lastWifiChan"`
	DsaLocalURL       string `json:"dsaLocalUrl"`
	HiddenMenuEnabled bool   `json:"hiddenMenuEnabled"`
	HideAdminPassword bool   `json:"hideAdminPassword"`
	CarrierSim        string `json:"carrierSim"`
}

type Modem struct {
	Hostname string
	Username string
	Password string

	token    string
	maxSMSId int

	client *http.Client
}

func New(hostname, username, password string) (*Modem, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cookie jar: %w", err)
	}

	m := &Modem{
		Hostname: hostname,
		Username: username,
		Password: password,
	}

	m.client = &http.Client{
		Jar: jar,
	}

	return m, nil
}

func (m *Modem) baseURL() string {
	return fmt.Sprintf("http://%s", m.Hostname)
}

func (m *Modem) url(path string) string {
	return fmt.Sprintf("%s/%s", m.baseURL(), path)
}

func (m *Modem) urlForm(form string) string {
	return fmt.Sprintf("%s/Forms/%s", m.baseURL(), form)
}

func (m *Modem) Logout() error {
	m.token = ""

	return nil
}

func (m *Modem) Login() error {
	resp, err := m.client.Get(m.url("model.json"))
	if err != nil {
		return fmt.Errorf("failed to get model.json: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed HTTP request for model.json: %s (%d)", resp.Status, resp.StatusCode)
	}

	j := &Model{}
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(j); err != nil {
		return fmt.Errorf("failed to decode model: %w", err)
	}

	m.token = j.Session.SecToken

	if err := m.Config("session.password", m.Password); err != nil {
		return err
	}

	return nil
}

func (m *Modem) Config(key, value string) error {
	return m.postForm("config", map[string]string{
		key: value,
	})
}

func (m *Modem) postForm(f string, cfg map[string]string) error {
	form := url.Values{}
	form.Add("token", m.token)
	form.Add("err_redirect", "/error.json")
	form.Add("ok_redirect", "/success.json")
	for key, value := range cfg {
		form.Add(key, value)
	}

	req, err := http.NewRequest("POST", m.urlForm(f), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get login: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed HTTP request: %s (%d)", resp.Status, resp.StatusCode)
	}

	return nil
}

func (m *Modem) SetWifiSSID(ssid string) error {
	return m.Config("wifi.SSID", ssid)
}

func (m *Modem) SetFailoverMode(mode FailOverMode) error {
	return m.Config("failover.mode", string(mode))
}

func (m *Modem) SetAutoConnectMode(mode AutoConnectMode) error {
	return m.Config("wwan.autoconnect", string(mode))
}

func (m *Modem) Restart() error {
	return m.Config("general.shutdown", "restart")
}

func (m *Modem) FactoryReset() error {
	return m.Config("general.factoryReset", "1")
}

func (m *Modem) ConnectLTE() error {
	return m.Config("wwan.connect", "DefaultProfile")
}

func (m *Modem) DisconnectLTE() error {
	return m.Config("wwan.connect", "Disconnect")
}

func (m *Modem) SendSMS(phone, message string) error {
	return m.postForm("smsSendMsg", map[string]string{
		"sms.sendMsg.receiver": phone,
		"sms.sendMsg.text":     message,
		"sms.sendMsg.clientId": "webUi",
		"action":               "send",
	})
}

func (m *Modem) DeleteAllSMS() error {
	return m.Config("sms.deleteAll", "1")
}

func (m *Modem) DeleteSMS(id int) error {
	return m.Config("sms.deleteId", fmt.Sprint(id))
}

func (m *Modem) GetModel() (*Model, error) {
	u := m.url("model.json")
	r, err := m.client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed to get model.json: %w", err)
	}

	j := &Model{}
	d := json.NewDecoder(r.Body)
	if err := d.Decode(j); err != nil {
		return nil, fmt.Errorf("failed to decode model: %w", err)
	}

	return j, nil
}

func (m *Modem) GetState() (*pb.ModemState, error) {
	n, err := m.GetModel()
	if err != nil {
		return nil, err
	}

	return &pb.ModemState{
		Wwan: &pb.ModemState_WWAN{
			Operator:       n.WWAN.Connection,
			ConnectionText: n.WWAN.ConnectionText,
		},
	}, nil
}
