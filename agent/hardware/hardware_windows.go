package hardware

import (
	"github.com/StackExchange/wmi"
	"log"
	"strings"
)

type Win32_Processor struct {
	Name                      string
	Manufacturer              string
	Description               string
	DeviceID                  string
	CurrentClockSpeed         uint32
	MaxClockSpeed             uint32
	ExtClock                  uint32
	L2CacheSize               uint32
	L2CacheSpeed              uint32
	L3CacheSize               uint32
	L3CacheSpeed              uint32
	NumberOfCores             uint32
	NumberOfLogicalProcessors uint32
	ThreadCount               uint32
	Architecture              uint16
	AddressWidth              uint16
	DataWidth                 uint16
	Family                    uint16
	Stepping                  string
	Revision                  uint16
	CpuStatus                 uint16
	PowerManagementSupported  bool
	VoltageCaps               uint32
	SocketDesignation         string
	ProcessorType             uint16
	Role                      string
}

func GetCPUInfo() (info []*CPUInfo, err error) {
	var processors []Win32_Processor

	q := wmi.CreateQuery(&processors, "")
	err = wmi.Query(q, &processors)
	if err != nil {
		log.Fatalf("WMI 查询失败: %v", err)
		return nil, err
	}
	if len(processors) == 0 {
		log.Fatal("未找到 CPU 信息")
		return nil, err
	}
	for _, processor := range processors {
		info = append(info, &CPUInfo{
			DeviceID:                  processor.DeviceID,
			Manufacturer:              processor.Manufacturer,
			Socket:                    processor.SocketDesignation,
			MaxClockSpeed:             processor.MaxClockSpeed,
			ExtClock:                  processor.ExtClock,
			L2CacheSize:               processor.L2CacheSize,
			L3CacheSize:               processor.L3CacheSize,
			NumberOfCores:             processor.NumberOfCores,
			NumberOfLogicalProcessors: processor.NumberOfLogicalProcessors,
			ThreadCount:               processor.ThreadCount,
			AddressWidth:              processor.AddressWidth,
			DataWidth:                 processor.DataWidth,
		})
	}
	return
}

type Win32_PhysicalMemory struct {
	BankLabel            string // 内存插槽标签
	Capacity             uint64 // 内存容量（字节）
	DataWidth            uint16 // 数据位宽
	Speed                uint32 // 运行频率（MHz）
	Manufacturer         string // 制造商（品牌）
	PartNumber           string // 部件号
	SerialNumber         string // 序列号
	ConfiguredClockSpeed uint32 // 配置频率（MHz）
	DeviceLocator        string // 设备位置
	FormFactor           uint16 // 物理规格
	MemoryType           uint16 // 内存类型
	MinVoltage           uint32 // 最小电压
	MaxVoltage           uint32 // 最大电压
	SMBIOSMemoryType     uint16 // SMBIOS 内存类型
}

// 内存类型映射表
var memoryTypeMap = map[uint16]string{
	0:  "Unknown",
	1:  "Other",
	2:  "DRAM",
	3:  "Synchronous DRAM",
	4:  "Cache DRAM",
	5:  "EDO",
	6:  "EDRAM",
	7:  "VRAM",
	8:  "SRAM",
	9:  "RAM",
	10: "ROM",
	11: "Flash",
	12: "EEPROM",
	13: "FEPROM",
	14: "EPROM",
	15: "CDRAM",
	16: "3DRAM",
	17: "SDRAM",
	18: "SGRAM",
	19: "RDRAM",
	20: "DDR",
	21: "DDR2",
	22: "DDR2 FB-DIMM",
	24: "DDR3",
	25: "FBD2",
	26: "DDR4",
	27: "DDR5",
}

