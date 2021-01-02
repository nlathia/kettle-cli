package preferences

// A ConfigChoice is used to enumerate a set of preferences
// that can be selected interactively by the user
type ConfigChoice struct {
	// The label, which is shown in the prompt to the end user
	Label string

	// The config key: the selection will be stored in viper using this
	Key string

	// Flags so that users can define this choice via an input flag
	// e.g. --cloud <value>
	FlagKey         string
	FlagValue       string
	FlagDescription string

	// A function to collect values if the user does not provide one via a flag
	CollectValuesFunc func() (map[string]string, error)

	// A function to validate the choice
	ValidationFunc func(string) error
}
