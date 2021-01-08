package config

// A ConfigChoice is used to enumerate a set of preferences
// that can be selected interactively by the user
type ConfigChoice struct {
	// The label, which is shown in the prompt to the end user
	// The config key: the selection will be stored in viper using this
	Label string
	Key   string

	// Flags so that users can define this choice via an input flag
	// e.g. --cloud <value>
	FlagKey         string
	FlagDescription string
	FlagValue       string

	// A function to collect values if the user does not provide one via a flag
	// A function to validate the choice
	CollectValuesFunc func() (map[string]string, error)
	ValidationFunc    func(string) error
}
