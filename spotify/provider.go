package spotify

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	s "github.com/zmb3/spotify"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = s.NewAuthenticator(redirectURI, s.ScopePlaylistModifyPrivate)
	ch    = make(chan *s.Client)
	state = "abc123"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SPOTIFY_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SPOTIFY_SECRET", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"spotify_playlist": resourcePlaylist(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"spotify_track":  dataSourceTrack(),
			"spotify_artist": dataSourceArtist(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	auth.SetAuthInfo(clientID, clientSecret)

	url := auth.AuthURL(state)
	cmd := exec.Command("open", url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// wait for auth to complete
	c := <-ch
	return c, diags
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
