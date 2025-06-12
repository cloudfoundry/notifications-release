package services

import "github.com/cloudfoundry/notifications-release/src/notifications/v81/v1/models"

type TemplateFinder struct {
	templatesRepo TemplatesRepo
}

func NewTemplateFinder(templatesRepo TemplatesRepo) TemplateFinder {
	return TemplateFinder{
		templatesRepo: templatesRepo,
	}
}

func (finder TemplateFinder) FindByID(database DatabaseInterface, templateID string) (models.Template, error) {
	template, err := finder.templatesRepo.FindByID(database.Connection(), templateID)
	if err != nil {
		return models.Template{}, err
	}

	return template, err
}
