package docker

import (
	"fmt"

	"github.com/megamsys/vertice/repository"
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

//can we list the marketplaces items for dockerimages
//we need add a  db.FetchCollection to list all from the bucket.
func listBoxImages(appName string) ([]string, error) {
	return []string{}, nil
}

//list all the tosca images from marketplace and figure out if this is valid.
func isValidBoxImage(appName, imageId string) (bool, error) {
	/*images, err := listBoxImages(appName)

	if err != nil {
		return false, err
	}
	for _, img := range images {
		if img == imageId {
			return true, nil
		}
	}
	return false, nil
	*/
	return true, nil

}

// getBuildImage returns the image name from box or plaftorm.
func (p *dockerProvisioner) getBuildImage(re *repository.Repo, version string) string {
	if p.usePlatformImage(re) {
		return platformImageName(re.Gitr())
	}
	return fmt.Sprintf("%s:%s", re.Gitr(), version)
}

func platformImageName(platformName string) string {
	return fmt.Sprintf("%s:latest", platformName)
}

func (p *dockerProvisioner) usePlatformImage(re *repository.Repo) bool {
	return true
}

func (p *dockerProvisioner) cleanImage(appName, imgName string) {
	//	shouldRemove := true
	err := p.Cluster().RemoveImage(imgName)
	if err != nil {
		//	shouldRemove = false
		//log.Errorf("Ignored error removing old image %q: %s. Image kept on list to retry later.",
		//	imgName, err.Error())
	}
	err = p.Cluster().RemoveFromRegistry(imgName)
	if err != nil {
		//shouldRemove = false
		//log.Errorf("Ignored error removing old image from registry %q: %s. Image kept on list to retry later.",
		//	imgName, err.Error())
	}
}
