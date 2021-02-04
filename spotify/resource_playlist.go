package spotify

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	s "github.com/zmb3/spotify"
)

func resourcePlaylist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaylistCreate,
		ReadContext:   resourcePlaylistRead,
		UpdateContext: resourcePlaylistUpdate,
		DeleteContext: resourcePlaylistDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"track_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePlaylistCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*s.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	trackIDs := make([]s.ID, 0)
	for _, trackID := range d.Get("track_ids").([]interface{}) {
		trackIDs = append(trackIDs, s.ID(trackID.(string)))
	}

	user, err := c.CurrentUser()
	if err != nil {
		return diag.FromErr(err)
	}

	playlist, err := c.CreatePlaylistForUser(user.ID, name, "", false)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(playlist.ID))

	c.AddTracksToPlaylist(playlist.ID, trackIDs...)

	resourcePlaylistRead(ctx, d, m)

	return diags
}

func resourcePlaylistRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*s.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	playlistID := d.Id()

	playlist, err := c.GetPlaylist(s.ID(playlistID))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", playlist.Name)

	trackIDs := make([]string, 0)
	for _, track := range playlist.Tracks.Tracks {
		trackIDs = append(trackIDs, string(track.Track.ID))
	}

	d.Set("track_ids", trackIDs)

	return diags
}

func resourcePlaylistUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*s.Client)

	playlistID := d.Id()

	if d.HasChange("name") {
		newName := d.Get("name").(string)

		c.ChangePlaylistName(s.ID(playlistID), newName)
	}

	if d.HasChange("track_ids") {
		playlist, err := c.GetPlaylist(s.ID(playlistID))
		if err != nil {
			return diag.FromErr(err)
		}

		oldTrackIDs := make([]s.ID, 0)
		for _, track := range playlist.Tracks.Tracks {
			oldTrackIDs = append(oldTrackIDs, track.Track.ID)
		}
		c.RemoveTracksFromPlaylist(s.ID(playlistID), oldTrackIDs...)

		newTrackIDs := make([]s.ID, 0)
		for _, trackID := range d.Get("track_ids").([]interface{}) {
			newTrackIDs = append(newTrackIDs, s.ID(trackID.(string)))
		}
		c.AddTracksToPlaylist(s.ID(playlistID), newTrackIDs...)
	}

	return resourcePlaylistRead(ctx, d, m)
}

// You can't delete a playlist through the API
func resourcePlaylistDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
