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
          export PGDATA=temp/database
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
          fi

          export ASSETS_SERVER_URI="$(realpath ./temp/assets)" # Supposed file server
        '';
      };
    };
}
