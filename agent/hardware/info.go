package hardware

type CPUInfo struct {
	DeviceID                  string `json:"device_id"`
	Manufacturer              string `json:"manufacturer"`
	Socket                    string `json:"socket"`
	MaxClockSpeed             uint32 `json:"max_clock_speed"`
	ExtClock                  uint32 `json:"ext_clock"`
	L2CacheSize               uint32 `json:"l2_cache_size"`
	L3CacheSize               uint32 `json:"l3_cache_size"`
	NumberOfCores             uint32 `json:"number_of_cores"`
	NumberOfLogicalProcessors uint32 `json:"number_of_logical_processors"`
	ThreadCount               uint32 `json:"thread_count"`
	ArchName                  string `json:"arch_name"`
	AddressWidth              uint16 `json:"address_width"`
	DataWidth                 uint16 `json:"data_width"`
}

type MemoryInfo struct {
	BankLabel            string `json:"bank_label"`
	DeviceLocator        string `json:"device_locator"`
	Capacity             uint64 `json:"capacity"`
	Speed                uint32 `json:"speed"`
	ConfiguredClockSpeed uint32 `json:"configured_clock_speed"`
	DataWidth            uint16 `json:"data_width"`
	Manufacturer         string `json:"manufacturer"`
	PartNumber           string `json:"part_number"`
	SerialNumber         string `json:"serial_number"`
	FormFactor           string `json:"form_factor"`
	MemoryType           string `json:"memory_type"`
	ChipBrand            string `json:"chip_brand"`
}
