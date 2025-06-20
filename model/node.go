package model

import (
	"time"
)

type CPUSlot struct {
	Slot      string  `json:"slot"`
	Vendor    string  `json:"vendor"`
	Model     string  `json:"model"`
	Frequency float32 `json:"frequency"`
}

type MemorySlot struct {
	Slot      string  `json:"slot"`
	Vendor    string  `json:"vendor"`
	Model     string  `json:"model"`
	Frequency float32 `json:"frequency"`
	Size      float32 `json:"size"`
}

type DiskSlot struct {
	Slot   string  `json:"slot"`
	Vendor string  `json:"vendor"`
	Model  string  `json:"model"`
	Size   float32 `json:"size"`
}

type NIC struct {
	MAC     string `json:"mac"`
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
}

type Node struct {
	BaseModel `bson:",inline"`

	UUID string `json:"uuid" bson:"uuid"`

	HostName string   `json:"hostname" bson:"hostname"`
	IP       string   `json:"ip" bson:"ip"`
	IPs      []string `json:"ips" bson:"ips"`

	CPUSlot    []CPUSlot    `json:"cpuslot" bson:"cpuslot"`
	MemorySlot []MemorySlot `json:"memoryslot" bson:"memoryslot"`
	DiskSlot   []DiskSlot   `json:"diskslot" bson:"diskslot"`
	NIC        []NIC        `json:"nic" bson:"nic"`

	HeartBeat       time.Time `json:"heartbeat" bson:"heartbeat"`
	HeartBeatStatus bool      `json:"heartbeat_status" bson:"heartbeat_status"`
}

func (model *Node) TableName() string {
	return "node"
}
