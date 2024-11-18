package snowflake_interf

type InterfaceSnowFlake interface {
	GetId() int64
	GetIdAsString() string
}
