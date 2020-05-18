package session

//FlagArguments is a copy from the main package, this should be modified in order not to repeat
type flagArguments struct {
	TimeOut     int
	TargetURL   string
	HardLimit   int
	DomainScope string
	NWorkers    int
	FileName    string
	StartZMQ    bool
}
