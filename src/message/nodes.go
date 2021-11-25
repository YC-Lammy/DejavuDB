package message

type NodeType byte

const (
	Reserved NodeType = iota
	Manager
	TeamLead
	TeamMember
)
