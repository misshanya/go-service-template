{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };

        goVersion = pkgs.go_1_25;
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            goVersion

            gopls
            gotools
            delve
            golangci-lint

            air

            sqlc
            goose
            go-swag
            oapi-codegen
            go-task
          ];

          shellHook = ''
            export GOROOT="${goVersion}/share/go"
          '';
        };
      }
    );
}
