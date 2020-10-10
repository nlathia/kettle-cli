package clouds

type GoogleCloudFunction struct{}

func (GoogleCloudFunction) Build(directory string) error {
	return nil
}

func (GoogleCloudFunction) Deploy(directory string) error {
	return nil
}
