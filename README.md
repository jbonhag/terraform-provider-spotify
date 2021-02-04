# terraform-provider-spotify

A Spotify provider for Terraform.

## What can it do?

- Create an empty playlist
- Rename a playlist
- Get a track data source by id (not super useful)
- Create a playlist with tracks
- Update tracks on a playlist
- Find a track by name
- Find a track by artist

## What can't it do?

- Add/remove tracks from your user library

## What could be better?

- Store the token code temporarily (logging in way too much)
- Playlist updating is pretty inefficient, it just removes/adds the entire set
  of tracks on every snapshot
- Fuzzy matching for artist/track names -- right now it just uses `strings.EqualFold` to compare
- The track <-> artist relationship is not correct, should be many:many

## Usage

Create an application in the [Spotify Developer Dashboard][1]. Spotify will
give you a client ID and a client secret. You can store these in the
`SPOTIFY_ID` and `SPOTIFY_SECRET` environment variables, or you can set them
in the provider config:

```hcl
provider "spotify" {
  client_id     = "8f7272f164f0c1e7e7bb680ab0722b70"
  client_secret = "70a3e08ec339460b40cbf18737e93dab"
}
```

When Terraform initializes, it will open a web page to authorize the provider
to use your Spotify account.

[1]: https://developer.spotify.com/dashboard/applications
