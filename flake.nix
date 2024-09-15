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
        packages = with pkgs; [
          go
          postgresql
          postman
        ];
        shellHook = ''
          export PGPORT=5432
          export PGHOST=localhost
          export PGDATA=./tmp/database
          export PGSUPERUSER=postgres
          export DB_NAME=db1
          export POSTGRES_DSN="host=$PGHOST user=$PGSUPERUSER password= database=$DB_NAME dbname=$DB_NAME port=$PGPORT sslmode=disable"

          if [ ! -d ./tmp/database ]; then
            # Setup development db
            mkdir -p "$PGDATA"
            initdb -D "$PGDATA"
            pg_ctl stop
            pg_ctl -l database.log -o "--unix_socket_directories='$PWD/tmp'" start
            createdb "$DB_NAME"
            dropuser "$PGSUPERUSER" --if-exists
            createuser "$PGSUPERUSER" -s

            # Setup development asset fileserver
            mkdir -p ./tmp/assets/tracks
            mkdir -p ./tmp/assets/albums
            mkdir -p ./tmp/assets/users
          else
            pg_ctl stop
            pg_ctl -l database.log -o "--unix_socket_directories='$PWD/tmp'" start
            createdb "$DB_NAME"
            createuser "$PGSUPERUSER" -s

            psql -c "INSERT INTO albums (id, title, artist_ids, track_ids) VALUES (0, 'Blank album', ARRAY [0], ARRAY [0]);" "$DB_NAME"
            psql -c "INSERT INTO artists (id, name, album_ids, track_ids, is_verified) VALUES (0, 'Blank artist', ARRAY [0], ARRAY [0], false);" "$DB_NAME"
            psql -c "INSERT INTO playlists (id, title, owner_id, is_public, track_ids, deletable) VALUES (0, 'Blank playlist', 0, false, ARRAY [0], false);" "$DB_NAME"
            psql -c "INSERT INTO tracks (id, title, artist_ids, album_id, duration) VALUES (0, 'Blank track', ARRAY [0], 0, 0);" "$DB_NAME"
            psql -c "INSERT INTO users (id, username, password_hash, playlist_ids, is_player_state_public, player_track_id, player_progress, player_device, player_is_shuffle_enabled, player_is_repeat_playlist_enabled, player_is_repeat_track_enabled, player_volume) VALUES (0, 'Blank user', 'Blank hash', ARRAY [0], true, 0, 0, 'Blank device', false, false, false, 0.0);" "$DB_NAME"
          fi
          export ASSETS_SERVER_URI="file://$(realpath ./tmp/assets)" # Supposed file server
        '';
      };
    };
}
