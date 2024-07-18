{
  description = "Standart developer shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
  };

  outputs = {nixpkgs, ...}: let
    pkgs = import nixpkgs {
      system = "x86_64-linux";
      config.allowUnfree = true;
    };
  in
    with pkgs; {
      devShells.x86_64-linux.default = mkShell {
        name = "Standart developer shell";
        buildInputs = with pkgs; [
          go
          postgresql
        ];
        shellHook = ''
          export PGDATA=./temp/database
          export PGHOST=localhost
          export POSTGRES_DSN="host=$PGHOST user=postgres password= database=db1 dbname=db1 port=5432 sslmode=disable"

          if [ ! -d ./temp/database ]; then
            # Setup development db
            mkdir -p ./temp/database
            initdb -D ./temp/database && sleep 1
            pg_ctl -l database.log -o "--unix_socket_directories='$PWD/temp'" start
            createdb db1
            dropuser postgres
            createuser postgres -s

            # Setup development asset fileserver
            mkdir -p ./temp/assets/tracks
            mkdir -p ./temp/assets/albums
            mkdir -p ./temp/assets/users
          else
            pg_ctl -l database.log -o "--unix_socket_directories='$PWD/temp'" start

            psql -c "INSERT INTO albums (id, title, artist_ids, track_ids) VALUES (0, 'Blank album', ARRAY [0], ARRAY [0]);" db1
            psql -c "INSERT INTO artists (id, name, album_ids, track_ids, is_verified) VALUES (0, 'Blank artist', ARRAY [0], ARRAY [0], false);" db1
            psql -c "INSERT INTO playlists (id, title, owner_id, is_public, track_ids) VALUES (0, 'Blank playlist', 0, false, ARRAY [0]);" db1
            psql -c "INSERT INTO tracks (id, title, artist_ids, album_id, duration) VALUES (0, 'Blank track', ARRAY [0], 0, 0);" db1
            psql -c "INSERT INTO users (id, username, password_hash, playlist_ids, player_track_id, player_progress, player_device, player_shuffle_enabled, player_repeat_playlist_enabled, player_repeat_track_enabled, player_volume) VALUES (0, 'Blank user', 'Blank hash', ARRAY [0], 0, 0, 'Blank device', false, false, false, 0.0);" db1
          fi

          export ASSETS_SERVER_URI="$(realpath ./temp/assets)" # Supposed file server
        '';
      };
    };
}