var formFactorMap = map[uint16]string{
	0:  "Unknown",
	1:  "Other",
	2:  "SIP",
	3:  "DIP",
	4:  "ZIP",
	5:  "SOJ",
	6:  "Proprietary",
	7:  "SIMM",
	8:  "DIMM",
	9:  "TSOP",
	10: "PGA",
	11: "RIMM",
	12: "SODIMM", // 笔记本内存
	13: "SRIMM",
	14: "SMD",
	15: "SSMP",
	16: "QFP",
	17: "TQFP",
	18: "SOIC",
	19: "LCC",
	20: "PLCC",
	21: "BGA",
	22: "FPBGA",
	23: "LGA",
}

func getMemoryBrand(manufacturer string) string {
	manufacturer = strings.TrimSpace(manufacturer)

	// 常见品牌映射
	brandMap := map[string]string{
		"0098": "Kingston",
		"029E": "Corsair",
		"04CD": "G.Skill",
		"04CB": "Crucial",
		"04F4": "Samsung",
		"059B": "Micron",
		"80AD": "Hynix",
		"80CE": "SK Hynix",
		"8551": "Nanya",
		"2C00": "AMD",
	}

	// 如果是编码形式，尝试转换
	if len(manufacturer) == 4 {
		if brand, ok := brandMap[strings.ToUpper(manufacturer)]; ok {
			return brand
		}
	}

	// 返回原始值（可能已经是可读名称）
	if manufacturer == "" {
		return "未知品牌"
	}
	return manufacturer
}

// 根据部件号推断内存颗粒品牌
func inferChipBrand(partNumber string) string {
	partNumber = strings.ToUpper(partNumber)

	// 常见颗粒品牌识别特征
	switch {
	case strings.Contains(partNumber, "SEC") || strings.Contains(partNumber, "SAMSUNG"):
		return "三星(Samsung)"
	case strings.Contains(partNumber, "HYNIX") || strings.Contains(partNumber, "HYK0") || strings.Contains(partNumber, "H5AN"):
		return "海力士(Hynix)"
	case strings.Contains(partNumber, "MICRON") || strings.Contains(partNumber, "MT") || strings.Contains(partNumber, "MTA"):
		return "美光(Micron)"
	case strings.Contains(partNumber, "NANYA") || strings.Contains(partNumber, "NT"):
		return "南亚(Nanya)"
	case strings.Contains(partNumber, "ELPIDA") || strings.Contains(partNumber, "EBJ"):
		return "尔必达(Elpida)"
	case strings.Contains(partNumber, "SPEC") || strings.Contains(partNumber, "SPECTEK"):
		return "Spectek"
	default:
		return ""
	}
}

func GetMemoryInfo() (info []MemoryInfo, err error) {
	var memory []Win32_PhysicalMemory

	// 执行 WMI 查询
	err = wmi.Query("SELECT * FROM Win32_PhysicalMemory", &memory)
	if err != nil {
		log.Fatalf("内存信息查询失败: %v", err)
		return nil, err
	}

	if len(memory) == 0 {
		log.Fatal("未找到内存信息")
		return nil, err
	}
	for _, mem := range memory {
		i := MemoryInfo{
			BankLabel:            mem.BankLabel,
			DeviceLocator:        mem.DeviceLocator,
			Capacity:             mem.Capacity,
			Speed:                mem.Speed,
			ConfiguredClockSpeed: mem.ConfiguredClockSpeed,
			DataWidth:            mem.DataWidth,
			Manufacturer:         getMemoryBrand(mem.Manufacturer),
			PartNumber:           getMemoryBrand(mem.PartNumber),
			SerialNumber:         getMemoryBrand(mem.SerialNumber),
		}
		if formFactor, ok := formFactorMap[mem.FormFactor]; ok {
			i.FormFactor = formFactor
		}
		if memType, ok := memoryTypeMap[mem.MemoryType]; ok {
			i.MemoryType = memType
		} else if mem.SMBIOSMemoryType > 0 {
			if memType, ok := memoryTypeMap[mem.SMBIOSMemoryType]; ok {
				i.MemoryType = memType
			}
		}
		if chipBrand := inferChipBrand(i.Manufacturer); chipBrand != "" {
			i.ChipBrand = chipBrand
		}
		info = append(info, i)
	}
	return info, nil

}
