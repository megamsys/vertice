package docker

import (
	"fmt"
	"strings"
)

type boxImages struct {
	BoxName string
	Images  []string
	Count   int
}

type ImageMetadata struct {
	Name string
}

func boxImageName(appName string) string {
	return fmt.Sprintf("%s", appName)
}

func listBoxImages(appName string) ([]string, error) {
	//call marketplaces bucket
	//we need to improve db.FetchCollection code to list all from the bucket.
	return []string{}, nil
}

//list all the tosca images from marketplace and figure out if this is valid.
func isValidBoxImage(appName, imageId string) (bool, error) {
	images, err := listBoxImages(appName)

	if err != nil {
		return false, err
	}
	for _, img := range images {
		if img == imageId {
			return true, nil
		}
	}
	//return false, nil
	return true, nil

}


// getBuildImage returns the image name from box or plaftorm.
func (p *dockerProvisioner) getBuildImage(box *provision.Box) string {
	if p.usePlatformImage(box) {
		return platformImageName(box.GetPlatform())
	}
	return box.GetFullName()
}


func platformImageName(platformName string) string {
	return fmt.Sprintf("%s/%s:latest", basicImageName(), platformName)
}

func basicImageName() string {
	parts := make([]string, 0, 2)
	registry, _ := "registry.dockerhub.com"
	if registry != "" {
		parts = append(parts, registry)
	}
	repoNamespace, _ := "megam"
	parts = append(parts, repoNamespace)
	return strings.Join(parts, "/")
}

func (p *dockerProvisioner) usePlatformImage(app provision.App) bool {
	return false
}

func (p *dockerProvisioner) cleanImage(appName, imgName string) {
	shouldRemove := true
	err := p.Cluster().RemoveImage(imgName)
	if err != nil {
		shouldRemove = false
		//log.Errorf("Ignored error removing old image %q: %s. Image kept on list to retry later.",
		//	imgName, err.Error())
	}
	err = p.Cluster().RemoveFromRegistry(imgName)
	if err != nil {
		shouldRemove = false
		//log.Errorf("Ignored error removing old image from registry %q: %s. Image kept on list to retry later.",
		//	imgName, err.Error())
	}
}
