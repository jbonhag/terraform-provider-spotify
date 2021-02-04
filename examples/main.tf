terraform {
  required_providers {
    spotify = {
      version = "0.1"
      source  = "hashicorp.com/jbonhag/spotify"
    }
  }
}

provider "spotify" {
}

data "spotify_track" "mmmbop" {
  name = "MMMBop"
  #id = "0lnxrQAd9ZxbhBBe7d8FO8"
}

data "spotify_track" "mmmbop_single_version" {
  name = "MMMBop - Single Version"
}

data "spotify_artist" "scary_pockets" {
  name = "Scary Pockets"
}

data "spotify_track" "scary_pockets_mmmbop" {
  artist_id = data.spotify_artist.scary_pockets.id
  name      = "Mmmbop"
}

data "spotify_track" "vsq_mmmbop" {
  artist_name = "Vitamin String Quartet"
  name        = "Mmmbop"
}

resource "spotify_playlist" "mmmbop" {
  name = "Nothing but MMMBop"
  track_ids = [
    data.spotify_track.mmmbop.id,
    data.spotify_track.scary_pockets_mmmbop.id,
    data.spotify_track.mmmbop_single_version.id,
    data.spotify_track.vsq_mmmbop.id,
  ]
}
