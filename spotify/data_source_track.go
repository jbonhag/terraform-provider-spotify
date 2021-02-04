package spotify

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	s "github.com/zmb3/spotify"
)

func dataSourceTrack() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrackRead,
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
			"artist_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"artist_name"},
			},
			"artist_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"artist_id"},
			},
		},
	}
}

func dataSourceTrackRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*s.Client)

	var diags diag.Diagnostics

	if v, ok := d.GetOk("id"); ok {
		track, err := findTrackById(c, v.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(track.ID.String())
		d.Set("name", track.Name)
		d.Set("artist_id", track.Artists[0].ID.String())
		d.Set("artist_name", track.Artists[0].Name)

		return diags
	}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)

		if artistID, ok := d.GetOk("artist_id"); ok {
			track, err := findTrackByNameAndArtistID(c, name, artistID.(string))
			if err != nil {
				return diag.FromErr(err)
			}

			d.SetId(track.ID.String())
			d.Set("name", track.Name)
			d.Set("artist_id", track.Artists[0].ID.String())
			d.Set("artist_name", track.Artists[0].Name)

			return diags
		} else if artistName, ok := d.GetOk("artist_name"); ok {
			track, err := findTrackByNameAndArtistName(c, name, artistName.(string))
			if err != nil {
				return diag.FromErr(err)
			}

			d.SetId(track.ID.String())
			d.Set("name", track.Name)
			d.Set("artist_id", track.Artists[0].ID.String())
			d.Set("artist_name", track.Artists[0].Name)

			return diags
		} else {
			track, err := findTrackByName(c, name)
			if err != nil {
				return diag.FromErr(err)
			}

			d.SetId(track.ID.String())
			d.Set("name", track.Name)
			d.Set("artist_id", track.Artists[0].ID.String())
			d.Set("artist_name", track.Artists[0].Name)

			return diags
		}
	}

	return diags
}

func findTrackById(c *s.Client, trackID string) (*s.FullTrack, error) {
	return c.GetTrack(s.ID(trackID))
}

func findTrackByName(c *s.Client, name string) (*s.FullTrack, error) {
	searchResult, err := c.Search(name, s.SearchTypeTrack)
	if err != nil {
		return nil, err
	}

	for _, track := range searchResult.Tracks.Tracks {
		if strings.EqualFold(track.Name, name) {
			return &track, nil
		}
	}

	return nil, err
}

func findTrackByNameAndArtistID(c *s.Client, name string, artistID string) (*s.FullTrack, error) {
	searchResult, err := c.Search(name, s.SearchTypeTrack)
	if err != nil {
		return nil, err
	}

	for _, track := range searchResult.Tracks.Tracks {
		if strings.EqualFold(track.Name, name) && track.Artists[0].ID.String() == artistID {
			return &track, nil
		}
	}

	return nil, err
}

func findTrackByNameAndArtistName(c *s.Client, name string, artistName string) (*s.FullTrack, error) {
	searchResult, err := c.Search(name, s.SearchTypeTrack)
	if err != nil {
		return nil, err
	}

	for _, track := range searchResult.Tracks.Tracks {
		if strings.EqualFold(track.Name, name) && strings.EqualFold(track.Artists[0].Name, artistName) {
			return &track, nil
		}
	}

	return &searchResult.Tracks.Tracks[0], err
}
