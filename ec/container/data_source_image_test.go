package container_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	containerv1 "gitlab.com/nitrado/b2b/ec/core/pkg/api/container/v1"
)

func TestDataSourceImages(t *testing.T) {
	img1 := &containerv1.Image{
		ImageObjectMeta: containerv1.ImageObjectMeta{
			ObjectMeta: metav1.ObjectMeta{
				Name:             "my-image-1",
				CreatedTimestamp: time.Now().Add(-1 * time.Hour),
			},
			Branch: "my-branch",
		},
		Spec: containerv1.ImageSpec{
			Image: "my-image-name",
			Tag:   "my-tag1",
		},
	}
	img2 := &containerv1.Image{
		ImageObjectMeta: containerv1.ImageObjectMeta{
			ObjectMeta: metav1.ObjectMeta{
				Name:             "my-image-2",
				CreatedTimestamp: time.Now(),
			},
			Branch: "my-branch",
		},
		Spec: containerv1.ImageSpec{
			Image: "my-image-name",
			Tag:   "my-tag2",
		},
	}

	pf, _ := providertest.SetupProviderFactories(t, img1, img2)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageNameConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_container_image.by_name", "metadata.0.object_meta.0.name", "my-image-1"),
					resource.TestCheckResourceAttr("data.ec_container_image.by_name", "metadata.0.branch", "my-branch"),
					resource.TestCheckResourceAttr("data.ec_container_image.by_name", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_container_image.by_name", "spec.0.image", "my-image-name"),
					resource.TestCheckResourceAttr("data.ec_container_image.by_name", "spec.0.tag", "my-tag1"),
				),
			},
			{
				Config: testDataSourceImageSpecConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_container_image.by_image", "metadata.0.object_meta.0.name", "my-image-2"),
					resource.TestCheckResourceAttr("data.ec_container_image.by_image", "metadata.0.branch", "my-branch"),
					resource.TestCheckResourceAttr("data.ec_container_image.by_image", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_container_image.by_image", "spec.0.image", "my-image-name"),
					resource.TestCheckResourceAttr("data.ec_container_image.by_image", "spec.0.tag", "my-tag2"),
				),
			},
		},
	})
}

func testDataSourceImageNameConfigRead() string {
	return `data "ec_container_image" "by_name" {
  metadata {
    object_meta {
      name   = "my-image-1"
    }
    branch = "my-branch" 
  }
}
`
}

func testDataSourceImageSpecConfigRead() string {
	return `data "ec_container_image" "by_image" {
  metadata {
    branch = "my-branch" 
  }
  spec {
    image = "my-image-name"
  }
}
`
}
