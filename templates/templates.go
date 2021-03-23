package templates

func GetTemplate(templatePath string) (string, bool, error) {
	// Match on a local path first
	exists, err := PathExists(templatePath)
	if err != nil {
		return "", false, err
	}
	if exists {
		return templatePath, false, nil
	}

	// Match against a github repo & clone the repo to a tmp directory
	if isGitRepository(templatePath) {
		tempDirectory, err := cloneRepository(templatePath)
		return tempDirectory, true, err
	}

	// Look for the template in the kettle-templates monorepo
	tempDirectory, err := searchTemplates(templatePath)
	if err != nil {
		return "", false, err
	}
	return tempDirectory, true, nil
}
