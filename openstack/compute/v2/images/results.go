package images

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// GetResult temporarily stores a Get response.
type GetResult struct {
	gophercloud.CommonResult
}

// Extract interprets a GetResult as an Image.
func (gr GetResult) Extract() (*Image, error) {
	if gr.Err != nil {
		return nil, gr.Err
	}

	var decoded struct {
		Image Image `mapstructure:"image"`
	}

	err := mapstructure.Decode(gr.Resp, &decoded)
	return &decoded.Image, err
}

// Image is used for JSON (un)marshalling.
// It provides a description of an OS image.
type Image struct {
	// ID contains the image's unique identifier.
	ID string

	Created string

	// MinDisk and MinRAM specify the minimum resources a server must provide to be able to install the image.
	MinDisk int
	MinRAM  int

	// Name provides a human-readable moniker for the OS image.
	Name string

	// The Progress and Status fields indicate image-creation status.
	// Any usable image will have 100% progress.
	Progress int
	Status   string

	Updated string
}

// ImagePage contains a single page of results from a List operation.
// Use ExtractImages to convert it into a slice of usable structs.
type ImagePage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no Image results.
func (page ImagePage) IsEmpty() (bool, error) {
	images, err := ExtractImages(page)
	if err != nil {
		return true, err
	}
	return len(images) == 0, nil
}

// NextPageURL uses the response's embedded link reference to navigate to the next page of results.
func (page ImagePage) NextPageURL() (string, error) {
	type link struct {
		Href string `mapstructure:"href"`
		Rel  string `mapstructure:"rel"`
	}
	type resp struct {
		Links []link `mapstructure:"images_links"`
	}

	var r resp
	err := mapstructure.Decode(page.Body, &r)
	if err != nil {
		return "", err
	}

	var url string
	for _, l := range r.Links {
		if l.Rel == "next" {
			url = l.Href
		}
	}
	if url == "" {
		return "", nil
	}

	return url, nil
}

// ExtractImages converts a page of List results into a slice of usable Image structs.
func ExtractImages(page pagination.Page) ([]Image, error) {
	casted := page.(ImagePage).Body
	var results struct {
		Images []Image `mapstructure:"images"`
	}

	err := mapstructure.Decode(casted, &results)
	return results.Images, err
}
