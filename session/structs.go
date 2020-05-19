package session

//FlagArguments is a copy from the main package, this should be modified in order not to repeat
type FlagArguments struct {
	TimeOut     int
	TargetURL   string
	HardLimit   int
	DomainScope string
	NWorkers    int
	FileName    string
	CheckEmails bool
}
