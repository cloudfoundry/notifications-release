package services

import "github.com/cloudfoundry/notifications-release/src/notifications/v81/cf"

type OrganizationLoader struct {
	cc cloudController
}

func NewOrganizationLoader(cc cloudController) OrganizationLoader {
	return OrganizationLoader{
		cc: cc,
	}
}

func (loader OrganizationLoader) Load(orgGUID string, token string) (cf.CloudControllerOrganization, error) {
	organization, err := loader.cc.LoadOrganization(orgGUID, token)
	if err != nil {
		return cf.CloudControllerOrganization{}, CCErrorFor(err)
	}

	return organization, nil
}
