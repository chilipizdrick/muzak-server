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
          postgresql
        ];
        shellHook = ''
          export PGDATA=temp/database
          export PGHOST=localhost
          export POSTGRES_DSN="host=$PGHOST user=postgres password= database=db1 dbname=db1 port=5432 sslmode=disable"
          export ASSETS_SERVER_URI="$(realpath ./temp/assets)" # Supposed file server
          pg_ctl -l database.log -o "--unix_socket_directories='$PWD/temp'" start
        '';
      };
    };
}
