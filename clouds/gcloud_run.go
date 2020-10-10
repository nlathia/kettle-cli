package clouds

type GoogleCloudRun struct{}

func (GoogleCloudRun) Build(directory string) error {
	return nil
}

func (GoogleCloudRun) Deploy(directory string) error {
	return nil
}
