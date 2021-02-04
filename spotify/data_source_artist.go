package spotify

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	s "github.com/zmb3/spotify"
)

func dataSourceArtist() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArtistRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"name"},
				Optional:      true,
			},
			"name": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"id"},
				Optional:      true,
			},
		},
	}
}

func dataSourceArtistRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*s.Client)

	var diags diag.Diagnostics

	if v, ok := d.GetOk("id"); ok {
		artist, err := findArtistById(c, v.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(artist.ID.String())
		d.Set("name", artist.Name)
	}

	if v, ok := d.GetOk("name"); ok {
		artist, err := findArtistByName(c, v.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(artist.ID.String())
		d.Set("name", artist.Name)
	}

	return diags
}

func findArtistById(c *s.Client, artistID string) (*s.FullArtist, error) {
	return c.GetArtist(s.ID(artistID))
}

func findArtistByName(c *s.Client, name string) (*s.FullArtist, error) {
	searchResult, err := c.Search(name, s.SearchTypeArtist)
	if err != nil {
		return nil, err
	}

	for _, artist := range searchResult.Artists.Artists {
		if strings.EqualFold(artist.Name, name) {
			return &artist, nil
		}
	}

	return nil, err
}
